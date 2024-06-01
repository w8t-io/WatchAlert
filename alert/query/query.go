package query

import (
	"regexp"
	"strconv"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/ctx"
)

type RuleQuery struct {
	alertEvent models.AlertCurEvent
	ctx        *ctx.Context
}

func (rq *RuleQuery) Query(ctx *ctx.Context, rule models.AlertRule) {
	rq.ctx = ctx

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
	firingKeys, err := rq.ctx.Redis.Rule().GetAlertFiringCacheKeys(models.AlertRuleQuery{
		TenantId:         rule.TenantId,
		RuleId:           rule.RuleId,
		DatasourceIdList: rule.DatasourceIdList,
	})
	if err != nil {
		return
	}
	// 获取已恢复告警的keys
	recoverKeys := process.GetSliceDifference(firingKeys, curKeys)
	if recoverKeys == nil {
		return
	}

	curTime := time.Now().Unix()
	for _, key := range recoverKeys {
		event := rq.ctx.Redis.Event().GetCache(key)
		if event.IsRecovered == true {
			return
		}

		if _, exists := queue.RecoverWaitMap[key]; !exists {
			// 如果没有，则记录当前时间
			queue.RecoverWaitMap[key] = curTime
			continue
		}

		// 判断是否在等待时间范围内
		rt := time.Unix(queue.RecoverWaitMap[key], 0).Add(time.Minute * time.Duration(global.Config.Server.RecoverWait)).Unix()
		if rt > curTime {
			continue
		}

		event.IsRecovered = true
		event.RecoverTime = curTime
		event.LastSendTime = 0

		rq.ctx.Redis.Event().SetCache("Firing", event, 0)

		// 触发恢复删除带恢复中的 key
		delete(queue.RecoverWaitMap, key)
	}
}

// Prometheus 数据源
func (rq *RuleQuery) prometheus(datasourceId string, rule models.AlertRule) {
	var (
		curFiringKeys  = &[]string{}
		curPendingKeys = &[]string{}
	)

	defer func() {
		go process.GcPendingCache(rq.ctx, rule, *curPendingKeys)
		rq.alertRecover(rule, *curFiringKeys)
		go process.GcRecoverWaitCache(rule, *curFiringKeys)
	}()

	r := models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
		Type:     "Prometheus",
	}
	datasourceInfo, err := rq.ctx.DB.Datasource().Get(r)
	if err != nil {
		return
	}

	resQuery, _, err := client.NewPromClient(datasourceInfo).Query(rule.PrometheusConfig.PromQL)

	if err != nil {
		return
	}

	if resQuery == nil {
		return
	}

	for _, v := range resQuery {
		for _, ruleExpr := range rule.PrometheusConfig.Rules {
			re := regexp.MustCompile(`([^\d]+)(\d+)`)
			matches := re.FindStringSubmatch(ruleExpr.Expr)
			t, _ := strconv.ParseFloat(matches[2], 64)
			process.CalIndicatorValue(rq.ctx, matches[1], t, rule, v, datasourceId, curFiringKeys, curPendingKeys, ruleExpr.Severity)
		}
	}

}

// AliCloudSLS 数据源
func (rq *RuleQuery) aliCloudSLS(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	curAt := time.Now()
	startsAt := process.ParserDuration(curAt, rule.AliCloudSLSConfig.LogScope, "m")
	args := client.AliCloudSlsQueryArgs{
		Project:  rule.AliCloudSLSConfig.Project,
		Logstore: rule.AliCloudSLSConfig.Logstore,
		StartsAt: int32(startsAt.Unix()),
		EndsAt:   int32(curAt.Unix()),
		Query:    rule.AliCloudSLSConfig.LogQL,
	}

	datasourceInfo, err := rq.ctx.DB.Datasource().Get(models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
	})
	if err != nil {
		return
	}

	res, err := client.NewAliCloudSlsClient(datasourceInfo).Query(args)
	if err != nil {
		global.Logger.Sugar().Error("查询 AliCloudSls 日志失败 ->", err.Error())
		return
	}

	count := len(res.Body)
	if count <= 0 {
		return
	}

	bodyList := client.GetSLSBodyData(res)

	for _, body := range bodyList.MetricList {

		event := func() {
			event := process.ParserDefaultEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = body.GetFingerprint()
			event.Annotations = body.GetAnnotations()
			event.Metric = body.GetMetric()

			key := event.GetFiringAlertCacheKey()
			curKeys = append(curKeys, key)

			ok := rq.ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				process.SaveEventCache(rq.ctx, event)
			}
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
		process.EvalCondition(event, count, options)
	}

}

// Loki 数据源
func (rq *RuleQuery) loki(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	curAt := time.Now().UTC()
	startsAt := process.ParserDuration(curAt, rule.LokiConfig.LogScope, "m")
	args := client.QueryOptions{
		Query:   rule.LokiConfig.LogQL,
		StartAt: startsAt.Format(time.RFC3339Nano),
		EndAt:   curAt.Format(time.RFC3339Nano),
	}

	datasourceInfo, err := rq.ctx.DB.Datasource().Get(models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
	})
	if err != nil {
		return
	}

	res, err := client.NewLokiClient(datasourceInfo).QueryRange(args)
	if err != nil {
		global.Logger.Sugar().Errorf("查询 Loki 日志失败 %s", err.Error())
		return
	}

	for _, v := range res {

		count := len(v.Values)
		if count <= 0 {
			continue
		}

		event := func() {
			event := process.ParserDefaultEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Annotations = v.GetAnnotations().(string)

			key := event.GetPendingAlertCacheKey()
			curKeys = append(curKeys, key)

			ok := rq.ctx.DB.Rule().GetRuleIsExist(event.RuleId)
			if ok {
				process.SaveEventCache(rq.ctx, event)
			}
		}

		options := models.EvalCondition{
			Type:     rule.LokiConfig.EvalCondition.Type,
			Operator: rule.LokiConfig.EvalCondition.Operator,
			Value:    rule.LokiConfig.EvalCondition.Value,
		}

		// 评估告警条件
		process.EvalCondition(event, count, options)

	}

}

// Jaeger 数据源
func (rq *RuleQuery) jaeger(datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		rq.alertRecover(rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	curAt := time.Now().UTC()
	startsAt := process.ParserDuration(curAt, rule.JaegerConfig.Scope, "m")

	rule.DatasourceType = "Jaeger"
	rule.DatasourceIdList = []string{"jaeger"}

	opt := client.JaegerQueryOptions{
		Tags:    rule.JaegerConfig.Tags,
		Service: rule.JaegerConfig.Service,
		StartAt: startsAt.UnixMicro(),
		EndAt:   curAt.UnixMicro(),
	}

	datasourceInfo, err := rq.ctx.DB.Datasource().Get(models.DatasourceQuery{
		Id: datasourceId,
	})
	if err != nil {
		return
	}

	res := client.NewJaegerClient(datasourceInfo).JaegerQuery(opt)
	if res.Data == nil {
		return
	}

	for _, v := range res.Data {
		event := process.ParserDefaultEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = v.GetFingerprint()
		event.Metric = v.GetMetric(rule)
		event.Annotations = v.GetAnnotations(rule)

		key := rq.alertEvent.GetFiringAlertCacheKey()
		curKeys = append(curKeys, key)

		ok := rq.ctx.DB.Rule().GetRuleIsExist(event.RuleId)
		if ok {
			process.SaveEventCache(rq.ctx, event)
		}
	}

}
