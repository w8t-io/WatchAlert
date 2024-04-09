package repo

import (
	"watchAlert/models"
	"watchAlert/public/globals"
)

type DashboardRepo struct{}

func (dr DashboardRepo) CreateDashboard(d models.Dashboard) error {
	err := DBCli.Create(&models.Dashboard{}, d)
	if err != nil {
		globals.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) UpdateDashboard(d models.Dashboard) error {
	u := Updates{
		Table:   &models.Dashboard{},
		Where:   []interface{}{"tenant_id = ? AND id = ?", d.TenantId, d.ID},
		Updates: d,
	}
	err := DBCli.Updates(u)
	if err != nil {
		globals.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) DeleteDashboard(d models.DashboardQuery) error {
	del := Delete{
		Table: &models.Dashboard{},
		Where: []interface{}{"tenant_id = ? AND id = ?", d.TenantId, d.ID},
	}
	err := DBCli.Delete(del)
	if err != nil {
		globals.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) SearchDashboard(d models.DashboardQuery) ([]models.Dashboard, error) {
	var db = globals.DBCli.Model(&models.Dashboard{})
	var data []models.Dashboard
	if d.Query != "" {
		db.Where("tenant_id = ? AND name LIKE ? OR description LIKE ? OR url LIKE ?", d.TenantId, "%"+d.Query+"%", "%"+d.Query+"%", "%"+d.Query+"%")
	} else {
		db.Where("tenant_id = ?", d.TenantId).Find(&data)
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
