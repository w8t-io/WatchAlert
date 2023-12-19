package services

import (
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
	"prometheus-manager/utils/cmd"
)

type DutyManageService struct{}

type InterDutyManageService interface {
	ListDutyManage() []dao.DutyManagement
	CreateDutyManage(dutyManage dao.DutyManagement) (dao.DutyManagement, error)
	UpdateDutyManage(dutyManage dao.DutyManagement) (dao.DutyManagement, error)
	DeleteDutyManage(id string) error
	GetDutyManage(id string) dao.DutyManagement
}

func NewInterDutyManageService() InterDutyManageService {
	return &DutyManageService{}
}

func (dms *DutyManageService) ListDutyManage() []dao.DutyManagement {

	var list []dao.DutyManagement
	globals.DBCli.Model(&dao.DutyManagement{}).Find(&list)
	return list

}

func (dms *DutyManageService) CreateDutyManage(dutyManage dao.DutyManagement) (dao.DutyManagement, error) {

	tx := globals.DBCli.Begin()
	dutyManage.ID = cmd.RandUuid()
	err := tx.Create(&dutyManage).Error
	if err != nil {
		tx.Rollback()
		return dao.DutyManagement{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return dao.DutyManagement{}, err
	}
	return dutyManage, nil

}

func (dms *DutyManageService) UpdateDutyManage(dutyManage dao.DutyManagement) (dao.DutyManagement, error) {

	tx := globals.DBCli.Begin()
	err := tx.Model(&dao.DutyManagement{}).Where("id = ?", dutyManage.ID).Updates(&dutyManage).Error
	if err != nil {
		tx.Rollback()
		return dao.DutyManagement{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return dao.DutyManagement{}, err
	}
	return dutyManage, nil

}

func (dms *DutyManageService) DeleteDutyManage(id string) error {

	tx := globals.DBCli.Begin()
	err := tx.Where("id = ?", id).Delete(&dao.DutyManagement{}).Error
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

func (dms *DutyManageService) GetDutyManage(id string) dao.DutyManagement {

	var data dao.DutyManagement
	globals.DBCli.Model(&dao.DutyManagement{}).Where("id = ?", id).Find(&data)
	return data

}
