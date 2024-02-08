package models

import (
	"gorm.io/gorm"
)

type People struct {
	gorm.Model
	UserID         string `gorm:"column:userId" json:"userId"`
	UserName       string `gorm:"column:userName" json:"userName"`
	Phone          string `gorm:"column:phone" json:"phone"`
	Email          string `gorm:"column:email" json:"email"`
	Notice         string `gorm:"column:notice" json:"notice"`
	FeiShuUserID   string `gorm:"column:feiShuUserID" json:"feiShuUserID"`
	DingDingUserID string `gorm:"column:dingDingUserID" json:"dingDingUserID"`
}

type PeopleGroup struct {
	gorm.Model
	GroupID   uint   `gorm:"column:groupID" json:"groupID"`
	GroupName string `gorm:"column:groupName" json:"groupName"`
}

type JoinsPeopleGroup struct {
	UserName  string
	GroupName string
}

type DutyManagement struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateBy    string `json:"create_by"`
	CreateAt    int64  `json:"create_at"`
}

type DutyScheduleCreate struct {
	DutyId     string  `json:"dutyId"`
	DutyPeriod int     `json:"dutyPeriod"`
	Month      string  `json:"month"`
	Users      []Users `json:"users"`
}

type Users struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
}

type DutySchedule struct {
	DutyId string `json:"dutyId"`
	Time   string `json:"time"`
	Users
}