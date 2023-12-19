package services

import (
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
	"prometheus-manager/utils/cmd"
)

type AlertNoticeService struct{}

type InterAlertNoticeService interface {
	SearchNoticeObject() []dao.AlertNotice
	CreateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error)
	UpdateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error)
	DeleteNoticeObject(uuid string) error
	GetNoticeObject(uuid string) dao.AlertNotice
}

func NewInterAlertNoticeService() InterAlertNoticeService {
	return &AlertNoticeService{}
}

func (ans *AlertNoticeService) SearchNoticeObject() []dao.AlertNotice {

	var alertNoticeObject []dao.AlertNotice
	globals.DBCli.Find(&alertNoticeObject)
	return alertNoticeObject

}

func (ans *AlertNoticeService) CreateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	alertNotice.Uuid = cmd.RandUuid()
	err := tx.Create(alertNotice).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("创建通知对象失败", err)
		return dao.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("事务提交失败", err)
		return dao.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) UpdateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	err := tx.Model(&dao.AlertNotice{}).Where("uuid = ?", alertNotice.Uuid).Updates(&alertNotice).Error
	if err != nil {
		tx.Rollback()
		return dao.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return dao.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) DeleteNoticeObject(uuid string) error {

	tx := globals.DBCli.Begin()
	err := tx.Where("uuid = ?", uuid).Delete(&dao.AlertNotice{}).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("删除通知对象失败", err)
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("事务提交失败", err)
		return err
	}
	return nil

}

func (ans *AlertNoticeService) GetNoticeObject(uuid string) dao.AlertNotice {

	var alertNoticeObject dao.AlertNotice
	globals.DBCli.Where("uuid = ?", uuid).Find(&alertNoticeObject)
	return alertNoticeObject

}
