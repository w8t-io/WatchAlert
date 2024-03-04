package query

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/client"
	"watchAlert/utils/cmd"
)

type RuleQuery struct {
	alertEvent models.AlertCurEvent
}

func (rq *RuleQuery) Query(rule models.AlertRule) {

	var recoverKeys []string

	for _, dsId := range rule.DatasourceIdList {

		var curKeys []string
		switch rule.DatasourceType {
		case "Prometheus":
			curKeys = rq.prometheus(dsId, rule)
		case "AliCloudSLS":
			curKeys = rq.aliCloudSLS(dsId, rule)
		}

		// 处理恢复逻辑
		allKey := alertCacheKeys(rule, dsId)
		recoverKeys = getSliceDifference(allKey, curKeys)

		for _, key := range recoverKeys {
			curTime := time.Now().Unix()
			go func(key string, curTime int64) {
				event := rq.alertEvent.GetCache(key)
				if event.IsRecovered == true {
					return
				}
				event.IsRecovered = true
				event.RecoverTime = curTime
				event.LastSendTime = 0
				rq.alertEvent.SetCache(event, 0)
			}(key, curTime)
		}

	}

}

// Prometheus 数据源
func (rq *RuleQuery) prometheus(datasourceId string, rule models.AlertRule) []string {

	resQuery, _, err := client.NewPromClient(datasourceId).Query(rule.RuleConfigJson.PromQL)
	if err != nil {
		return nil
	}

	var curKeys []string

	if resQuery != nil {
		for _, v := range resQuery {
			fingerprint := v.Labels.FastFingerprint().String()
			key := rq.alertEvent.CurAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
			curKeys = append(curKeys, key)

			// handle series tags
			metricMap := make(map[string]interface{})
			for label, value := range v.Labels {
				metricMap[string(label)] = string(value)
			}
			metricMap["value"] = v.Value

			metricArr := labelMapToArr(metricMap)
			sort.Strings(metricArr)

			event := parserDefaultEvent(key, rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = fingerprint
			event.Metric = strings.Join(metricArr, ",,")
			event.MetricMap = metricMap
			event.Annotations = cmd.ParserVariables(rule.Annotations, event.MetricMap)

			saveEventCache(event)
		}
	}

	return curKeys

}

// AliCloudSLS 数据源
func (rq *RuleQuery) aliCloudSLS(datasourceId string, rule models.AlertRule) []string {

	curAt := time.Now()
	duration, err := time.ParseDuration(strconv.Itoa(rule.AliCloudSLSConfigJson.AliCloudQueryLogScope) + "m")
	if err != nil {
		globals.Logger.Sugar().Error("解析相对时间失败 ->", err.Error())
		return nil
	}

	startsAt := curAt.Add(-duration)
	args := client.AliCloudSlsQueryArgs{
		Project:  rule.AliCloudSLSConfigJson.AliCloudProject,
		Logstore: rule.AliCloudSLSConfigJson.AliCloudLogstore,
		StartsAt: int32(startsAt.Unix()),
		EndsAt:   int32(curAt.Unix()),
		Query:    rule.AliCloudSLSConfigJson.AliCloudQuerySQL,
	}

	res, err := client.NewAliCloudSlsClient(datasourceId).Query(args)
	if err != nil {
		globals.Logger.Sugar().Error("查询 AliCloudSls 日志失败 ->", err.Error())
		return nil
	}

	count := len(res.Body)
	if count <= 0 {
		return nil
	}

	var curKeys []string
	h := md5.New()
	// 使用 Query 查询条件进行 Hash 作为告警指纹，可以有效地作为恢复逻辑的判断条件。
	h.Write([]byte((*res.Headers["x-log-where-query"])))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	key := rq.alertEvent.CurAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
	curKeys = append(curKeys, key)

	/*
		触发告警的条件
		- 有数据 > number	// 有数据并大于多少条。
	*/
	t := rule.AliCloudSLSConfigJson.AliCloudEvalConditionJson.Type
	operator := rule.AliCloudSLSConfigJson.AliCloudEvalConditionJson.Operator
	value := rule.AliCloudSLSConfigJson.AliCloudEvalConditionJson.Value

	event := func() {
		event := parserDefaultEvent(key, rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint
		bodyString, _ := json.Marshal(res.Body[0])
		event.Annotations = string(bodyString)
		saveEventCache(event)
	}

	switch t {
	case "count":
		switch operator {
		case ">":
			if count > value {
				event()
			}
		case ">=":
			if count >= value {
				event()
			}
		case "<":
			if count < value {
				event()
			}
		case "<=":
			if count <= value {
				event()
			}
		case "==":
			if count == value {
				event()
			}
		case "!=":
			if count != value {
				event()
			}
		default:
			globals.Logger.Sugar().Error("无效的评估条件", t, operator, value)
		}
	default:
		globals.Logger.Sugar().Error("无效的评估类型", t)
	}

	return curKeys

}
