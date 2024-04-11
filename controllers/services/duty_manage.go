package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type DutyManageService struct {
	repo.DutyRepo
}

type InterDutyManageService interface {
	ListDutyManage(ctx *gin.Context) []models.DutyManagement
	CreateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error)
	UpdateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error)
	DeleteDutyManage(tid, id string) error
	GetDutyManage(tid, id string) models.DutyManagement
}

func NewInterDutyManageService() InterDutyManageService {
	return &DutyManageService{}
}

func (dms *DutyManageService) ListDutyManage(ctx *gin.Context) []models.DutyManagement {

	var list []models.DutyManagement
	tid, _ := ctx.Get("TenantID")
	globals.DBCli.Model(&models.DutyManagement{}).Where("tenant_id = ?", tid.(string)).Find(&list)
	return list

}

func (dms *DutyManageService) CreateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error) {

	ok := dms.DutyRepo.GetQuota(dutyManage.TenantId)
	if !ok {
		return models.DutyManagement{}, fmt.Errorf("创建失败, 配额不足")
	}

	tx := globals.DBCli.Begin()
	dutyManage.ID = "dt-" + cmd.RandId()
	dutyManage.CreateAt = time.Now().Unix()

	err := tx.Create(&dutyManage).Error
	if err != nil {
		tx.Rollback()
		return models.DutyManagement{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return models.DutyManagement{}, err
	}
	return dutyManage, nil

}

func (dms *DutyManageService) UpdateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error) {

	tx := globals.DBCli.Begin()
	err := tx.Model(&models.DutyManagement{}).Where("tenant_id = ? AND id = ?", dutyManage.TenantId, dutyManage.ID).Updates(&dutyManage).Error
	if err != nil {
		tx.Rollback()
		return models.DutyManagement{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return models.DutyManagement{}, err
	}
	return dutyManage, nil

}

func (dms *DutyManageService) DeleteDutyManage(tid, id string) error {

	var noticeNum int64
	globals.DBCli.Model(&models.AlertNotice{}).Where("tenant_id = ? AND duty_id = ?", tid, id).Count(&noticeNum)
	if noticeNum != 0 {
		return fmt.Errorf("无法删除值班表 %s, 因为已有通知对象绑定", id)
	}

	tx := globals.DBCli.Begin()
	err := tx.Where("tenant_id = ? AND id = ?", tid, id).Delete(&models.DutyManagement{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&models.DutySchedule{}).Where("tenant_id = ? AND duty_id = ?", tid, id).Delete(&models.DutySchedule{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

func (dms *DutyManageService) GetDutyManage(tid, id string) models.DutyManagement {

	var data models.DutyManagement
	globals.DBCli.Model(&models.DutyManagement{}).Where("tenant_id = ? AND id = ?", tid, id).Find(&data)
	return data

}
