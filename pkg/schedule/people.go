package schedule

import (
	"log"
	"prometheus-manager/globals"
	"prometheus-manager/models/dao"
	"prometheus-manager/utils"
)

func CreateDutyUser(userInfo dao.People) error {

	userInfo.UserID = utils.RandUserId()
	err := globals.DBCli.Create(&userInfo).Error
	if err != nil {
		log.Println("值班用户创建失败 ->", err)
		return err
	}

	return nil
}

func SelectDutyUser() []dao.People {

	var people []dao.People

	err := globals.DBCli.Model(&dao.People{}).Find(&people).Error
	if err != nil {
		log.Println("用户查询失败失败 ->", err)
		return nil
	}

	return people

}

func GetDutyUser(user string) ([]dao.People, error) {

	var userInfo []dao.People

	err := globals.DBCli.Model(dao.People{}).Where("userId LIKE ? OR userName LIKE ? OR phone LIKE ? OR email LIKE ?", "%"+user+"%", "%"+user+"%", "%"+user+"%", "%"+user+"%").Find(&userInfo).Error
	if err != nil {
		log.Println("值班用户查询失败 ->", err)
		return []dao.People{}, err
	}

	return userInfo, nil

}
