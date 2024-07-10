package models

import "watchAlert/config"

type Settings struct {
	IsInit      int                `json:"isInit"`
	AlarmConfig config.AlarmConfig `json:"alarmConfig" gorm:"alarmConfig;serializer:json"`
	EmailConfig emailConfig        `json:"emailConfig" gorm:"emailConfig;serializer:json"`
	AppVersion  string             `json:"appVersion" gorm:"-"`
}

type emailConfig struct {
	ServerAddress string `json:"serverAddress"`
	Email         string `json:"email"`
	Token         string `json:"token"`
}
