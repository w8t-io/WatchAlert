package services

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
)

type DutyScheduleService struct{}

var layout = "2006-01"

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

	var (
		dutyScheduleList []models.DutySchedule
		wg               sync.WaitGroup
	)
	ch := make(chan string)

	go func() {
		for dutyTime := range ch {
			if dutyTime == "" {
				continue
			}
			for _, value := range dutyInfo.Users {
				for t := 1; t <= dutyInfo.DutyPeriod; t++ {
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
				if err := dms.UpdateDutySystem(v); err != nil {
					globals.Logger.Sugar().Errorf("值班系统更新失败 %s", err)
				}
			} else {
				if err := globals.DBCli.Create(&v).Error; err != nil {
					globals.Logger.Sugar().Errorf("值班系统创建失败 %s", err)
				}
			}
		}
	}()

	// 默认从当前月份顺延到年底
	curYear, curMonth, _ := parseTime(dutyInfo.Month)
	wg.Add(1)
	go func(curYear int) {
		defer wg.Done()
		for i := int(curMonth); i <= 12; i++ {
			t := getDaysInMonth(curYear, time.Month(i))
			daysInMonth, _ := getDayNumber(t.Format(layout))
			fmt.Println(t, daysInMonth)
			for i := 1; i <= daysInMonth; i++ {
				dutyTime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(i)
				ch <- dutyTime
			}
		}
	}(curYear)
	wg.Wait()
	close(ch)

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

// 获取月份的天数
func getDayNumber(month string) (int, error) {

	// 获取当前年月份
	curYear, curMonth, _ := parseTime(month)
	// 构建下个月的第一天
	nextMonth := getDaysInMonth(curYear, curMonth)
	// 计算当前月份有多少天
	daysInMonth := nextMonth.Add(-time.Hour * 24).Day()
	return daysInMonth, nil
}

func parseTime(month string) (int, time.Month, int) {
	parsedTime, err := time.Parse(layout, month)
	if err != nil {
		return 0, time.Month(0), 0
	}
	curYear, curMonth, curDay := parsedTime.Date()
	return curYear, curMonth, curDay
}

func getDaysInMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}
