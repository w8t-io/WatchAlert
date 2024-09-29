package consumer

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/sender"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/hash"
)

type Consume struct {
	ctx *ctx.Context
	sync.RWMutex
	// 从 Redis 中读取当前告警事件提到内存做处理.
	alertsMap                  map[string][]models.AlertCurEvent
	preStoreFiringAlertEvents  map[string][]models.AlertCurEvent
	preStoreRecoverAlertEvents map[string][]models.AlertCurEvent
	Timing                     map[string]int
}

type InterEvalConsume interface {
	Run()
}

func NewInterEvalConsumeWork(ctx *ctx.Context) InterEvalConsume {
	return &Consume{
		ctx:                        ctx,
		alertsMap:                  make(map[string][]models.AlertCurEvent),
		preStoreFiringAlertEvents:  make(map[string][]models.AlertCurEvent),
		preStoreRecoverAlertEvents: make(map[string][]models.AlertCurEvent),
		Timing:                     make(map[string]int),
	}
}

// Run 启动告警消费进程
func (ec *Consume) Run() {
	go func() {
		for {
			ec.processAlerts()
			time.Sleep(time.Second)
		}
	}()
}

// 处理告警的主循环
func (ec *Consume) processAlerts() {
	alertKeys := process.GetRedisFiringKeys(ec.ctx)
	ec.loadAlertsToMem(alertKeys)
	for key, alerts := range ec.alertsMap {
		if len(alerts) == 0 {
			continue
		}
		waitTime := ec.calculateWaitTime(key)
		if ec.Timing[key] >= waitTime {
			curEvents := ec.filterAlerts(alerts)
			ec.fireAlertEvent(curEvents)
			ec.clear(key)
		}
		ec.Timing[key]++
	}
}

// 加载告警到内存
func (ec *Consume) loadAlertsToMem(alertKeys []string) {
	for _, key := range alertKeys {
		alert := ec.ctx.Redis.Event().GetCache(key)
		if alert.Fingerprint != "" {
			ec.addAlertToRuleIdMap(alert)
		}
	}
}

// 根据告警的状态计算等待时间
func (ec *Consume) calculateWaitTime(key string) int {
	alert := ec.ctx.Redis.Event().GetCache(key)
	if alert.LastSendTime == 0 {
		return global.Config.Server.AlarmConfig.GroupWait
	}
	return global.Config.Server.AlarmConfig.GroupInterval
}

// 告警事件提取到内存中
func (ec *Consume) addAlertToRuleIdMap(alert models.AlertCurEvent) {
	ec.Lock()
	defer ec.Unlock()

	ec.alertsMap[alert.RuleId] = append(ec.alertsMap[alert.RuleId], alert)
}

// 清楚本地缓存
func (ec *Consume) clear(ruleId string) {
	ec.Lock()
	defer ec.Unlock()

	delete(ec.alertsMap, ruleId)
	delete(ec.preStoreFiringAlertEvents, ruleId)
	delete(ec.preStoreRecoverAlertEvents, ruleId)
	ec.Timing[ruleId] = 0
}

// 过滤告警
func (ec *Consume) filterAlerts(alerts []models.AlertCurEvent) map[string][]models.AlertCurEvent {
	var (
		newAlertsMap = make(map[string][]models.AlertCurEvent)
		latestAlert  = make(map[string]models.AlertCurEvent)
	)

	// 基于指纹去重，保留最新的告警
	for _, alert := range alerts {
		if existingAlert, exists := latestAlert[alert.Fingerprint]; !exists || alert.LastEvalTime > existingAlert.LastEvalTime {
			latestAlert[alert.Fingerprint] = alert
		}
	}

	// 进一步处理重复通知
	for _, alert := range latestAlert {
		if !alert.IsRecovered && (alert.LastSendTime == 0 || alert.LastEvalTime >= alert.LastSendTime+alert.RepeatNoticeInterval*60) {
			newAlertsMap[alert.RuleId] = append(newAlertsMap[alert.RuleId], alert)
		} else if alert.IsRecovered {
			newAlertsMap[alert.RuleId] = append(newAlertsMap[alert.RuleId], alert)
		}
	}
	return newAlertsMap
}

// 触发告警通知
func (ec *Consume) fireAlertEvent(alertsMap map[string][]models.AlertCurEvent) {
	for _, alerts := range alertsMap {
		for _, alert := range alerts {
			ec.addAlertToGroup(alert)
			if alert.IsRecovered {
				ec.removeAlertFromCache(alert)
				err := process.RecordAlertHisEvent(ec.ctx, alert)
				if err != nil {
					global.Logger.Sugar().Error(err.Error())
					return
				}
			}
		}
	}
	ec.handleAlert(ec.preStoreFiringAlertEvents)
	ec.handleAlert(ec.preStoreRecoverAlertEvents)
}

// 删除缓存
func (ec *Consume) removeAlertFromCache(alert models.AlertCurEvent) {
	key := alert.GetFiringAlertCacheKey()
	ec.ctx.Redis.Event().DelCache(key)
}

// 添加告警到组(分组)
func (ec *Consume) addAlertToGroup(alert models.AlertCurEvent) {
	// 如果没有定义通知组，则直接添加到 ruleId 组中
	if alert.NoticeGroup == nil || len(alert.NoticeGroup) == 0 {
		ec.addAlertToGroupByRuleId(alert)
		return
	}

	// 遍历所有的 Metric
	matched := false
	for key, value := range alert.Metric {
		// 遍历所有的通知组
		for _, noticeGroup := range alert.NoticeGroup {
			// 如果当前 Metric 的 key 和 value 与通知组中的相匹配
			if noticeGroup["key"] == key && noticeGroup["value"] == value.(string) {
				// 计算分组的 ID 并添加警报到对应的组
				groupId := ec.calculateGroupHash(key, value.(string))
				ec.addAlertToGroupByGroupId(groupId+"_"+alert.RuleId, alert)
				matched = true
				break
			}
		}
		if matched {
			break
		}
	}

	// 如果没有找到任何匹配的组，则添加到 ruleId 组中
	if !matched {
		ec.addAlertToGroupByRuleId(alert)
	}
}

// 以Id作为key添加到组
func (ec *Consume) addAlertToGroupByGroupId(groupId string, alert models.AlertCurEvent) {
	ec.Lock()
	defer ec.Unlock()

	if alert.IsRecovered {
		ec.preStoreRecoverAlertEvents[groupId] = append(ec.preStoreRecoverAlertEvents[groupId], alert)
	} else {
		ec.preStoreFiringAlertEvents[groupId] = append(ec.preStoreFiringAlertEvents[groupId], alert)
	}
}

// 以ruleName作为key添加到组
func (ec *Consume) addAlertToGroupByRuleId(alert models.AlertCurEvent) {
	ec.Lock()
	defer ec.Unlock()

	if alert.IsRecovered {
		ec.preStoreRecoverAlertEvents[alert.RuleId] = append(ec.preStoreRecoverAlertEvents[alert.RuleId], alert)
	} else {
		ec.preStoreFiringAlertEvents[alert.RuleId] = append(ec.preStoreFiringAlertEvents[alert.RuleId], alert)
	}
}

// hash
func (ec *Consume) calculateGroupHash(key, value string) string {
	return hash.Md5Hash([]byte(key + ":" + value))
}

// 推送告警
func (ec *Consume) handleAlert(alertMapping map[string][]models.AlertCurEvent) {
	curTime := time.Now().Unix()
	for key, alerts := range alertMapping {
		if strings.Contains(key, "_") {
			i := strings.Split(key, "_")
			key = i[1]
		}

		object := ec.ctx.DB.Rule().GetRuleObject(key)
		if object.RuleId == "" {
			return
		}

		if *object.AlarmAggregation {
			alerts = ec.groupAlert(curTime, alerts)
		}

		if len(alerts) <= 0 {
			return
		}

		for _, alert := range alerts {
			noticeId := process.GetNoticeGroupId(alert)

			r := models.NoticeQuery{
				TenantId: alert.TenantId,
				Uuid:     noticeId,
			}
			noticeData, _ := ec.ctx.DB.Notice().Get(r)
			alert.DutyUser = process.GetDutyUser(ec.ctx, noticeData)
			err := sender.Sender(ec.ctx, alert, noticeData) //发送告警
			if err != nil {
				global.Logger.Sugar().Errorf(err.Error())
				return
			}

			if !alert.IsRecovered {
				alert.LastSendTime = curTime
				ctx.Redis.Event().SetCache("Firing", alert, 0)
			}
		}
	}
}

// 聚合告警
func (ec *Consume) groupAlert(timeInt int64, alerts []models.AlertCurEvent) []models.AlertCurEvent {
	var (
		alertOne []models.AlertCurEvent
		content  string
	)

	if len(alerts) > 1 {
		content = fmt.Sprintf("聚合 %d 条告警\n", len(alerts))
	}

	for _, alert := range alerts {
		if !ec.isSilence(alert) {
			alertOne = []models.AlertCurEvent{alert}
			alertOne[0].Annotations += "\n" + content
		}

		if !alert.IsRecovered {
			alert.LastSendTime = timeInt
			ctx.Redis.Event().SetCache("Firing", alert, 0)
		}
	}

	return alertOne
}

// 判断是否静默
func (ec *Consume) isSilence(alert models.AlertCurEvent) bool {
	_, ok := ctx.Redis.Silence().GetCache(models.AlertSilenceQuery{
		TenantId:    alert.TenantId,
		Fingerprint: alert.Fingerprint,
	})

	if ok {
		return true
	} else {
		ttl, _ := ctx.Redis.Redis().TTL(alert.TenantId + ":" + models.SilenceCachePrefix + alert.Fingerprint).Result()
		// 如果剩余生存时间小于0，表示键已过期
		if ttl < 0 {
			// 过期后标记为1
			ctx.DB.DB().Model(models.AlertSilences{}).
				Where("fingerprint = ? and status = ?", alert.Fingerprint, 0).
				Update("status", 1)
		}
	}

	return false
}
