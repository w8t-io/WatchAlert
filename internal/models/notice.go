package models

import (
	"gorm.io/gorm"
)

type AlertNotice struct {
	TenantId        string `json:"tenantId"`
	Uuid            string `json:"uuid"`
	Name            string `json:"name"`
	Env             string `json:"env"`
	DutyId          string `json:"dutyId"`
	NoticeType      string `json:"noticeType"`
	EnableCard      string `json:"enableCard"`
	Hook            string `json:"hook"`
	Email           Email  `json:"email" gorm:"email;serializer:json"`
	Template        string `json:"template"`
	TemplateFiring  string `json:"templateFiring"`
	TemplateRecover string `json:"templateRecover"`
}

type Email struct {
	Subject string   `json:"subject"`
	To      []string `json:"to" gorm:"to;serializer:json"`
	CC      []string `json:"cc" gorm:"cc;serializer:json"`
}

type AlertRecord struct {
	gorm.Model
	AlertName   string `json:"alertName"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

type NoticeTemplateExample struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Template    string `json:"template"`
}

type NoticeQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	Uuid     string `json:"uuid" form:"uuid"`
	Name     string `json:"name" form:"name"`
	Query    string `json:"query" form:"query"`
}

type NoticeTemplateExampleQuery struct {
	Id    string `json:"id" form:"id"`
	Name  string `json:"name" form:"name"`
	Query string `json:"query" form:"query"`
}
