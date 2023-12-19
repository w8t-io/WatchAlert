package repo

import (
	"fmt"
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
)

type DutyScheduleRepo struct{}

// GetDutyScheduleInfo 获取值班信息
func (dsr *DutyScheduleRepo) GetDutyScheduleInfo(dutyId, time string) (dao.DutySchedule, string) {

	var (
		dutySchedule dao.DutySchedule
		dutyPeople   dao.People
	)

	globals.DBCli.Where("duty_id = ? AND time = ?", dutyId, time).Find(&dutySchedule)

	globals.DBCli.Where("userName = ?", dutySchedule.UserName).Find(&dutyPeople)

	if len(dutyPeople.FeiShuUserID) == 0 {
		dutyPeople.FeiShuUserID = "暂无安排值班人员"
	} else {
		dutyPeople.FeiShuUserID = fmt.Sprintf("**👤 值班人员：**<at id=%s></at>", dutyPeople.FeiShuUserID)
	}

	return dutySchedule, dutyPeople.FeiShuUserID

}
