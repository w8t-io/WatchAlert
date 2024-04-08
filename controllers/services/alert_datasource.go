package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
	"watchAlert/utils/http"
)

type AlertDataSourceService struct {
	repo.DatasourceRepo
}

type InterAlertDataSourceService interface {
	Create(dataSource models.AlertDataSource) error
	Update(dataSource models.AlertDataSource) error
	Delete(tid, id string) error
	List(ctx *gin.Context) ([]models.AlertDataSource, error)
	Get(tid, id, dsType string) []models.AlertDataSource
	Search(req interface{}) (interface{}, interface{})
}

func NewInterAlertDataSourceService() InterAlertDataSourceService {
	return &AlertDataSourceService{}
}

func (adss *AlertDataSourceService) Create(dataSource models.AlertDataSource) error {

	err := adss.Check(dataSource)
	if err != nil {
		return err
	}

	id := "ds-" + cmd.RandId()

	data := models.AlertDataSource{
		TenantId:         dataSource.TenantId,
		Id:               id,
		Name:             dataSource.Name,
		Type:             dataSource.Type,
		HTTP:             dataSource.HTTP,
		AliCloudEndpoint: dataSource.AliCloudEndpoint,
		AliCloudAk:       dataSource.AliCloudAk,
		AliCloudSk:       dataSource.AliCloudSk,
		Enabled:          strconv.FormatBool(dataSource.EnabledBool),
		Description:      dataSource.Description,
	}

	err = repo.DBCli.Create(models.AlertDataSource{}, &data)
	if err != nil {
		return err
	}

	return nil

}

func (adss *AlertDataSourceService) Update(dataSource models.AlertDataSource) error {

	data := repo.Updates{
		Table: models.AlertDataSource{},
		Where: []interface{}{"id = ? AND tenant_id = ?", dataSource.Id, dataSource.TenantId},
		Updates: models.AlertDataSource{
			Id:               dataSource.Id,
			Name:             dataSource.Name,
			Type:             dataSource.Type,
			HTTP:             dataSource.HTTP,
			AliCloudEndpoint: dataSource.AliCloudEndpoint,
			AliCloudAk:       dataSource.AliCloudAk,
			AliCloudSk:       dataSource.AliCloudSk,
			Enabled:          strconv.FormatBool(dataSource.EnabledBool),
			Description:      dataSource.Description,
		},
	}

	err := repo.DBCli.Updates(data)
	if err != nil {
		return err
	}

	return nil

}

func (adss *AlertDataSourceService) Delete(tid, id string) error {

	var ruleNum int64
	globals.DBCli.Model(&models.AlertRule{}).Where("tenant_id = ? AND datasource_id LIKE ?", tid, "%"+id+"%").Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("无法删除数据源 %s, 因为已有告警规则绑定", id)
	}

	data := repo.Delete{
		Table: models.AlertDataSource{},
		Where: []interface{}{"tid = ? AND id = ?", tid, id},
	}

	err := repo.DBCli.Delete(data)
	if err != nil {
		return err
	}

	return nil

}

func (adss *AlertDataSourceService) List(ctx *gin.Context) ([]models.AlertDataSource, error) {

	var (
		data []models.AlertDataSource
	)
	tid, _ := ctx.Get("TenantID")

	db := globals.DBCli.Model(&models.AlertDataSource{})
	db.Where("tenant_id = ?", tid.(string))
	db.Find(&data)

	for k, v := range data {
		data[k].EnabledBool, _ = strconv.ParseBool(v.Enabled)
	}

	return data, nil

}

func (adss *AlertDataSourceService) Get(tid, id, dsType string) []models.AlertDataSource {

	db := globals.DBCli.Model(&models.AlertDataSource{})
	db.Where("tenant_id = ?", tid)
	db.Where("type = ?", dsType)

	if id != "" {
		db.Where("id = ?", id)
	}

	var data []models.AlertDataSource
	err := db.Find(&data).Error
	if err != nil {
		return []models.AlertDataSource{}
	}

	for k := range data {
		data[k].EnabledBool, _ = strconv.ParseBool(data[k].Enabled)
	}

	return data

}

func (adss *AlertDataSourceService) Search(req interface{}) (interface{}, interface{}) {
	var newData []models.AlertDataSource
	r := req.(*models.DatasourceQuery)
	data, err := adss.DatasourceRepo.SearchDatasource(*r)
	if err != nil {
		return nil, err
	}
	newData = data

	for k := range data {
		newData[k].EnabledBool, _ = strconv.ParseBool(data[k].Enabled)
	}

	return newData, nil
}

func (adss *AlertDataSourceService) Check(dataSource models.AlertDataSource) error {

	switch dataSource.Type {
	case "Prometheus":
		path := "/api/v1/format_query?query=foo/bar"
		fullPath := dataSource.HTTP.URL + path
		res, err := http.Get(fullPath)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("StatusCode 非预期 -> %d", res.StatusCode)
		}
	}

	return nil

}
