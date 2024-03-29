package query

import (
	"time"
	"watchAlert/alert/queue"
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
		switch rule.DatasourceType {
		case "Prometheus":
			rq.prometheus(dsId, rule)
		case "AliCloudSLS":
			rq.aliCloudSLS(dsId, rule)
		case "Loki":
			rq.loki(dsId, rule)
		}
	}

}

func (rq *RuleQuery) alertRecover(rule models.AlertRule, dsId string, curKeys []string) {
	firingKeys := getFiringAlertCacheKeys(rule, dsId)
	// 获取已恢复告警的keys
	recoverKeys := getSliceDifference(firingKeys, curKeys)
	if recoverKeys == nil {
		return
	}

	curTime := time.Now().Unix()
	for _, key := range recoverKeys {
		event := rq.alertEvent.GetCache(key)
		if event.IsRecovered == true {
			return
		}

		if _, exists := queue.RecoverWaitMap[key]; !exists {
			// 如果没有，则记录当前时间
			queue.RecoverWaitMap[key] = curTime
			continue
		}

		// 判断是否在等待时间范围内
		rt := time.Unix(queue.RecoverWaitMap[key], 0).Add(time.Minute * time.Duration(globals.Config.Server.RecoverWait)).Unix()
		if rt > curTime {
			continue
		}

		go func(key string, curTime int64, event models.AlertCurEvent) {
			event.IsRecovered = true
			event.RecoverTime = curTime
			event.LastSendTime = 0
			event.SetFiringCache(0)
		}(key, curTime, event)

		// 触发恢复删除带恢复中的 key
		delete(queue.RecoverWaitMap, key)
	}
}

// Prometheus 数据源
func (rq *RuleQuery) prometheus(datasourceId string, rule models.AlertRule) {
	var curFiringKeys, curPendingKeys []string
	defer func() {
		go gcPendingCache(rule, datasourceId, curPendingKeys)
		rq.alertRecover(rule, datasourceId, curFiringKeys)
		go gcRecoverWaitCache(rule, curFiringKeys)
	}()

	resQuery, _, err := client.NewPromClient(datasourceId).Query(rule.PrometheusConfig.PromQL)
	if err != nil {
		return
	}

	if resQuery == nil {
		return
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

}

// AliCloudSLS 数据源
func (rq *RuleQuery) aliCloudSLS(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, datasourceId, curKeys)
		go gcRecoverWaitCache(rule, curKeys)
	}()

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
		return
	}

	count := len(res.Body)
	if count <= 0 {
		return
	}

	bodyList := client.GetSLSBodyData(res)

	for _, body := range bodyList.MetricList {
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
	}

}

// Loki 数据源
func (rq *RuleQuery) loki(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, datasourceId, curKeys)
		go gcRecoverWaitCache(rule, curKeys)
	}()

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
		return
	}

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

}
