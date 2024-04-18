package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/repo"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type AlertNoticeService struct {
	repo.NoticeRepo
}

type InterAlertNoticeService interface {
	ListNoticeObject(ctx *gin.Context) []models.AlertNotice
	CreateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error)
	UpdateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error)
	DeleteNoticeObject(tid, uuid string) error
	GetNoticeObject(tid, uuid string) models.AlertNotice
	CheckNoticeObjectStatus(tid, uuid string) string
	Search(req interface{}) (interface{}, interface{})
}

func NewInterAlertNoticeService() InterAlertNoticeService {
	return &AlertNoticeService{}
}

func (ans *AlertNoticeService) ListNoticeObject(ctx *gin.Context) []models.AlertNotice {
	db := globals.DBCli.Model(&models.AlertNotice{})
	tid, _ := ctx.Get("TenantID")

	var alertNoticeObject []models.AlertNotice
	db.Where("tenant_id = ?", tid.(string))
	db.Find(&alertNoticeObject)
	return alertNoticeObject
}

func (ans *AlertNoticeService) CreateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error) {

	ok := ans.NoticeRepo.GetQuota(alertNotice.TenantId)
	if !ok {
		return models.AlertNotice{}, fmt.Errorf("创建失败, 配额不足")
	}

	tx := globals.DBCli.Begin()
	alertNotice.Uuid = "n-" + cmd.RandId()
	err := tx.Create(alertNotice).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("创建通知对象失败", err)
		return models.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("事务提交失败", err)
		return models.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) UpdateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	err := tx.Model(&models.AlertNotice{}).Where("tenant_id = ? AND uuid = ?", alertNotice.TenantId, alertNotice.Uuid).Updates(&alertNotice).Error
	if err != nil {
		tx.Rollback()
		return models.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return models.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) DeleteNoticeObject(tid, uuid string) error {

	var ruleNum1, ruleNum2 int64
	db := globals.DBCli.Model(&models.AlertRule{})
	db.Where("notice_id = ?", uuid).Count(&ruleNum1)
	db.Where("notice_group LIKE ?", "%"+uuid+"%").Count(&ruleNum2)
	if ruleNum1 != 0 || ruleNum2 != 0 {
		return fmt.Errorf("无法删除通知对象 %s, 因为已有告警规则绑定", uuid)
	}

	tx := globals.DBCli.Begin()
	err := tx.Where("tenant_id = ? AND uuid = ?", tid, uuid).Delete(&models.AlertNotice{}).Error
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

func (ans *AlertNoticeService) GetNoticeObject(tid, uuid string) models.AlertNotice {

	var alertNoticeObject models.AlertNotice
	globals.DBCli.Where("tenant_id = ? AND uuid = ?", tid, uuid).Find(&alertNoticeObject)
	return alertNoticeObject

}

func (ans *AlertNoticeService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeQuery)
	data, err := ans.NoticeRepo.Search(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ans *AlertNoticeService) CheckNoticeObjectStatus(tid, uuid string) string {

	// ToDo

	return ""
}
