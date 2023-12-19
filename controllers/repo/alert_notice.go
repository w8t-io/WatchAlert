package repo

import (
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
)

type AlertNoticeRepo struct{}

func (anr *AlertNoticeRepo) Get(uuid string) dao.AlertNotice {

	var alertNoticeData dao.AlertNotice
	globals.DBCli.Model(&dao.AlertNotice{}).Where("uuid = ?", uuid).Find(&alertNoticeData)
	return alertNoticeData

}
