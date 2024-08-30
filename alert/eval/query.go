package eval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"watchAlert/alert/process"
	"watchAlert/alert/queue"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/community/aws/cloudwatch"
	"watchAlert/pkg/community/aws/cloudwatch/types"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
)

func alertRecover(ctx *ctx.Context, rule models.AlertRule, curKeys []string) {
	firingKeys, err := ctx.Redis.Rule().GetAlertFiringCacheKeys(models.AlertRuleQuery{
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
		event := ctx.Redis.Event().GetCache(key)
		if event.IsRecovered == true {
			return
		}

		if _, exists := queue.RecoverWaitMap[key]; !exists {
			// 如果没有，则记录当前时间
			queue.RecoverWaitMap[key] = curTime
			continue
		}

		// 判断是否在等待时间范围内
		rt := time.Unix(queue.RecoverWaitMap[key], 0).Add(time.Minute * time.Duration(global.Config.Server.AlarmConfig.RecoverWait)).Unix()
		if rt > curTime {
			continue
		}

		event.IsRecovered = true
		event.RecoverTime = curTime
		event.LastSendTime = 0

		ctx.Redis.Event().SetCache("Firing", event, 0)

		// 触发恢复删除带恢复中的 key
		delete(queue.RecoverWaitMap, key)
	}
}

// Prometheus 数据源
func prometheus(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var (
		curFiringKeys  []string
		curPendingKeys []string
	)

	defer func() {
		go process.GcPendingCache(ctx, rule, curPendingKeys)
		alertRecover(ctx, rule, curFiringKeys)
		go process.GcRecoverWaitCache(rule, curFiringKeys)
	}()

	r := models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
		Type:     "Prometheus",
	}
	datasourceInfo, err := ctx.DB.Datasource().Get(r)
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

			f := func() models.AlertCurEvent {
				event := process.BuildEvent(rule)
				event.DatasourceId = datasourceId
				event.Fingerprint = v.GetFingerprint()
				event.Metric = v.GetMetric()
				event.Metric["severity"] = ruleExpr.Severity
				event.Severity = ruleExpr.Severity
				event.Annotations = cmd.ParserVariables(rule.PrometheusConfig.Annotations, event.Metric)

				firingKey := event.GetFiringAlertCacheKey()
				pendingKey := event.GetPendingAlertCacheKey()

				curFiringKeys = append(curFiringKeys, firingKey)
				curPendingKeys = append(curPendingKeys, pendingKey)

				return event
			}

			option := models.EvalCondition{
				Type:     "metric",
				Operator: matches[1],
				Value:    t,
			}

			process.EvalCondition(ctx, f, v.Value, option)
		}
	}

}

// VictorMetrics 数据源
func victoriametrics(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var (
		curFiringKeys  []string
		curPendingKeys []string
	)

	defer func() {
		go process.GcPendingCache(ctx, rule, curPendingKeys)
		alertRecover(ctx, rule, curFiringKeys)
		go process.GcRecoverWaitCache(rule, curFiringKeys)
	}()

	r := models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
		Type:     "VictoriaMetrics",
	}
	datasourceInfo, err := ctx.DB.Datasource().Get(r)
	if err != nil {
		return
	}

	cmCli := client.NewVictoriaMetricsClient(datasourceInfo)
	resQuery, err := cmCli.Query(rule.PrometheusConfig.PromQL)
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

			f := func() models.AlertCurEvent {
				event := process.BuildEvent(rule)
				event.DatasourceId = datasourceId
				event.Fingerprint = v.GetFingerprint()
				event.Metric = v.GetMetric()
				event.Metric["severity"] = ruleExpr.Severity
				event.Severity = ruleExpr.Severity
				event.Annotations = cmd.ParserVariables(rule.PrometheusConfig.Annotations, event.Metric)

				firingKey := event.GetFiringAlertCacheKey()
				pendingKey := event.GetPendingAlertCacheKey()

				curFiringKeys = append(curFiringKeys, firingKey)
				curPendingKeys = append(curPendingKeys, pendingKey)

				return event
			}

			option := models.EvalCondition{
				Type:     "metric",
				Operator: matches[1],
				Value:    t,
			}

			process.EvalCondition(ctx, f, v.Value, option)
		}
	}

}

// AliCloudSLS 数据源
func aliCloudSLS(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		alertRecover(ctx, rule, curKeys)
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

	datasourceInfo, err := ctx.DB.Datasource().Get(models.DatasourceQuery{
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

		event := func() models.AlertCurEvent {
			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = body.GetFingerprint()
			event.Annotations = body.GetAnnotations()
			event.Metric = body.GetMetric()

			key := event.GetFiringAlertCacheKey()
			curKeys = append(curKeys, key)

			return event
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
		process.EvalCondition(ctx, event, float64(count), options)
	}

}

// Loki 数据源
func loki(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		alertRecover(ctx, rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	curAt := time.Now().UTC()
	startsAt := process.ParserDuration(curAt, rule.LokiConfig.LogScope, "m")
	args := client.QueryOptions{
		Query:   rule.LokiConfig.LogQL,
		StartAt: startsAt.Format(time.RFC3339Nano),
		EndAt:   curAt.Format(time.RFC3339Nano),
	}

	datasourceInfo, err := ctx.DB.Datasource().Get(models.DatasourceQuery{
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

	// count 用于统计日志条数
	var count int
	for _, v := range res {
		count += len(v.Values)
		event := func() models.AlertCurEvent {
			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Annotations = fmt.Sprintf("\\_\\_count\\_\\_: %d\n%s", count, v.GetAnnotations().(string))

			key := event.GetPendingAlertCacheKey()
			curKeys = append(curKeys, key)

			return event
		}

		options := models.EvalCondition{
			Type:     rule.LokiConfig.EvalCondition.Type,
			Operator: rule.LokiConfig.EvalCondition.Operator,
			Value:    rule.LokiConfig.EvalCondition.Value,
		}

		// 评估告警条件
		process.EvalCondition(ctx, event, float64(count), options)
	}
}

// Jaeger 数据源
func jaeger(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		alertRecover(ctx, rule, curKeys)
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

	datasourceInfo, err := ctx.DB.Datasource().Get(models.DatasourceQuery{
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
		event := process.BuildEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = v.GetFingerprint()
		event.Metric = v.GetMetric(rule)
		event.Annotations = v.GetAnnotations(rule, datasourceInfo)

		key := event.GetFiringAlertCacheKey()
		curKeys = append(curKeys, key)

		process.SaveAlertEvent(ctx, event)
	}

}

func cloudWatch(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		alertRecover(ctx, rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	datasourceObj, err := ctx.DB.Datasource().GetInstance(datasourceId)
	if err != nil {
		return
	}

	cfg, err := client.NewAWSCredentialCfg(datasourceObj.AWSCloudWatch.Region, datasourceObj.AWSCloudWatch.AccessKey, datasourceObj.AWSCloudWatch.SecretKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cli := cfg.CloudWatchCli()

	curAt := time.Now().UTC()
	startsAt := process.ParserDuration(curAt, rule.CloudWatchConfig.Period, "m")

	for _, endpoint := range rule.CloudWatchConfig.Endpoints {
		query := types.CloudWatchQuery{
			Endpoint:   endpoint,
			Dimension:  rule.CloudWatchConfig.Dimension,
			Period:     int32(rule.CloudWatchConfig.Period * 60),
			Namespace:  rule.CloudWatchConfig.Namespace,
			MetricName: rule.CloudWatchConfig.MetricName,
			Statistic:  rule.CloudWatchConfig.Statistic,
			Form:       startsAt,
			To:         curAt,
		}
		_, values := cloudwatch.MetricDataQuery(cli, query)
		if len(values) == 0 {
			return
		}

		event := func() models.AlertCurEvent {
			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = query.GetFingerprint()
			event.Metric = query.GetMetrics()
			event.Annotations = fmt.Sprintf("%s %s %s %s %d", query.Namespace, query.MetricName, query.Statistic, rule.CloudWatchConfig.Expr, rule.CloudWatchConfig.Threshold)

			return event
		}

		options := models.EvalCondition{
			Type:     "metric",
			Operator: rule.CloudWatchConfig.Expr,
			Value:    float64(rule.CloudWatchConfig.Threshold),
		}

		process.EvalCondition(ctx, event, values[0], options)
	}
}

func kubernetesEvent(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		alertRecover(ctx, rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	datasourceObj, err := ctx.DB.Datasource().GetInstance(datasourceId)
	if err != nil {
		return
	}

	cli, err := client.NewKubernetesClient(ctx.Ctx, datasourceObj.KubeConfig)
	if err != nil {
		return
	}
	event, err := cli.GetWarningEvent(rule.KubernetesConfig.Reason, rule.KubernetesConfig.Scope)
	if err != nil {
		return
	}

	if len(event.Items) < rule.KubernetesConfig.Value {
		return
	}

	var eventMapping = make(map[string][]string)
	for _, item := range process.FilterKubeEvent(event, rule.KubernetesConfig.Filter).Items {
		// 同一个资源可能有多条不同的事件信息
		eventMapping[item.InvolvedObject.Name] = append(eventMapping[item.InvolvedObject.Name], "\n"+strings.ReplaceAll(item.Message, "\"", "'"))
		k8sItem := process.KubernetesAlertEvent(ctx, item)
		alertEvent := process.BuildEvent(rule)
		alertEvent.DatasourceId = datasourceId
		alertEvent.Fingerprint = k8sItem.GetFingerprint()
		alertEvent.Metric = k8sItem.GetMetrics()
		alertEvent.Annotations = fmt.Sprintf("\n- 环境: %s\n- 命名空间: %s\n- 资源类型: %s\n- 资源名称: %s\n- 事件类型: %s\n- 事件详情: %s\n",
			datasourceObj.Name, item.Namespace, item.InvolvedObject.Kind,
			item.InvolvedObject.Name, item.Reason, eventMapping[item.InvolvedObject.Name],
		)

		process.SaveAlertEvent(ctx, alertEvent)
	}
}

func elasticSearch(ctx *ctx.Context, datasourceId string, rule models.AlertRule) {
	var curKeys []string
	defer func() {
		alertRecover(ctx, rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	datasourceInfo, err := ctx.DB.Datasource().Get(models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
	})
	if err != nil {
		return
	}

	cli, err := client.NewElasticSearchClient(ctx.Ctx, datasourceInfo)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	res, err := cli.Query(ctx.Ctx, rule.ElasticSearchConfig.Index, rule.ElasticSearchConfig.Filter, rule.ElasticSearchConfig.Scope)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	count := len(res)
	if count <= 0 {
		return
	}

	for _, v := range res {
		event := func() models.AlertCurEvent {
			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Annotations = v.GetAnnotations()

			key := event.GetPendingAlertCacheKey()
			curKeys = append(curKeys, key)

			return event
		}

		options := models.EvalCondition{
			Type:     "count",
			Operator: ">",
			Value:    1,
		}

		// 评估告警条件
		process.EvalCondition(ctx, event, float64(count), options)
	}
}
