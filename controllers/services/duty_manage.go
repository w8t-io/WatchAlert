package services

import (
	"fmt"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type DutyManageService struct{}

type InterDutyManageService interface {
	ListDutyManage() []models.DutyManagement
	CreateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error)
	UpdateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error)
	DeleteDutyManage(id string) error
	GetDutyManage(id string) models.DutyManagement
}

func NewInterDutyManageService() InterDutyManageService {
	return &DutyManageService{}
}

func (dms *DutyManageService) ListDutyManage() []models.DutyManagement {

	var list []models.DutyManagement
	globals.DBCli.Model(&models.DutyManagement{}).Find(&list)
	return list

}

func (dms *DutyManageService) CreateDutyManage(dutyManage models.DutyManagement) (models.DutyManagement, error) {

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
	err := tx.Model(&models.DutyManagement{}).Where("id = ?", dutyManage.ID).Updates(&dutyManage).Error
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

func (dms *DutyManageService) DeleteDutyManage(id string) error {

	var noticeNum int64
	globals.DBCli.Model(&models.AlertNotice{}).Where("duty_id = ?", id).Count(&noticeNum)
	if noticeNum != 0 {
		return fmt.Errorf("无法删除值班表 %s, 因为已有通知对象绑定", id)
	}

	tx := globals.DBCli.Begin()
	err := tx.Where("id = ?", id).Delete(&models.DutyManagement{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&models.DutySchedule{}).Where("duty_id = ?", id).Delete(&models.DutySchedule{}).Error
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

func (dms *DutyManageService) GetDutyManage(id string) models.DutyManagement {

	var data models.DutyManagement
	globals.DBCli.Model(&models.DutyManagement{}).Where("id = ?", id).Find(&data)
	return data

}
