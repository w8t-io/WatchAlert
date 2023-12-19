package services

import (
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
	"prometheus-manager/utils/cmd"
)

type DutyPeopleService struct{}

type InterDutyPeopleService interface {
	CreateDutyUser(userInfo dao.People) (dao.People, error)
	SelectDutyUser() []dao.People
	UpdateDutyUser(userInfo dao.People) (dao.People, error)
	DeleteDutyUser(userId string) error
	GetDutyUser(user string) ([]dao.People, error)
}

func NewInterDutyPeopleService() InterDutyPeopleService {
	return &DutyPeopleService{}
}

func (dps *DutyPeopleService) CreateDutyUser(userInfo dao.People) (dao.People, error) {

	userInfo.UserID = cmd.RandUserId()
	err := globals.DBCli.Create(&userInfo).Error
	if err != nil {
		globals.Logger.Sugar().Error("值班用户创建失败 ->", err)
		return dao.People{}, err
	}

	return userInfo, nil
}

func (dps *DutyPeopleService) SelectDutyUser() []dao.People {

	var people []dao.People

	err := globals.DBCli.Model(&dao.People{}).Find(&people).Error
	if err != nil {
		globals.Logger.Sugar().Error("用户查询失败失败 ->", err)
		return nil
	}

	return people

}

func (dps *DutyPeopleService) UpdateDutyUser(userInfo dao.People) (dao.People, error) {

	var newInfo dao.People

	tx := globals.DBCli.Begin()
	err := tx.Where("userId = ?", userInfo.UserID).Updates(&userInfo).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Error("更新用户信息失败")
		return userInfo, err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		globals.Logger.Error("提交事务失败")
		return userInfo, err
	}

	globals.DBCli.Model(&dao.People{}).Where("userId = ?", userInfo.UserID).Find(&newInfo)

	return newInfo, nil

}

func (dps *DutyPeopleService) DeleteDutyUser(userId string) error {

	tx := globals.DBCli.Begin()
	err := tx.Where("userId = ?", userId).Delete(&dao.People{}).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Error("删除用户信息失败")
		return err
	}

	err = tx.Exec("update duty_systems set user_id = '', user_name = '' where user_id = ?", userId).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Error("更新值班表失败")
		return err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		globals.Logger.Error("提交事务失败")
		return err
	}

	return nil

}

func (dps *DutyPeopleService) GetDutyUser(user string) ([]dao.People, error) {

	var userInfo []dao.People

	err := globals.DBCli.Model(dao.People{}).Where("userId LIKE ? OR userName LIKE ? OR phone LIKE ? OR email LIKE ?",
		"%"+user+"%", "%"+user+"%", "%"+user+"%", "%"+user+"%").Find(&userInfo).Error
	if err != nil {
		globals.Logger.Sugar().Error("值班用户查询失败 ->", err)
		return []dao.People{}, err
	}

	return userInfo, nil

}
