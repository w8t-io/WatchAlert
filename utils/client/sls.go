package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sls20201230 "github.com/alibabacloud-go/sls-20201230/v6/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
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

func NewAliCloudSlsClient(datasourceId string) AliCloudSlsClientApi {

	var datasource models.AlertDataSource
	globals.DBCli.Model(&models.AlertDataSource{}).Where("id = ?", datasourceId).First(&datasource)

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
	metric map[string]interface{}
}

func GetSLSBodyData(res *sls20201230.GetLogsResponse) SlsBody {
	bodyString, _ := json.Marshal(res.Body[0])
	// 标签，用于推送告警消息时 获取相关 label 信息
	metricMap := make(map[string]interface{})
	err := json.Unmarshal(bodyString, &metricMap)
	if err != nil {
		globals.Logger.Sugar().Errorf("解析 SLS Metric Label 失败, %s", err.Error())
	}
	return SlsBody{metric: metricMap}
}

func (sb SlsBody) GetMetric() map[string]interface{} {
	// 删除多余 label
	delete(sb.metric, "_image_name_")
	delete(sb.metric, "__topic__")
	delete(sb.metric, "_container_ip_")
	delete(sb.metric, "_pod_uid_")
	delete(sb.metric, "_source_")
	delete(sb.metric, "_time_")
	delete(sb.metric, "__time__")
	delete(sb.metric, "__tag__:__pack_id__")
	return sb.metric
}

func (sb SlsBody) GetAnnotations() string {
	var annotation string
	if sb.metric["content"] != nil {
		annotation = sb.metric["content"].(string)
		if cmd.IsJSON(annotation) {
			a := cmd.FormatJson(annotation)
			annotation = a
		}
	}
	delete(sb.metric, "content")
	return annotation
}

func (sb SlsBody) GetFingerprint() string {
	h := md5.New()
	// 使用 label 进行 Hash 作为告警指纹，可以有效地作为恢复逻辑的判断条件。
	h.Write([]byte(cmd.JsonMarshal(sb.metric)))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}
