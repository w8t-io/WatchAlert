package eval

import (
	"fmt"
	"sync"
	"time"
	"watchAlert/alert/notice"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/services"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/hash"
)

type EvalConsume struct {
	sync.RWMutex
	models.AlertCurEvent
	// 从 Redis 中读取当前告警事件提到内存做处理.
	alertsMap map[string][]models.AlertCurEvent
	// 告警分组
	preStoreAlertGroup map[string][]models.AlertCurEvent
	Timing             map[string]int
}

type InterEvalConsume interface {
	Run()
}

func NewInterEvalConsumeWork() InterEvalConsume {

	return &EvalConsume{
		alertsMap:          make(map[string][]models.AlertCurEvent),
		preStoreAlertGroup: make(map[string][]models.AlertCurEvent),
		Timing:             make(map[string]int),
	}

}

// Run 启动告警消费进程
func (ec *EvalConsume) Run() {

	action := func() {
		alertsCurEventKeys := ec.getRedisKeys()
		for _, key := range alertsCurEventKeys {
			alert := ec.GetCache(key)
			// 过滤空指纹告警
			if alert.Fingerprint == "" {
				continue
			}
			ec.addAlertToRuleIdMap(alert)
		}

		for key, alerts := range ec.alertsMap {
			if len(alerts) == 0 {
				continue
			}

			// 计算告警组的等待时间
			var waitTime int
			alert := ec.GetCache(key)
			if alert.LastSendTime == 0 {
				// 如果是初次告警, 那么等当前告警组时间到达 groupWait 的时间则推送告警
				waitTime = globals.Config.Server.GroupWait
			} else {
				// 当前告警组时间到达 groupInterval 的时间则推送告警
				waitTime = globals.Config.Server.GroupInterval
			}
			if ec.Timing[key] >= waitTime {
				curEvent := ec.filterAlerts(ec.alertsMap[key])
				ec.fireAlertEvent(curEvent)
				// 执行一波后 必须重新清空alerts组中的数据。
				ec.clear(key)
			}
			ec.Timing[key]++
		}
	}

	ticker := time.Tick(time.Second)

	go func() {
		for range ticker {
			action()
		}
	}()

}

func (ec *EvalConsume) addAlertToRuleIdMap(alert models.AlertCurEvent) {
	ec.Lock()
	ec.alertsMap[alert.RuleId] = append(ec.alertsMap[alert.RuleId], alert)
	ec.Unlock()
}

func (ec *EvalConsume) clear(ruleId string) {

	for key := range ec.alertsMap {
		delete(ec.alertsMap, key)
	}
	for key := range ec.preStoreAlertGroup {
		delete(ec.preStoreAlertGroup, key)
	}
	ec.Timing[ruleId] = 0

}

// 获取缓存所有Firing的Keys
func (ec *EvalConsume) getRedisKeys() []string {
	var keys []string
	cursor := uint64(0)
	pattern := "*" + ":" + models.FiringAlertCachePrefix + "*"
	// 每次获取的键数量
	count := int64(100)

	for {
		var curKeys []string
		var err error

		curKeys, cursor, err = globals.RedisCli.Scan(cursor, pattern, count).Result()
		if err != nil {
			break
		}

		keys = append(keys, curKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys
}

// 过滤告警
func (ec *EvalConsume) filterAlerts(alerts []models.AlertCurEvent) map[string][]models.AlertCurEvent {

	var newAlertsMap = make(map[string][]models.AlertCurEvent)

	// 根据相同指纹进行去重
	newAlert := ec.removeDuplicates(alerts)
	// 将通过指纹去重后以Fingerprint为Key的Map转换成以原来RuleName为Key的Map (同一告警类型聚合)
	for _, alert := range newAlert {
		// 重复通知，如果是初次推送不用进一步判断。
		if !alert.IsRecovered {
			if alert.LastSendTime == 0 || alert.LastEvalTime >= alert.LastSendTime+alert.RepeatNoticeInterval*60 {
				newAlertsMap[alert.RuleName] = append(newAlertsMap[alert.RuleName], alert)
			}
		}
		if alert.IsRecovered {
			newAlertsMap[alert.RuleName] = append(newAlertsMap[alert.RuleName], alert)
		}
	}

	return newAlertsMap

}

// 指纹去重
func (ec *EvalConsume) removeDuplicates(alerts []models.AlertCurEvent) []models.AlertCurEvent {
	/*
		alert中有不重复字段，last_eval_time。
	*/

	latestAlert := make(map[string]models.AlertCurEvent)
	var newAlerts []models.AlertCurEvent

	for _, alert := range alerts {
		// 以最新为准
		latestAlert[alert.Fingerprint] = alert
	}

	for _, alert := range latestAlert {
		newAlerts = append(newAlerts, alert)
	}

	return newAlerts
}

// 触发告警通知
func (ec *EvalConsume) fireAlertEvent(alertsMap map[string][]models.AlertCurEvent) {
	var wg sync.WaitGroup

	for _, alerts := range alertsMap {
		for _, alert := range alerts {
			wg.Add(1)
			go func(alert models.AlertCurEvent) {
				defer wg.Done()
				ec.addAlertToGroup(alert)
				if alert.IsRecovered {
					ec.removeAlertFromCache(alert)
					err := ec.RecordAlertHisEvent(alert)
					if err != nil {
						globals.Logger.Sugar().Error(err.Error())
						return
					}
				}
			}(alert)
		}
	}

	wg.Wait()

	for _, alerts := range ec.preStoreAlertGroup {
		ec.handleAlert(alerts)
	}
}

// 删除缓存
func (ec *EvalConsume) removeAlertFromCache(alert models.AlertCurEvent) {
	key := alert.GetFiringAlertCacheKey()
	ec.DelCache(key)
}

// 添加告警到组
func (ec *EvalConsume) addAlertToGroup(alert models.AlertCurEvent) {
	// 如果没有定义通知组，则直接添加到 ruleId 组中
	if alert.NoticeGroupList == nil || len(alert.NoticeGroupList) == 0 {
		ec.addAlertToGroupByRuleId(alert)
		return
	}

	// 遍历所有的 Metric
	matched := false
	for key, value := range alert.Metric {
		// 遍历所有的通知组
		for _, noticeGroup := range alert.NoticeGroupList {
			// 如果当前 Metric 的 key 和 value 与通知组中的相匹配
			if noticeGroup["key"] == key && noticeGroup["value"] == value.(string) {
				// 计算分组的 ID 并添加警报到对应的组
				groupId := ec.calculateGroupHash(key, value.(string))
				ec.addAlertToGroupById(groupId, alert)
				matched = true
				break // 找到匹配的组后，跳出内层循环
			}
		}
		if matched {
			break // 找到匹配的组后，跳出外层循环
		}
	}

	// 如果没有找到任何匹配的组，则添加到 ruleId 组中
	if !matched {
		ec.addAlertToGroupByRuleId(alert)
	}
}

// 以Id作为key添加到组
func (ec *EvalConsume) addAlertToGroupById(groupId string, alert models.AlertCurEvent) {
	ec.Lock()
	defer ec.Unlock()

	// 将告警和恢复消息再分组
	if alert.IsRecovered {
		groupId = "recovered-" + groupId
	}

	ec.preStoreAlertGroup[groupId] = append(ec.preStoreAlertGroup[groupId], alert)
}

// 以ruleName作为key添加到组
func (ec *EvalConsume) addAlertToGroupByRuleId(alert models.AlertCurEvent) {
	ec.Lock()
	defer ec.Unlock()

	// 将告警和恢复消息再分组
	if alert.IsRecovered {
		alert.RuleId = "recovered-" + alert.RuleId
	}
	ec.preStoreAlertGroup[alert.RuleId] = append(ec.preStoreAlertGroup[alert.RuleId], alert)
}

// hash
func (ec *EvalConsume) calculateGroupHash(key, value string) string {
	return hash.Md5Hash([]byte(key + ":" + value))
}

// 推送告警
func (ec *EvalConsume) handleAlert(alerts []models.AlertCurEvent) {

	if alerts == nil {
		return
	}

	var (
		content  string
		alertOne models.AlertCurEvent
		curTime  = time.Now().Unix()
	)

	if len(alerts) > 1 {
		content = fmt.Sprintf("聚合 %d 条告警\n", len(alerts))
	}

	var wg sync.WaitGroup
	for _, alert := range alerts {
		if !alert.IsRecovered {
			wg.Add(1)
			go func(alert models.AlertCurEvent) {
				defer wg.Done()
				alert.LastSendTime = curTime
				alert.SetFiringCache(0)
			}(alert)
		}
	}
	wg.Wait()

	// 聚合, 每组告警取第一位的告警数据
	alertOne = alerts[0]
	alertOne.Annotations += "\n" + content

	noticeId := ec.getNoticeGroupId(alertOne)

	noticeData := services.NewInterAlertNoticeService().GetNoticeObject(alertOne.TenantId, noticeId)

	var tmpl notice.Template
	notice.NewEntryNotice(&tmpl, alertOne, noticeData)

}

// 获取告警分组的通知ID
func (ec *EvalConsume) getNoticeGroupId(alert models.AlertCurEvent) string {

	if len(alert.NoticeGroupList) != 0 {
		var noticeGroup []map[string]string
		for _, v := range alert.NoticeGroupList {
			noticeGroup = append(noticeGroup, map[string]string{
				v["key"]:   v["value"],
				"noticeId": v["noticeId"],
			})
		}

		// 从Metric中获取Key/Value
		for metricKey, metricValue := range alert.Metric {
			// 如果配置分组的Key/Value 和 Metric中的Key/Value 一致，则使用分组的 noticeId，匹配不到则用默认的。
			for _, noticeInfo := range noticeGroup {
				value, ok := noticeInfo[metricKey]
				if ok && metricValue == value {
					noticeId := noticeInfo["noticeId"]
					return noticeId
				}
			}
		}
	}

	return alert.NoticeId

}

// RecordAlertHisEvent 记录历史告警
func (ec *EvalConsume) RecordAlertHisEvent(alert models.AlertCurEvent) error {

	hisData := models.AlertHisEvent{
		TenantID:         alert.TenantId,
		DatasourceType:   alert.DatasourceType,
		DatasourceId:     alert.DatasourceId,
		Fingerprint:      alert.Fingerprint,
		RuleId:           alert.RuleId,
		RuleName:         alert.RuleName,
		Severity:         alert.Severity,
		Metric:           alert.Metric,
		EvalInterval:     alert.EvalInterval,
		Annotations:      alert.Annotations,
		IsRecovered:      true,
		FirstTriggerTime: alert.FirstTriggerTime,
		LastEvalTime:     alert.LastEvalTime,
		LastSendTime:     alert.LastSendTime,
		RecoverTime:      alert.RecoverTime,
	}

	err := repo.DBCli.Create(models.AlertHisEvent{}, &hisData)
	if err != nil {
		return fmt.Errorf("RecordAlertHisEvent -> %s", err)
	}

	return nil

}
