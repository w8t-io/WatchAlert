package services

import (
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type DutyPeopleService struct{}

type InterDutyPeopleService interface {
	CreateDutyUser(userInfo models.People) (models.People, error)
	SelectDutyUser() []models.People
	UpdateDutyUser(userInfo models.People) (models.People, error)
	DeleteDutyUser(userId string) error
	GetDutyUser(user string) ([]models.People, error)
}

func NewInterDutyPeopleService() InterDutyPeopleService {
	return &DutyPeopleService{}
}

func (dps *DutyPeopleService) CreateDutyUser(userInfo models.People) (models.People, error) {

	userInfo.UserID = cmd.RandId()
	err := globals.DBCli.Create(&userInfo).Error
	if err != nil {
		globals.Logger.Sugar().Error("值班用户创建失败 ->", err)
		return models.People{}, err
	}

	return userInfo, nil
}

func (dps *DutyPeopleService) SelectDutyUser() []models.People {

	var people []models.People

	err := globals.DBCli.Model(&models.People{}).Find(&people).Error
	if err != nil {
		globals.Logger.Sugar().Error("用户查询失败失败 ->", err)
		return nil
	}

	return people

}

func (dps *DutyPeopleService) UpdateDutyUser(userInfo models.People) (models.People, error) {

	var newInfo models.People

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

	globals.DBCli.Model(&models.People{}).Where("userId = ?", userInfo.UserID).Find(&newInfo)

	return newInfo, nil

}

func (dps *DutyPeopleService) DeleteDutyUser(userId string) error {

	tx := globals.DBCli.Begin()
	err := tx.Where("userId = ?", userId).Delete(&models.People{}).Error
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

func (dps *DutyPeopleService) GetDutyUser(user string) ([]models.People, error) {

	var userInfo []models.People

	err := globals.DBCli.Model(models.People{}).Where("userId LIKE ? OR userName LIKE ? OR phone LIKE ? OR email LIKE ?",
		"%"+user+"%", "%"+user+"%", "%"+user+"%", "%"+user+"%").Find(&userInfo).Error
	if err != nil {
		globals.Logger.Sugar().Error("值班用户查询失败 ->", err)
		return []models.People{}, err
	}

	return userInfo, nil

}
