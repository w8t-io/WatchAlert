package repo

import (
	"watchAlert/globals"
	"watchAlert/models"
)

type DutyRepo struct{}

func (nr DutyRepo) GetQuota(id string) bool {
	var (
		db     = globals.DBCli.Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	globals.DBCli.Model(&models.DutyManagement{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.NoticeNumber {
		return true
	}

	return false
}
