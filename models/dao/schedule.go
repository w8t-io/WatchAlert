package dao

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

type DutySystem struct {
	Time     string `json:"time"`
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}
