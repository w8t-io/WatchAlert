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
	"watchAlert/pkg/community/aws/cloudwatch"
	"watchAlert/pkg/community/aws/cloudwatch/types"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
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

// Metrics 包含 Prometheus、VictoriaMetrics 数据源
func metrics(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) {
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
		Type:     datasourceType,
	}
	datasourceInfo, err := ctx.DB.Datasource().Get(r)
	if err != nil {
		return
	}

	health := provider.CheckDatasourceHealth(datasourceInfo)
	if !health {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	var resQuery []provider.Metrics
	switch datasourceType {
	case provider.PrometheusDsProvider:
		prometheusClient, err := provider.NewPrometheusClient(datasourceInfo)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		resQuery, err = prometheusClient.Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}
	case provider.VictoriaMetricsDsProvider:
		vmClient, err := provider.NewVictoriaMetricsClient(datasourceInfo)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		resQuery, err = vmClient.Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}
	default:
		global.Logger.Sugar().Errorf("Unsupported metrics type, type: %s", datasourceType)
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

// Logs 包含 AliSLS、Loki、ElasticSearch 数据源
func logs(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) {
	var (
		curKeys     []string
		queryRes    []provider.Logs
		count       int
		err         error
		evalOptions models.EvalCondition
	)
	defer func() {
		alertRecover(ctx, rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	r := models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
		Type:     datasourceType,
	}
	datasourceInfo, err := ctx.DB.Datasource().Get(r)
	if err != nil {
		return
	}

	health := provider.CheckDatasourceHealth(datasourceInfo)
	if !health {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	switch datasourceType {
	case provider.LokiDsProviderName:
		lokiCli, err := provider.NewLokiClient(datasourceInfo)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		curAt := time.Now()
		startsAt := process.ParserDuration(curAt, rule.LokiConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			Loki: provider.Loki{
				Query: rule.LokiConfig.LogQL,
			},
			StartAt: startsAt.Unix(),
			EndAt:   curAt.Unix(),
		}
		queryRes, count, err = lokiCli.Query(queryOptions)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		evalOptions = models.EvalCondition{
			Type:     rule.LokiConfig.EvalCondition.Type,
			Operator: rule.LokiConfig.EvalCondition.Operator,
			Value:    rule.LokiConfig.EvalCondition.Value,
		}
	case provider.AliCloudSLSDsProviderName:
		slsClient, err := provider.NewAliCloudSlsClient(datasourceInfo)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		curAt := time.Now()
		startsAt := process.ParserDuration(curAt, rule.AliCloudSLSConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			AliCloudSLS: provider.AliCloudSLS{
				Query:    rule.AliCloudSLSConfig.LogQL,
				Project:  rule.AliCloudSLSConfig.Project,
				LogStore: rule.AliCloudSLSConfig.Logstore,
			},
			StartAt: int32(startsAt.Unix()),
			EndAt:   int32(curAt.Unix()),
		}
		queryRes, count, err = slsClient.Query(queryOptions)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		evalOptions = models.EvalCondition{
			Type:     rule.AliCloudSLSConfig.EvalCondition.Type,
			Operator: rule.AliCloudSLSConfig.EvalCondition.Operator,
			Value:    rule.AliCloudSLSConfig.EvalCondition.Value,
		}
	case provider.ElasticSearchDsProviderName:
		searchClient, err := provider.NewElasticSearchClient(ctx.Ctx, datasourceInfo)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		curAt := time.Now()
		startsAt := process.ParserDuration(curAt, int(rule.ElasticSearchConfig.Scope), "m")
		queryOptions := provider.LogQueryOptions{
			ElasticSearch: provider.Elasticsearch{
				Index:       rule.ElasticSearchConfig.Index,
				QueryFilter: rule.ElasticSearchConfig.Filter,
			},
			StartAt: cmd.FormatTimeToUTC(startsAt.Unix()),
			EndAt:   cmd.FormatTimeToUTC(curAt.Unix()),
		}
		queryRes, count, err = searchClient.Query(queryOptions)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		evalOptions = models.EvalCondition{
			Type:     "count",
			Operator: ">",
			Value:    1,
		}
	}

	if count <= 0 {
		return
	}

	for _, v := range queryRes {
		event := func() models.AlertCurEvent {
			event := process.BuildEvent(rule)
			event.DatasourceId = datasourceId
			event.Fingerprint = v.GetFingerprint()
			event.Metric = v.GetMetric()
			event.Annotations = fmt.Sprintf("统计日志条数: %d 条\n%s", count, cmd.FormatJson(v.GetAnnotations()[0].(string)))

			key := event.GetPendingAlertCacheKey()
			curKeys = append(curKeys, key)

			return event
		}

		// 评估告警条件
		process.EvalCondition(ctx, event, float64(count), evalOptions)
	}
}

// Traces 包含 Jaeger 数据源
func traces(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) {
	var (
		curKeys  []string
		queryRes []provider.Traces
		err      error
	)
	defer func() {
		alertRecover(ctx, rule, curKeys)
		go process.GcRecoverWaitCache(rule, curKeys)
	}()

	r := models.DatasourceQuery{
		TenantId: rule.TenantId,
		Id:       datasourceId,
		Type:     datasourceType,
	}
	datasourceInfo, err := ctx.DB.Datasource().Get(r)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	health := provider.CheckDatasourceHealth(datasourceInfo)
	if !health {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	switch datasourceType {
	case provider.JaegerDsProviderName:
		curAt := time.Now().UTC()
		startsAt := process.ParserDuration(curAt, rule.JaegerConfig.Scope, "m")

		jaegerClient, err := provider.NewJaegerClient(datasourceInfo)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}

		queryOptions := provider.TraceQueryOptions{
			Tags:    rule.JaegerConfig.Tags,
			Service: rule.JaegerConfig.Service,
			StartAt: startsAt.UnixMicro(),
			EndAt:   curAt.UnixMicro(),
		}
		queryRes, err = jaegerClient.Query(queryOptions)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}
	}

	for _, v := range queryRes {
		event := process.BuildEvent(rule)
		event.DatasourceId = datasourceId
		event.Fingerprint = v.GetFingerprint()
		event.Metric = v.GetMetric()
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

	cfg, err := provider.NewAWSCredentialCfg(datasourceObj.AWSCloudWatch.Region, datasourceObj.AWSCloudWatch.AccessKey, datasourceObj.AWSCloudWatch.SecretKey)
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
		global.Logger.Sugar().Error(err.Error())
		return
	}

	cli, err := provider.NewKubernetesClient(ctx.Ctx, datasourceObj.KubeConfig)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}
	event, err := cli.GetWarningEvent(rule.KubernetesConfig.Reason, rule.KubernetesConfig.Scope)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
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
		alertEvent.Annotations = fmt.Sprintf("- 环境: %s\n- 命名空间: %s\n- 资源类型: %s\n- 资源名称: %s\n- 事件类型: %s\n- 事件详情: %s\n",
			datasourceObj.Name, item.Namespace, item.InvolvedObject.Kind,
			item.InvolvedObject.Name, item.Reason, eventMapping[item.InvolvedObject.Name],
		)

		process.SaveAlertEvent(ctx, alertEvent)
	}
}
