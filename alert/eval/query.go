package eval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"watchAlert/alert/process"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/pkg/community/aws/cloudwatch"
	"watchAlert/pkg/community/aws/cloudwatch/types"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

// Metrics 包含 Prometheus、VictoriaMetrics 数据源
func metrics(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) (curFiringKeys, curPendingKeys []string) {
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
		return
	}

	pools := ctx.Redis.ProviderPools()
	var resQuery []provider.Metrics
	switch datasourceType {
	case provider.PrometheusDsProvider:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		resQuery, err = cli.(provider.PrometheusProvider).Query(rule.PrometheusConfig.PromQL)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return
		}
	case provider.VictoriaMetricsDsProvider:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		resQuery, err = cli.(provider.VictoriaMetricsProvider).Query(rule.PrometheusConfig.PromQL)
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
				event.Annotations = tools.ParserVariables(rule.PrometheusConfig.Annotations, event.Metric)

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

	return
}

// Logs 包含 AliSLS、Loki、ElasticSearch 数据源
func logs(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) (curFiringKeys []string) {
	var (
		queryRes    []provider.Logs
		count       int
		err         error
		evalOptions models.EvalCondition
	)

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
		return
	}

	pools := ctx.Redis.ProviderPools()
	switch datasourceType {
	case provider.LokiDsProviderName:
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		curAt := time.Now()
		startsAt := tools.ParserDuration(curAt, rule.LokiConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			Loki: provider.Loki{
				Query: rule.LokiConfig.LogQL,
			},
			StartAt: startsAt.Unix(),
			EndAt:   curAt.Unix(),
		}
		queryRes, count, err = cli.(provider.LokiProvider).Query(queryOptions)
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
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		curAt := time.Now()
		startsAt := tools.ParserDuration(curAt, rule.AliCloudSLSConfig.LogScope, "m")
		queryOptions := provider.LogQueryOptions{
			AliCloudSLS: provider.AliCloudSLS{
				Query:    rule.AliCloudSLSConfig.LogQL,
				Project:  rule.AliCloudSLSConfig.Project,
				LogStore: rule.AliCloudSLSConfig.Logstore,
			},
			StartAt: int32(startsAt.Unix()),
			EndAt:   int32(curAt.Unix()),
		}
		queryRes, count, err = cli.(provider.AliCloudSlsDsProvider).Query(queryOptions)
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
		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		curAt := time.Now()
		startsAt := tools.ParserDuration(curAt, int(rule.ElasticSearchConfig.Scope), "m")
		queryOptions := provider.LogQueryOptions{
			ElasticSearch: provider.Elasticsearch{
				Index:       rule.ElasticSearchConfig.Index,
				QueryFilter: rule.ElasticSearchConfig.Filter,
			},
			StartAt: tools.FormatTimeToUTC(startsAt.Unix()),
			EndAt:   tools.FormatTimeToUTC(curAt.Unix()),
		}
		queryRes, count, err = cli.(provider.ElasticSearchDsProvider).Query(queryOptions)
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
			event.Annotations = fmt.Sprintf("统计日志条数: %d 条\n%s", count, tools.FormatJson(v.GetAnnotations()[0].(string)))

			key := event.GetPendingAlertCacheKey()
			curFiringKeys = append(curFiringKeys, key)

			return event
		}

		// 评估告警条件
		process.EvalCondition(ctx, event, float64(count), evalOptions)
	}

	return
}

// Traces 包含 Jaeger 数据源
func traces(ctx *ctx.Context, datasourceId, datasourceType string, rule models.AlertRule) (curFiringKeys []string) {
	var (
		queryRes []provider.Traces
		err      error
	)

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
		return
	}

	pools := ctx.Redis.ProviderPools()
	switch datasourceType {
	case provider.JaegerDsProviderName:
		curAt := time.Now().UTC()
		startsAt := tools.ParserDuration(curAt, rule.JaegerConfig.Scope, "m")

		cli, err := pools.GetClient(datasourceId)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		queryOptions := provider.TraceQueryOptions{
			Tags:    rule.JaegerConfig.Tags,
			Service: rule.JaegerConfig.Service,
			StartAt: startsAt.UnixMicro(),
			EndAt:   curAt.UnixMicro(),
		}
		queryRes, err = cli.(provider.JaegerDsProvider).Query(queryOptions)
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
		curFiringKeys = append(curFiringKeys, key)

		process.SaveAlertEvent(ctx, event)
	}

	return
}

func cloudWatch(ctx *ctx.Context, datasourceId string, rule models.AlertRule) (curFiringKeys []string) {
	pools := ctx.Redis.ProviderPools()
	cfg, err := pools.GetClient(datasourceId)
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return
	}
	cli := cfg.(provider.AwsConfig).CloudWatchCli()
	curAt := time.Now().UTC()
	startsAt := tools.ParserDuration(curAt, rule.CloudWatchConfig.Period, "m")

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
			return []string{}
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

	return
}

func kubernetesEvent(ctx *ctx.Context, datasourceId string, rule models.AlertRule) (curFiringKeys []string) {
	datasourceObj, err := ctx.DB.Datasource().GetInstance(datasourceId)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	pools := ctx.Redis.ProviderPools()
	cli, err := pools.GetClient(datasourceId)
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return
	}

	event, err := cli.(provider.KubernetesClient).GetWarningEvent(rule.KubernetesConfig.Reason, rule.KubernetesConfig.Scope)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return
	}

	if len(event.Items) < rule.KubernetesConfig.Value {
		return []string{}
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

	return
}
