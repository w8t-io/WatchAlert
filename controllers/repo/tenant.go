package repo

import (
	"watchAlert/models"
	"watchAlert/public/globals"
)

type TenantRepo struct{}

func (tr TenantRepo) CreateTenant(t models.Tenant) error {
	err := DBCli.Create(&models.Tenant{}, t)
	if err != nil {
		globals.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (tr TenantRepo) UpdateTenant(t models.Tenant) error {
	u := Updates{
		Table:   &models.Tenant{},
		Where:   []interface{}{"id = ?", t.ID},
		Updates: t,
	}
	err := DBCli.Updates(u)
	if err != nil {
		globals.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (tr TenantRepo) DeleteTenant(t models.Tenant) error {
	err := DBCli.Delete(Delete{
		Table: &models.Tenant{},
		Where: []interface{}{"id = ?", t.ID},
	})
	if err != nil {
		globals.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (tr TenantRepo) ListTenant() (data []models.Tenant, err error) {
	var d []models.Tenant
	err = globals.DBCli.Model(&models.Tenant{}).Find(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}
