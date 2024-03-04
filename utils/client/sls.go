package client

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sls20201230 "github.com/alibabacloud-go/sls-20201230/v6/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"watchAlert/globals"
	"watchAlert/models"
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
