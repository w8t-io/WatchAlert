package provider

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sls20201230 "github.com/alibabacloud-go/sls-20201230/v6/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

type AliCloudSlsDsProvider struct {
	client *sls20201230.Client
}

func NewAliCloudSlsClient(source models.AlertDataSource) (LogsFactoryProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     &source.AliCloudAk,
		AccessKeySecret: &source.AliCloudSk,
	}
	config.Endpoint = tea.String(source.AliCloudEndpoint)
	result, err := sls20201230.NewClient(config)
	if err != nil {
		return AliCloudSlsDsProvider{}, err
	}

	return AliCloudSlsDsProvider{client: result}, nil
}

func (a AliCloudSlsDsProvider) Query(query LogQueryOptions) ([]Logs, int, error) {
	var err error
	getLogsRequest := &sls20201230.GetLogsRequest{
		To:    tea.Int32(query.EndAt.(int32)),
		From:  tea.Int32(query.StartAt.(int32)),
		Query: tea.String(query.AliCloudSLS.Query),
	}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	defer func() {
		if r := tea.Recover(recover()); r != nil {
			err = r
		}
	}()

	res, err := a.client.GetLogsWithOptions(tea.String(query.AliCloudSLS.Project), tea.String(query.AliCloudSLS.LogStore), getLogsRequest, headers, runtime)
	if err != nil {
		return nil, 0, err
	}

	var (
		msgList []interface{}
		metric  = map[string]interface{}{}
	)
	for _, body := range res.Body {
		msg := cmd.FormatJson(cmd.JsonMarshal(body))
		msgList = append(msgList, msg)

		metric["_container_name_"] = body["__tag__:_container_name_"]
		metric["_namespace_"] = body["__tag__:_namespace_"]
	}

	var data []Logs
	data = append(data, Logs{
		ProviderName: AliCloudSLSDsProviderName,
		Metric:       metric,
		Message:      msgList,
	})

	return data, len(msgList), nil
}

func (a AliCloudSlsDsProvider) Check() (bool, error) {

	return true, nil
}
