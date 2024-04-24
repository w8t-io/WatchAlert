package query

import (
	"time"
	"watchAlert/alert/queue"
	"watchAlert/models"
	"watchAlert/public/client"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
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
		case "Jaeger":
			rq.jaeger(dsId, rule)
		}
	}

}

func (rq *RuleQuery) alertRecover(rule models.AlertRule, curKeys []string) {
	firingKeys := rule.GetFiringAlertCacheKeys()
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

		event.IsRecovered = true
		event.RecoverTime = curTime
		event.LastSendTime = 0
		event.SetFiringCache(0)

		// 触发恢复删除带恢复中的 key
		delete(queue.RecoverWaitMap, key)
	}
}

// Prometheus 数据源
func (rq *RuleQuery) prometheus(datasourceId string, rule models.AlertRule) {
	var curFiringKeys, curPendingKeys []string
	defer func() {
		go gcPendingCache(rule, datasourceId, curPendingKeys)
		rq.alertRecover(rule, curFiringKeys)
		go gcRecoverWaitCache(rule, curFiringKeys)
	}()

	resQuery, _, err := client.NewPromClient(rule.TenantId, datasourceId).Query(rule.PrometheusConfig.PromQL)
	if err != nil {
		return
	}

	if resQuery == nil {
		return
	}

	for _, v := range resQuery {
		event := parserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = v.GetFingerprint()
		event.Metric = v.GetMetric()
		event.Annotations = cmd.ParserVariables(rule.Annotations, event.Metric)

		firingKey := event.GetFiringAlertCacheKey()
		pendingKey := event.GetPendingAlertCacheKey()
		curFiringKeys = append(curFiringKeys, firingKey)
		curPendingKeys = append(curPendingKeys, pendingKey)

		saveEventCache(event)
	}

}

// AliCloudSLS 数据源
func (rq *RuleQuery) aliCloudSLS(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, curKeys)
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

	res, err := client.NewAliCloudSlsClient(rule.TenantId, datasourceId).Query(args)
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

		event := func() {
			event := parserDefaultEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = body.GetFingerprint()
			event.Annotations = body.GetAnnotations()
			event.Metric = body.GetMetric()

			key := event.GetFiringAlertCacheKey()
			curKeys = append(curKeys, key)

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
		rq.alertRecover(rule, curKeys)
		go gcRecoverWaitCache(rule, curKeys)
	}()

	curAt := time.Now().UTC()
	startsAt := parserDuration(curAt, rule.LokiConfig.LogScope, "m")
	args := client.QueryOptions{
		Query:   rule.LokiConfig.LogQL,
		StartAt: startsAt.Format(time.RFC3339Nano),
		EndAt:   curAt.Format(time.RFC3339Nano),
	}

	res, err := client.NewLokiClient(rule.TenantId, datasourceId).QueryRange(args)
	if err != nil {
		globals.Logger.Sugar().Errorf("查询 Loki 日志失败 %s", err.Error())
		return
	}

	for _, v := range res {

		count := len(v.Values)
		if count <= 0 {
			continue
		}

		event := func() {
			event := parserDefaultEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Annotations = v.GetAnnotations().(string)

			key := event.GetPendingAlertCacheKey()
			curKeys = append(curKeys, key)

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

// Jaeger 数据源
func (rq *RuleQuery) jaeger(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, curKeys)
		go gcRecoverWaitCache(rule, curKeys)
	}()

	curAt := time.Now().UTC()
	startsAt := parserDuration(curAt, rule.JaegerConfig.Scope, "m")

	rule.DatasourceType = "Jaeger"
	rule.DatasourceIdList = []string{"jaeger"}

	opt := client.JaegerQueryOptions{
		Tags:    rule.JaegerConfig.Tags,
		Service: rule.JaegerConfig.Service,
		StartAt: startsAt.UnixMicro(),
		EndAt:   curAt.UnixMicro(),
	}

	res := client.NewJaegerClient(datasourceId).JaegerQuery(opt)
	if res.Data == nil {
		return
	}

	for _, v := range res.Data {
		event := parserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = v.GetFingerprint()
		event.Metric = v.GetMetric(rule)
		event.Annotations = v.GetAnnotations(rule)

		key := rq.alertEvent.GetFiringAlertCacheKey()
		curKeys = append(curKeys, key)

		saveEventCache(event)
	}

}
