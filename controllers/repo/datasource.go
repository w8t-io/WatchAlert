package repo

import (
	"watchAlert/models"
	"watchAlert/public/globals"
)

type DatasourceRepo struct{}

func (ds DatasourceRepo) SearchDatasource(r models.DatasourceQuery) ([]models.AlertDataSource, error) {
	var db = globals.DBCli.Model(&models.AlertDataSource{})
	var data []models.AlertDataSource

	db.Where("tenant_id = ?", r.TenantId)
	if r.Id != "" {
		db.Where("id = ?", r.Id)
	}
	if r.Type != "" {
		db.Where("type = ?", r.Type)
	}
	if r.Query != "" {
		db.Where("id = ? OR name = ? OR description = ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	}
	if r.Id == "" && r.Type == "" && r.Query == "" {
		err := db.Find(&data).Error
		if err != nil {
			return nil, err
		}
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
