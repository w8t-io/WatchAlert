package services

import (
	"encoding/json"
	"fmt"
	"io"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/http"
)

type ClientService struct{}

type InterClientService interface {
	GetJaegerService(req interface{}) (interface{}, interface{})
}

func NewInterClientService() InterClientService {
	return &ClientService{}
}

type JaegerServiceData struct {
	Data []string `json:"data"`
}

func (cs ClientService) GetJaegerService(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DatasourceQuery)
	var data models.AlertDataSource
	err := globals.DBCli.Model(&models.AlertDataSource{}).Where("id = ?", r.Id).First(&data).Error
	if err != nil {
		return nil, err
	}

	url := data.HTTP.URL + "/api/services"
	res, err := http.Get(url)
	if err != nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("后端服务请求异常, 上游返回 %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	var resData JaegerServiceData
	err = json.Unmarshal(body, &resData)
	if err != nil {
		return nil, err
	}
	return resData, nil
}
