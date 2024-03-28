package query

import (
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
		rq.alertRecover(rule, dsId, curKeys)
	}

}

func (rq *RuleQuery) alertRecover(rule models.AlertRule, dsId string, curKeys []string) {
	var recoverKeys []string
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
		fingerprint := v.GetFingerprint()
		firingKey := rq.alertEvent.FiringAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
		pendingKey := rq.alertEvent.PendingAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
		curFiringKeys = append(curFiringKeys, firingKey)
		curPendingKeys = append(curPendingKeys, pendingKey)

		event := parserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint
		event.Metric = v.GetMetric()
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

	body := client.GetSLSBodyData(res)

	var curKeys []string
	fingerprint := body.GetFingerprint()
	key := rq.alertEvent.FiringAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
	curKeys = append(curKeys, key)

	event := func() {
		event := parserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = fingerprint
		event.Annotations = body.GetAnnotations()
		event.Metric = body.GetMetric()

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

		fingerprint := v.GetFingerprint()
		key := rq.alertEvent.FiringAlertCacheKey(rule.RuleId, datasourceId, fingerprint)
		curKeys = append(curKeys, key)

		event := func() {
			event := parserDefaultEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = fingerprint
			event.Metric = v.GetMetric()
			event.Annotations = v.GetAnnotations().(string)

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
