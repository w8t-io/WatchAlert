package repo

import (
	"watchAlert/controllers/dao"
	"watchAlert/globals"
)

type AlertNoticeRepo struct{}

func (anr *AlertNoticeRepo) Get(uuid string) dao.AlertNotice {

	var alertNoticeData dao.AlertNotice
	globals.DBCli.Model(&dao.AlertNotice{}).Where("uuid = ?", uuid).Find(&alertNoticeData)
	return alertNoticeData

}
