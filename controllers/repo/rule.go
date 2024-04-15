package repo

import (
	"watchAlert/models"
	"watchAlert/public/globals"
)

type RuleRepo struct{}

func (rr RuleRepo) GetQuota(id string) bool {
	var (
		db     = globals.DBCli.Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	globals.DBCli.Model(&models.AlertRule{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.RuleNumber {
		return true
	}

	return false
}
