package query

import (
	"fmt"
	"sort"
	"strconv"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
)

// 获取缓存中与 Rule 相关所有的Firing key
func getFiringAlertCacheKeys(rule models.AlertRule, dsId string) []string {
	var alertCurEvent models.AlertCurEvent
	keyPrefix := alertCurEvent.FiringAlertCacheKey(rule.RuleId, dsId, "*")
	keys, _ := globals.RedisCli.Keys(keyPrefix).Result()
	return keys
}

// 获取缓存中与 Rule 相关所有的Pending key
func getPendingAlertCacheKeys(rule models.AlertRule, dsId string) []string {
	var alertCurEvent models.AlertCurEvent
	keyPrefix := alertCurEvent.PendingAlertCacheKey(rule.RuleId, dsId, "*")
	keys, _ := globals.RedisCli.Keys(keyPrefix).Result()
	return keys
}

// 获取差异key. 当slice1中存在, slice2不存在则标记为可恢复告警
func getSliceDifference(slice1 []string, slice2 []string) []string {
	difference := []string{}

	// 遍历缓存
	for _, item1 := range slice1 {
		found := false
		// 遍历当前key
		for _, item2 := range slice2 {
			if item1 == item2 {
				found = true
				break
			}
		}
		// 添加到差异切片中
		if !found {
			difference = append(difference, item1)
		}
	}

	return difference
}

func labelMapToArr(m map[string]interface{}) []string {
	numLabels := len(m)

	labelStrings := make([]string, 0, numLabels)
	for label, value := range m {
		labelStrings = append(labelStrings, fmt.Sprintf("%s=%s", label, value))
	}

	if numLabels > 1 {
		sort.Strings(labelStrings)
	}

	return labelStrings
}

func parserDefaultEvent(rule models.AlertRule) models.AlertCurEvent {

	event := models.AlertCurEvent{
		DatasourceType:       rule.DatasourceType,
		RuleId:               rule.RuleId,
		RuleName:             rule.RuleName,
		Severity:             rule.Severity,
		Labels:               rule.Labels,
		EvalInterval:         rule.EvalInterval,
		ForDuration:          rule.ForDuration,
		NoticeId:             rule.NoticeId,
		NoticeGroupList:      rule.NoticeGroupList,
		IsRecovered:          false,
		RepeatNoticeInterval: rule.RepeatNoticeInterval,
		DutyUser:             "暂无", // 默认暂无值班人员, 渲染模版时会实际判断 Notice 是否存在值班人员
	}

	return event

}

func saveEventCache(event models.AlertCurEvent) {

	firingKey := event.FiringAlertCacheKey(event.RuleId, event.DatasourceId, event.Fingerprint)
	pendingKey := event.PendingAlertCacheKey(event.RuleId, event.DatasourceId, event.Fingerprint)

	// 判断改事件是否是Firing状态, 如果不是Firing状态 则标记Pending状态
	resFiring := event.GetCache(firingKey)
	if resFiring.Fingerprint != "" {
		event.FirstTriggerTime = resFiring.FirstTriggerTime
		event.LastEvalTime = event.GetLastEvalTime(firingKey)
		event.LastSendTime = resFiring.LastSendTime
	} else {
		event.FirstTriggerTime = event.GetFirstTime(pendingKey)
		event.LastEvalTime = event.GetLastEvalTime(pendingKey)
		event.LastSendTime = event.GetLastSendTime(pendingKey)
		event.SetPendingCache(0)
	}

	// 持续时间
	if event.LastEvalTime-event.FirstTriggerTime < event.ForDuration {
		return
	}

	event.SetFiringCache(0)
	event.DelCache(pendingKey)

}

// 获取时间区间的开始时间
func parserDuration(curTime time.Time, logScope int, timeType string) time.Time {

	duration, err := time.ParseDuration(strconv.Itoa(logScope) + timeType)
	if err != nil {
		globals.Logger.Sugar().Error("解析相对时间失败 ->", err.Error())
		return time.Time{}
	}
	startsAt := curTime.Add(-duration)

	return startsAt

}

// 评估告警条件
func evalCondition(f func(), count int, ec models.EvalCondition) {

	switch ec.Type {
	case "count":
		switch ec.Operator {
		case ">":
			if count > ec.Value {
				f()
			}
		case ">=":
			if count >= ec.Value {
				f()
			}
		case "<":
			if count < ec.Value {
				f()
			}
		case "<=":
			if count <= ec.Value {
				f()
			}
		case "==":
			if count == ec.Value {
				f()
			}
		case "!=":
			if count != ec.Value {
				f()
			}
		default:
			globals.Logger.Sugar().Error("无效的评估条件", ec.Type, ec.Operator, ec.Value)
		}
	default:
		globals.Logger.Sugar().Error("无效的评估类型", ec.Type)
	}

}

/*
	清理 Pending 数据的缓存.
	场景: 第一次查询到有异常的指标会写入 Pending 缓存, 当该指标持续 Pending 到达持续时间后才会写入 Firing 缓存,
	那么未到达持续时间并且该指标恢复正常, 那么就需要清理该指标的 Pending 数据.
*/
func gcPendingCache(rule models.AlertRule, dsId string, curKeys []string) {
	var ae models.AlertCurEvent
	pendingKeys := getPendingAlertCacheKeys(rule, dsId)
	gcPendingKeys := getSliceDifference(pendingKeys, curKeys)
	for _, key := range gcPendingKeys {
		ae.DelCache(key)
	}
}
