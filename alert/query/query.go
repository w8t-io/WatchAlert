package query

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sort"
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
		case "Loki":
			curKeys = rq.loki(dsId, rule)
		}

		// 处理恢复逻辑
		firingKeys := getFiringAlertCacheKeys(rule, dsId)
		recoverKeys = getSliceDifference(firingKeys, curKeys)
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
				event.SetFiringCache(0)
			}(key, curTime)
		}

	}

}

// Prometheus 数据源
func (rq *RuleQuery) prometheus(datasourceId string, rule models.AlertRule) []string {

	resQuery, _, err := client.NewPromClient(datasourceId).Query(rule.PrometheusConfig.PromQL)
	if err != nil {
		return nil
	}

	var curFiringKeys, curPendingKeys []string

	if resQuery == nil {
		go gcPendingCache(rule, datasourceId, curPendingKeys)
		return nil
	}

	for _, v := range resQuery {
		fingerprint := v.Labels.FastFingerprint().String()
		firingKey := rq.alertEvent.FiringAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
		pendingKey := rq.alertEvent.PendingAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
		curFiringKeys = append(curFiringKeys, firingKey)
		curPendingKeys = append(curPendingKeys, pendingKey)

		// handle series tags
		metricMap := make(map[string]interface{})
		for label, value := range v.Labels {
			metricMap[string(label)] = string(value)
		}
		metricMap["value"] = v.Value

		metricArr := labelMapToArr(metricMap)
		sort.Strings(metricArr)

		event := parserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint
		event.Metric = metricMap
		event.Annotations = cmd.ParserVariables(rule.Annotations, event.Metric)

		saveEventCache(event)
	}

	go gcPendingCache(rule, datasourceId, curPendingKeys)

	return curFiringKeys

}

// AliCloudSLS 数据源
func (rq *RuleQuery) aliCloudSLS(datasourceId string, rule models.AlertRule) []string {

	curAt := time.Now()
	startsAt := parserDuration(curAt, rule.AliCloudSLSConfig.LogScope, "m")
	args := client.AliCloudSlsQueryArgs{
		Project:  rule.AliCloudSLSConfig.Project,
		Logstore: rule.AliCloudSLSConfig.Logstore,
		StartsAt: int32(startsAt.Unix()),
		EndsAt:   int32(curAt.Unix()),
		Query:    rule.AliCloudSLSConfig.LogQL,
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

	bodyString, _ := json.Marshal(res.Body[0])
	// 标签，用于推送告警消息时 获取相关 label 信息
	metricMap := make(map[string]interface{})
	err = json.Unmarshal(bodyString, &metricMap)
	if err != nil {
		globals.Logger.Sugar().Errorf("解析 SLS Metric Label 失败, %s", err.Error())
	}

	// 删除多余 label
	delete(metricMap, "_image_name_")
	delete(metricMap, "__topic__")
	delete(metricMap, "_container_ip_")
	delete(metricMap, "_pod_uid_")
	delete(metricMap, "_source_")
	delete(metricMap, "_time_")
	delete(metricMap, "__time__")
	delete(metricMap, "__tag__:__pack_id__")
	var annotation string
	if metricMap["content"] != nil {
		annotation = metricMap["content"].(string)
		if isJSON(annotation) {
			// 将字符串解析为map类型
			var data map[string]interface{}
			err = json.Unmarshal([]byte(annotation), &data)
			if err != nil {
				globals.Logger.Sugar().Errorf("Error parsing JSON: %s", err.Error())
			} else {
				// 格式化JSON并输出
				formattedJson, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					globals.Logger.Sugar().Errorf("Error marshalling JSON: %s", err.Error())
				} else {
					annotation = string(formattedJson)
				}
			}
		}
	}
	delete(metricMap, "content")

	var curKeys []string
	h := md5.New()
	// 使用 label 进行 Hash 作为告警指纹，可以有效地作为恢复逻辑的判断条件。
	h.Write([]byte(cmd.JsonMarshal(metricMap)))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	key := rq.alertEvent.FiringAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
	curKeys = append(curKeys, key)

	event := func() {
		event := parserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint
		event.Annotations = annotation
		event.Metric = metricMap

		saveEventCache(event)
	}

	options := models.EvalCondition{
		/*
			触发告警的条件
			- 有数据 > number	// 有数据并大于多少条。
		*/
		Type:     rule.AliCloudSLSConfig.EvalCondition.Type,
		Operator: rule.AliCloudSLSConfig.EvalCondition.Operator,
		Value:    rule.AliCloudSLSConfig.EvalCondition.Value,
	}

	// 评估告警条件
	evalCondition(event, count, options)

	return curKeys

}

// Loki 数据源
func (rq *RuleQuery) loki(datasourceId string, rule models.AlertRule) []string {

	curAt := time.Now().UTC()
	startsAt := parserDuration(curAt, rule.LokiConfig.LogScope, "m")
	args := client.QueryOptions{
		Query:   rule.LokiConfig.LogQL,
		StartAt: startsAt.Format(time.RFC3339Nano),
		EndAt:   curAt.Format(time.RFC3339Nano),
	}

	res, err := client.NewLokiClient(datasourceId).QueryRange(args)
	if err != nil {
		globals.Logger.Sugar().Errorf("查询 Loki 日志失败 %s", err.Error())
		return nil
	}

	var curKeys []string

	for _, v := range res {

		count := len(v.Values)
		if count <= 0 {
			continue
		}

		// 使用 Loki 提供的 Stream label 进行 Hash 作为告警指纹.
		h := md5.New()
		streamString := cmd.JsonMarshal(v.Stream)
		h.Write([]byte(streamString))
		fingerprint := hex.EncodeToString(h.Sum(nil))
		key := rq.alertEvent.FiringAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
		curKeys = append(curKeys, key)

		// 标签，用于推送告警消息时 获取相关 label 信息
		metricMap := make(map[string]interface{})
		for label, value := range v.Stream {
			metricMap[label] = value
		}

		delete(metricMap, "stream")
		delete(metricMap, "filename")

		event := func() {
			event := parserDefaultEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = fingerprint
			event.Metric = metricMap

			var logValue string
			if v.Values[0] != nil {
				if v.Values[0].([]interface{}) != nil {
					logValue = v.Values[0].([]interface{})[1].(string)
				}
			}

			var logV10 client.LogValueV10
			err := json.Unmarshal([]byte(logValue), &logV10)
			if err != nil {
				event.Annotations = logValue
			} else {
				if logV10.Log == "" {
					return
				}
				event.Annotations = logV10.Log.(string)
			}

			saveEventCache(event)
		}

		options := models.EvalCondition{
			Type:     rule.LokiConfig.EvalCondition.Type,
			Operator: rule.LokiConfig.EvalCondition.Operator,
			Value:    rule.LokiConfig.EvalCondition.Value,
		}

		// 评估告警条件
		evalCondition(event, count, options)

	}

	return curKeys

}
