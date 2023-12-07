package schedule

import (
	"log"
	"prometheus-manager/globals"
	"prometheus-manager/models/dao"
)

func CreateDutyGroup(groupInfo dao.PeopleGroup) error {

	err := globals.DBCli.Create(&groupInfo).Error
	if err != nil {
		log.Println("值班组创建失败 ->", err)
		return err
	}

	return nil
}
