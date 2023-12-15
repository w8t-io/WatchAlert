package schedule

import (
	"prometheus-manager/globals"
	"prometheus-manager/models/dao"
	"prometheus-manager/utils"
)

func CreateDutyUser(userInfo dao.People) (dao.People, error) {

	userInfo.UserID = utils.RandUserId()
	err := globals.DBCli.Create(&userInfo).Error
	if err != nil {
		globals.Logger.Sugar().Error("值班用户创建失败 ->", err)
		return dao.People{}, err
	}

	return userInfo, nil
}

func SelectDutyUser() []dao.People {

	var people []dao.People

	err := globals.DBCli.Model(&dao.People{}).Find(&people).Error
	if err != nil {
		globals.Logger.Sugar().Error("用户查询失败失败 ->", err)
		return nil
	}

	return people

}

func UpdateDutyUser(userInfo dao.People) (dao.People, error) {

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

func DeleteDutyUser(userId string) error {

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

func GetDutyUser(user string) ([]dao.People, error) {

	var userInfo []dao.People

	err := globals.DBCli.Model(dao.People{}).Where("userId LIKE ? OR userName LIKE ? OR phone LIKE ? OR email LIKE ?",
		"%"+user+"%", "%"+user+"%", "%"+user+"%", "%"+user+"%").Find(&userInfo).Error
	if err != nil {
		globals.Logger.Sugar().Error("值班用户查询失败 ->", err)
		return []dao.People{}, err
	}

	return userInfo, nil

}
