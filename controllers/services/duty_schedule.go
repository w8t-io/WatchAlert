package services

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
)

type DutyScheduleService struct{}

type InterDutyScheduleService interface {
	CreateAndUpdateDutySystem(dutyUserInfo models.DutyScheduleCreate) ([]models.DutySchedule, error)
	UpdateDutySystem(dutySchedule models.DutySchedule) error
	SelectDutySystem(dutyId, date string) ([]models.DutySchedule, error)
}

func NewInterDutyScheduleService() InterDutyScheduleService {
	return &DutyScheduleService{}
}

// CreateAndUpdateDutySystem 创建和更新值班表
func (dms *DutyScheduleService) CreateAndUpdateDutySystem(dutyInfo models.DutyScheduleCreate) ([]models.DutySchedule, error) {

	var dutyScheduleList []models.DutySchedule
	ch := make(chan string)

	layout := "2006-01"
	parsedTime, err := time.Parse(layout, dutyInfo.Month)
	if err != nil {
		return nil, err
	}

	// 获取当前月份
	year, month, _ := parsedTime.Date()
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
		for _, value := range dutyInfo.Users {
			for t := 1; t <= dutyInfo.DutyPeriod; t++ {
				dutyTime := <-ch
				if dutyTime == "" {
					break
				}
				ds := models.DutySchedule{
					DutyId: dutyInfo.DutyId,
					Time:   dutyTime,
					Users: models.Users{
						UserId:   value.UserId,
						Username: value.Username,
					},
				}
				dutyScheduleList = append(dutyScheduleList, ds)
			}
		}
	}

	for _, v := range dutyScheduleList {

		// 更新当前已发布的日程表
		dutyScheduleInfo := dutySchedule.GetDutyScheduleInfo(dutyInfo.DutyId, v.Time)

		if dutyScheduleInfo.Time != "" {

			if err = dms.UpdateDutySystem(v); err != nil {
				return nil, err
			}

		} else {

			if err = globals.DBCli.Create(&v).Error; err != nil {
				log.Println("值班系统创建失败", err)
				return nil, err
			}

		}

	}

	return dutyScheduleList, nil

}

// UpdateDutySystem 更新值班表
func (dms *DutyScheduleService) UpdateDutySystem(dutySchedule models.DutySchedule) error {

	err := globals.DBCli.Model(&models.DutySchedule{}).Where("duty_id = ? AND time = ?", dutySchedule.DutyId, dutySchedule.Time).Updates(&dutySchedule).Error
	if err != nil {
		return err
	}
	return nil

}

// SelectDutySystem 查询值班表
func (dms *DutyScheduleService) SelectDutySystem(dutyId, date string) ([]models.DutySchedule, error) {

	var (
		dutyScheduleList []models.DutySchedule
	)

	if date != "" {
		globals.DBCli.Model(&models.DutySchedule{}).Where("duty_id = ? AND time = ?", dutyId, date).Find(&dutyScheduleList)
		return dutyScheduleList, nil
	}

	yearMonth := fmt.Sprintf("%d-%d-", time.Now().Year(), time.Now().Month())

	globals.DBCli.Model(&models.DutySchedule{}).Where("duty_id = ? AND time LIKE ?", dutyId, yearMonth+"%").Find(&dutyScheduleList)

	return dutyScheduleList, nil

}
