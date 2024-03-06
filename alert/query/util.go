package query

import (
	"fmt"
	"sort"
	"strconv"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	"watchAlert/models"
)

func alertCacheKeys(rule models.AlertRule, dsId string) []string {

	var alertCurEvent models.AlertCurEvent
	// 获取所有keys
	keyPrefix := alertCurEvent.CurAlertCacheKey(rule.RuleId, dsId, "*")
	keys, _ := globals.RedisCli.Keys(keyPrefix).Result()

	return keys

}

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

func parserDefaultEvent(key string, rule models.AlertRule) models.AlertCurEvent {

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
	event.FirstTriggerTime = event.GetFirstTime(key)
	event.LastEvalTime = event.GetLastEvalTime(key)
	event.LastSendTime = event.GetLastSendTime(key)

	return event

}

func saveEventCache(event models.AlertCurEvent) {

	event.SetCache(event, 0)
	err := repo.DBCli.Create(models.AlertCurEvent{}, &event)
	if err != nil {
		globals.Logger.Sugar().Errorf("Failed inserting AlertCurEvent into the database: %s", err)
		return
	}

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