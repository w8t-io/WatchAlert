package utils

import (
	"prometheus-manager/globals"
	"time"
)

// Person user list
type Person struct {
	ActionUser []string
}

// 创建排班表
var schedule = make(map[interface{}]string)

func CreateAndReturnSchedule(date string) string {

	people := Person{ActionUser: globals.Config.FeiShu.DutyUser}

	// 设置值班周期和总共排班天数
	totalDays := 30

	// 获取当前日期
	currentDate := time.Now()

	// 初始化排班表，按照顺序循环排班
	rotationIndex := 0

	for day := 0; day < totalDays; day++ {
		person := people.ActionUser[rotationIndex]
		date := currentDate.AddDate(0, 0, day)
		schedule[date.Format("2006-01-02")] = person

		rotationIndex = (rotationIndex + 1) % len(people.ActionUser)
	}

	return schedule[date]

}

func GetCurrentDutyUser() string {

	timeNow := time.Now().Format("2006-01-02")
	user := schedule[timeNow]
	if user != "" {
		return user
	}

	return CreateAndReturnSchedule(timeNow)

}
