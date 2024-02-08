package repo

import (
	"watchAlert/globals"
	"watchAlert/models"
)

type AlertNoticeRepo struct{}

func (anr *AlertNoticeRepo) Get(uuid string) models.AlertNotice {

	var alertNoticeData models.AlertNotice
	globals.DBCli.Model(&models.AlertNotice{}).Where("uuid = ?", uuid).Find(&alertNoticeData)
	return alertNoticeData

}
