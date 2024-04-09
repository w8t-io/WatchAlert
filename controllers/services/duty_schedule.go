package services

import (
	"fmt"
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
	SelectDutySystem(tid, dutyId, date string) ([]models.DutySchedule, error)
}

func NewInterDutyScheduleService() InterDutyScheduleService {
	return &DutyScheduleService{}
}

// CreateAndUpdateDutySystem 创建和更新值班表
func (dms *DutyScheduleService) CreateAndUpdateDutySystem(dutyInfo models.DutyScheduleCreate) ([]models.DutySchedule, error) {

	var (
		dutyScheduleList []models.DutySchedule
		timeC            = make(chan string, 370)
		wg               sync.WaitGroup
	)
	// 默认从当前月份顺延到年底
	curYear, curMonth, _ := parseTime(dutyInfo.Month)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 生产值班日期
		for mon := int(curMonth); mon <= 12; mon++ {
			for day := 1; day <= 31; day++ {
				dutyTime := fmt.Sprintf("%d-%d-%d", curYear, mon, day)
				timeC <- dutyTime
			}
		}
		// 产出值班表数据结构
		for len(timeC) != 0 {
			for _, value := range dutyInfo.Users {
				for t := 1; t <= dutyInfo.DutyPeriod; t++ {
					dutyTime := <-timeC
					ds := models.DutySchedule{
						TenantId: dutyInfo.TenantId,
						DutyId:   dutyInfo.DutyId,
						Time:     dutyTime,
						Users: models.Users{
							UserId:   value.UserId,
							Username: value.Username,
						},
					}
					dutyScheduleList = append(dutyScheduleList, ds)
				}
			}
		}
	}()
	wg.Wait()
	close(timeC)

	go func(dutyScheduleList []models.DutySchedule) {
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
	}(dutyScheduleList)

	return dutyScheduleList, nil

}

// UpdateDutySystem 更新值班表
func (dms *DutyScheduleService) UpdateDutySystem(dutySchedule models.DutySchedule) error {

	err := globals.DBCli.Model(&models.DutySchedule{}).Where("tenant_id = ? AND duty_id = ? AND time = ?", dutySchedule.TenantId, dutySchedule.DutyId, dutySchedule.Time).Updates(&dutySchedule).Error
	if err != nil {
		return err
	}
	return nil

}

// SelectDutySystem 查询值班表
func (dms *DutyScheduleService) SelectDutySystem(tid, dutyId, date string) ([]models.DutySchedule, error) {

	var dutyScheduleList []models.DutySchedule
	db := globals.DBCli.Model(&models.DutySchedule{})

	if date != "" {
		db.Where("tenant_id = ? AND duty_id = ? AND time = ?", tid, dutyId, date).Find(&dutyScheduleList)
		return dutyScheduleList, nil
	}

	yearMonth := fmt.Sprintf("%d-%d-", time.Now().Year(), time.Now().Month())

	db.Where("tenant_id = ? AND duty_id = ? AND time LIKE ?", tid, dutyId, yearMonth+"%").Find(&dutyScheduleList)

	return dutyScheduleList, nil

}

func parseTime(month string) (int, time.Month, int) {
	parsedTime, err := time.Parse(layout, month)
	if err != nil {
		return 0, time.Month(0), 0
	}
	curYear, curMonth, curDay := parsedTime.Date()
	return curYear, curMonth, curDay
}
