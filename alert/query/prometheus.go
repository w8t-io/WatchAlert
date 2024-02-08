package query

import (
	"encoding/json"
	"io/ioutil"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/services"
	"watchAlert/globals"
	models "watchAlert/models"
	"watchAlert/utils/hash"
	utilHttp "watchAlert/utils/http"
)

type Prometheus struct {
	alertEvent models.AlertCurEvent
}

func (p *Prometheus) Query(rule models.AlertRule) {

	for _, dsId := range rule.DatasourceIdList {

		datasource := services.NewInterAlertDataSourceService().Get(dsId, rule.DatasourceType)

		url := datasource[0].HTTPJson.URL + "/api/v1/query?query=" + rule.RuleConfigJson.PromQL
		res, err := utilHttp.Get(url)
		if err != nil {
			return
		}

		var (
			ResQuery models.PromQueryResponse
			curKeys  []string
		)
		byteBody, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(byteBody, &ResQuery)

		var curValue []string

		if ResQuery.Data.Result != nil {

			for _, v := range ResQuery.Data.Result {
				curValue = v.Value

				metricJson, _ := json.Marshal(v.Metric)
				// fingerprint 用于报警指纹，根据 Prom Query 查到的Metric Label
				fingerprint := hash.Md5Hash(metricJson)
				fingerprint = fingerprint[:15]

				key := p.alertEvent.CurAlertCacheKey(rule.RuleId, fingerprint)
				curKeys = append(curKeys, key)

				event := models.AlertCurEvent{
					DatasourceType:       rule.DatasourceType,
					DatasourceIdList:     []string{dsId},
					Fingerprint:          fingerprint,
					RuleId:               rule.RuleId,
					RuleName:             rule.RuleName,
					Severity:             rule.RuleConfigJson.Severity,
					Instance:             v.Metric["instance"],
					Metric:               string(metricJson),
					MetricMap:            v.Metric,
					CurValue:             v.Value,
					PromQl:               rule.RuleConfigJson.PromQL,
					LabelsMap:            rule.LabelsMap,
					Labels:               rule.Labels,
					EvalInterval:         rule.EvalInterval,
					ForDuration:          rule.ForDuration,
					NoticeId:             rule.NoticeId,
					NoticeGroupList:      rule.NoticeGroupList,
					IsRecovered:          false,
					RepeatNoticeInterval: rule.RepeatNoticeInterval,
					DutyUser:             "暂无", // 默认暂无值班人员, 渲染模版时会实际判断 Notice 是否存在值班人员
				}
				event.Annotations = event.ParserAnnotation(rule.Annotations)
				event.FirstTriggerTime = event.GetFirstTime()
				event.LastEvalTime = event.GetLastEvalTime()
				event.LastSendTime = event.GetLastSendTime()

				p.alertEvent.SetCache(event, 0)
				err = repo.DBCli.Create(models.AlertCurEvent{}, &event)
				if err != nil {
					return
				}

			}

		}

		allKey := p.alertCacheKeys(rule)

		recoverKeys := p.getSliceDifference(allKey, curKeys)

		for _, key := range recoverKeys {
			event := p.alertEvent.GetCache(key)
			if event.IsRecovered == true {
				continue
			}
			event.CurValue = curValue
			event.IsRecovered = true
			event.RecoverTime = time.Now().Unix()
			event.LastSendTime = 0
			p.alertEvent.SetCache(event, 0)
		}

	}

}

func (p *Prometheus) alertCacheKeys(rule models.AlertRule) []string {

	// 获取所有keys
	keyPrefix := p.alertEvent.CurAlertCacheKey(rule.RuleId, "*")
	keys, _ := globals.RedisCli.Keys(keyPrefix).Result()

	return keys

}

func (p *Prometheus) getSliceDifference(slice1 []string, slice2 []string) []string {
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
