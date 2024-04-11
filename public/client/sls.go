package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sls20201230 "github.com/alibabacloud-go/sls-20201230/v6/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type AliCloudSlsClientApi struct {
	client *sls20201230.Client
}

type AliCloudSlsQueryArgs struct {
	Project  string
	Logstore string
	StartsAt int32
	EndsAt   int32
	Query    string
}

func NewAliCloudSlsClient(tid, datasourceId string) AliCloudSlsClientApi {

	var datasource models.AlertDataSource
	globals.DBCli.Model(&models.AlertDataSource{}).Where("tenant_id = ? AND id = ?", tid, datasourceId).First(&datasource)

	config := &openapi.Config{
		AccessKeyId:     &datasource.AliCloudAk,
		AccessKeySecret: &datasource.AliCloudSk,
	}
	config.Endpoint = tea.String(datasource.AliCloudEndpoint)
	result, err := sls20201230.NewClient(config)
	if err != nil {
		globals.Logger.Sugar().Errorf("创建 SLS 客户端失败 -> %s", err.Error())
		return AliCloudSlsClientApi{}
	}

	return AliCloudSlsClientApi{client: result}

}

func (sca AliCloudSlsClientApi) Query(args AliCloudSlsQueryArgs) (res *sls20201230.GetLogsResponse, err error) {

	getLogsRequest := &sls20201230.GetLogsRequest{
		To:    tea.Int32(args.EndsAt),
		From:  tea.Int32(args.StartsAt),
		Query: tea.String(args.Query),
	}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	defer func() {
		if r := tea.Recover(recover()); r != nil {
			err = r
		}
	}()

	res, err = sca.client.GetLogsWithOptions(tea.String(args.Project), tea.String(args.Logstore), getLogsRequest, headers, runtime)
	if err != nil {
		return nil, err
	}

	return res, nil

}

type SlsBody struct {
	MetricList []Metric
}
type Metric map[string]interface{}

func GetSLSBodyData(res *sls20201230.GetLogsResponse) SlsBody {

	var metricMapList []Metric
	for _, body := range res.Body {
		// 标签，用于推送告警消息时 获取相关 label 信息
		metricMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(cmd.JsonMarshal(body)), &metricMap)
		if err != nil {
			globals.Logger.Sugar().Errorf("解析 SLS Metric Label 失败, %s", err.Error())
		}
		metricMapList = append(metricMapList, metricMap)
	}

	return SlsBody{MetricList: metricMapList}
}

func (m Metric) GetMetric() map[string]interface{} {
	// 删除多余 label
	delete(m, "_image_name_")
	delete(m, "__topic__")
	delete(m, "_container_ip_")
	delete(m, "_pod_uid_")
	delete(m, "_source_")
	delete(m, "_time_")
	delete(m, "__time__")
	delete(m, "__tag__:__pack_id__")
	return m
}

func (m Metric) GetAnnotations() string {
	var annotation string
	if m["content"] != nil {
		annotation = m["content"].(string)
		if cmd.IsJSON(annotation) {
			a := cmd.FormatJson(annotation)
			annotation = a
		}
	}
	delete(m, "content")
	return annotation
}

func (m Metric) GetFingerprint() string {
	// 使用 label 进行 Hash 作为告警指纹，可以有效地作为恢复逻辑的判断条件。
	newMetric := map[string]interface{}{
		"_namespace_":      m["_namespace_"],
		"_container_name_": m["_container_name_"],
	}
	h := md5.New()
	h.Write([]byte(cmd.JsonMarshal(newMetric)))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}
