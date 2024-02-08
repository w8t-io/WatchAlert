package repo

import (
	"watchAlert/globals"
	"watchAlert/models"
)

type DutyScheduleRepo struct{}

// GetDutyScheduleInfo 获取值班表信息
func (dsr *DutyScheduleRepo) GetDutyScheduleInfo(dutyId, time string) models.DutySchedule {

	var dutySchedule models.DutySchedule

	globals.DBCli.Where("duty_id = ? AND time = ?", dutyId, time).Find(&dutySchedule)

	return dutySchedule

}

// GetDutyUserInfo 获取值班用户信息
func (dsr *DutyScheduleRepo) GetDutyUserInfo(dutyId, time string) models.Member {

	var user models.Member

	schedule := dsr.GetDutyScheduleInfo(dutyId, time)
	globals.DBCli.Model(&models.Member{}).Where("user_id = ?", schedule.UserId).First(&user)

	return user

}
