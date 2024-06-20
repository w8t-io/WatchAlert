package services

import (
	"fmt"
	"sync"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type dutyCalendarService struct {
	ctx *ctx.Context
}

var layout = "2006-01"

type InterDutyCalendarService interface {
	CreateAndUpdate(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
}

func newInterDutyCalendarService(ctx *ctx.Context) InterDutyCalendarService {
	return &dutyCalendarService{
		ctx: ctx,
	}
}

// CreateAndUpdate 创建和更新值班表
func (dms dutyCalendarService) CreateAndUpdate(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyScheduleCreate)

	var (
		dutyScheduleList []models.DutySchedule
		timeC            = make(chan string, 370)
		wg               sync.WaitGroup
	)
	// 默认从当前月份顺延到年底
	curYear, curMonth, _ := parseTime(r.Month)

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
		close(timeC)

		// 产出值班表数据结构
		for {
			if len(timeC) == 0 {
				break
			}

			for _, value := range r.Users {
				for t := 1; t <= r.DutyPeriod; t++ {
					tc := <-timeC
					ds := models.DutySchedule{
						TenantId: r.TenantId,
						DutyId:   r.DutyId,
						Time:     tc,
						Users: models.Users{
							UserId:   value.UserId,
							Username: value.Username,
						},
					}
					dutyScheduleList = append(dutyScheduleList, ds)
				}
			}
		}
		fmt.Println("done")

	}()

	wg.Wait()

	go func(dutyScheduleList []models.DutySchedule) {
		for _, v := range dutyScheduleList {
			// 更新当前已发布的日程表
			dutyScheduleInfo := dms.ctx.DB.DutyCalendar().GetCalendarInfo(r.DutyId, v.Time)
			if dutyScheduleInfo.Time != "" {
				if err := dms.ctx.DB.DutyCalendar().Update(v); err != nil {
					global.Logger.Sugar().Errorf("值班系统更新失败 %s", err)
				}
			} else {
				err := dms.ctx.DB.DutyCalendar().Create(v)
				if err != nil {
					global.Logger.Sugar().Errorf("值班系统创建失败 %s", err)
					return
				}
			}
		}
	}(dutyScheduleList)

	return dutyScheduleList, nil

}

// Update 更新值班表
func (dms dutyCalendarService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutySchedule)
	err := dms.ctx.DB.DutyCalendar().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Search 查询值班表
func (dms dutyCalendarService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DutyScheduleQuery)
	data, err := dms.ctx.DB.DutyCalendar().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseTime(month string) (int, time.Month, int) {
	parsedTime, err := time.Parse(layout, month)
	if err != nil {
		return 0, time.Month(0), 0
	}
	curYear, curMonth, curDay := parsedTime.Date()
	return curYear, curMonth, curDay
}
