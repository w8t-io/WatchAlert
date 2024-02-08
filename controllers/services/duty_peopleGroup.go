package services

import (
	"log"
	"watchAlert/globals"
	"watchAlert/models"
)

type DutyPeopleGroupService struct{}

type InterDutyPeopleGroupService interface {
	CreateDutyGroup(groupInfo models.PeopleGroup) error
}

func NewInterDutyPeopleGroupService() InterDutyPeopleGroupService {
	return &DutyPeopleGroupService{}
}

func (dpgs *DutyPeopleGroupService) CreateDutyGroup(groupInfo models.PeopleGroup) error {

	err := globals.DBCli.Create(&groupInfo).Error
	if err != nil {
		log.Println("值班组创建失败 ->", err)
		return err
	}

	return nil
}
