package schedule

import (
	"fmt"
	"log"
	"prometheus-manager/globals"
	"prometheus-manager/models/dao"
	"strconv"
	"time"
)

// CreateAndUpdateDutySystem 创建和更新值班表
func CreateAndUpdateDutySystem(dutyUserInfo []dao.DutySystem, dutyPeriod int) ([]dao.DutySystem, error) {

	var dutySystem []dao.DutySystem
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
				ds := dao.DutySystem{
					Time:     dutyTime,
					UserName: value.UserName,
					UserId:   value.UserId,
				}
				dutySystem = append(dutySystem, ds)
			}
		}
	}

	for _, v := range dutySystem {

		dutySystemInfo, _ := GetCurrentDutyInfo(v.Time)

		if dutySystemInfo.Time != "" {

			if err := UpdateDutySystem(v); err != nil {
				return []dao.DutySystem{}, err
			}

		} else {

			if err := globals.DBCli.Create(&v).Error; err != nil {
				log.Println("值班系统创建失败", err)
				return []dao.DutySystem{}, err
			}

		}

	}

	return dutySystem, nil

}

// UpdateDutySystem 更新值班表
func UpdateDutySystem(dutySystem dao.DutySystem) error {

	err := globals.DBCli.Model(&dao.DutySystem{}).Where("time = ?", dutySystem.Time).Updates(dutySystem).Error
	if err != nil {
		return err
	}
	return nil

}

// SelectDutySystem 查询值班表
func SelectDutySystem(date string) ([]dao.DutySystem, error) {

	var (
		dutySystemList []dao.DutySystem
	)

	if date != "" {
		globals.DBCli.Model(&dao.DutySystem{}).Where("time = ?", date).Find(&dutySystemList)
		return dutySystemList, nil
	}

	yearMonth := fmt.Sprintf("%d-%d-", time.Now().Year(), time.Now().Month())

	globals.DBCli.Model(&dao.DutySystem{}).Where("time LIKE ?", yearMonth+"%").Find(&dutySystemList)

	return dutySystemList, nil

}

// GetCurrentDutyInfo 获取值班信息
func GetCurrentDutyInfo(time string) (dao.DutySystem, dao.People) {

	var (
		dutySystem dao.DutySystem
		dutyPeople dao.People
	)

	globals.DBCli.Where("time = ?", time).Find(&dutySystem)

	globals.DBCli.Where("userName = ?", dutySystem.UserName).Find(&dutyPeople)

	return dutySystem, dutyPeople

}
