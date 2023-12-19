package services

import (
	"fmt"
	"log"
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
	"strconv"
	"time"
)

type DutyScheduleService struct{}

type InterDutyScheduleService interface {
	CreateAndUpdateDutySystem(dutyUserInfo []dao.DutySchedule, dutyPeriod int, dutyId string) ([]dao.DutySchedule, error)
	UpdateDutySystem(dutySchedule dao.DutySchedule, dutyId string) error
	SelectDutySystem(dutyId, date string) ([]dao.DutySchedule, error)
}

func NewInterDutyScheduleService() InterDutyScheduleService {
	return &DutyScheduleService{}
}

// CreateAndUpdateDutySystem 创建和更新值班表
func (dms *DutyScheduleService) CreateAndUpdateDutySystem(dutyUserInfo []dao.DutySchedule, dutyPeriod int, dutyId string) ([]dao.DutySchedule, error) {

	var dutyScheduleList []dao.DutySchedule
	ch := make(chan string)

	// 获取当前月份
	year, month, _ := time.Now().Date()
	// 构建下个月的第一天
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	// 计算当前月份有多少天
	daysInMonth := nextMonth.Add(-time.Hour * 24).Day()

	go func() {
		for i := 1; i <= daysInMonth; i++ {
			dutyTime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(i)
			ch <- dutyTime
		}
		defer close(ch)
	}()

	for i := 0; i <= daysInMonth; i++ {
		for _, value := range dutyUserInfo {
			for t := 1; t <= dutyPeriod; t++ {
				dutyTime := <-ch
				if dutyTime == "" {
					break
				}
				ds := dao.DutySchedule{
					DutyId:   dutyId,
					Time:     dutyTime,
					UserName: value.UserName,
					UserId:   value.UserId,
				}
				dutyScheduleList = append(dutyScheduleList, ds)
			}
		}
	}

	for _, v := range dutyScheduleList {

		dutyScheduleInfo, _ := dutySchedule.GetDutyScheduleInfo(dutyId, v.Time)

		if dutyScheduleInfo.Time != "" {

			if err := dms.UpdateDutySystem(v, dutyId); err != nil {
				return []dao.DutySchedule{}, err
			}

		} else {

			if err := globals.DBCli.Create(&v).Error; err != nil {
				log.Println("值班系统创建失败", err)
				return []dao.DutySchedule{}, err
			}

		}

	}

	return dutyScheduleList, nil

}

// UpdateDutySystem 更新值班表
func (dms *DutyScheduleService) UpdateDutySystem(dutySchedule dao.DutySchedule, dutyId string) error {

	err := globals.DBCli.Model(&dao.DutySchedule{}).Where("duty_id = ? AND time = ?", dutyId, dutySchedule.Time).Updates(&dutySchedule).Error
	if err != nil {
		return err
	}
	return nil

}

// SelectDutySystem 查询值班表
func (dms *DutyScheduleService) SelectDutySystem(dutyId, date string) ([]dao.DutySchedule, error) {

	var (
		dutyScheduleList []dao.DutySchedule
	)

	if date != "" {
		globals.DBCli.Model(&dao.DutySchedule{}).Where("duty_id = ? AND time = ?", dutyId, date).Find(&dutyScheduleList)
		return dutyScheduleList, nil
	}

	yearMonth := fmt.Sprintf("%d-%d-", time.Now().Year(), time.Now().Month())

	globals.DBCli.Model(&dao.DutySchedule{}).Where("duty_id = ? AND time LIKE ?", dutyId, yearMonth+"%").Find(&dutyScheduleList)

	return dutyScheduleList, nil

}
