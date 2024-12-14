package repo

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"watchAlert/internal/models"
)

type (
	DutyCalendarRepo struct {
		entryRepo
	}

	InterDutyCalendar interface {
		GetCalendarInfo(dutyId, time string) models.DutySchedule
		GetDutyUserInfo(dutyId, time string) (models.Member, bool)
		Create(r models.DutySchedule) error
		Update(r models.DutySchedule) error
		Search(r models.DutyScheduleQuery) ([]models.DutySchedule, error)
		GetCalendarUsers(r models.DutyScheduleQuery) ([]models.Users, error)
	}
)

func newDutyCalendarInterface(db *gorm.DB, g InterGormDBCli) InterDutyCalendar {
	return &DutyCalendarRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// GetCalendarInfo 获取值班表信息
func (dc DutyCalendarRepo) GetCalendarInfo(dutyId, time string) models.DutySchedule {
	var dutySchedule models.DutySchedule

	dc.db.Model(models.DutySchedule{}).
		Where("duty_id = ? AND time = ?", dutyId, time).
		First(&dutySchedule)

	return dutySchedule
}

// GetDutyUserInfo 获取值班用户信息
func (dc DutyCalendarRepo) GetDutyUserInfo(dutyId, time string) (models.Member, bool) {
	var user models.Member
	schedule := dc.GetCalendarInfo(dutyId, time)
	db := dc.db.Model(models.Member{}).
		Where("user_id = ?", schedule.UserId)
	if err := db.First(&user).Error; err != nil {
		return user, false
	}
	if user.JoinDuty == "true" {
		return user, true
	}

	return user, false
}

func (dc DutyCalendarRepo) Create(r models.DutySchedule) error {
	err := dc.g.Create(models.DutySchedule{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (dc DutyCalendarRepo) Update(r models.DutySchedule) error {
	u := Updates{
		Table: models.DutySchedule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"duty_id = ?":   r.DutyId,
			"time = ?":      r.Time,
		},
		Updates: r,
	}

	err := dc.g.Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (dc DutyCalendarRepo) Search(r models.DutyScheduleQuery) ([]models.DutySchedule, error) {
	var dutyScheduleList []models.DutySchedule
	db := dc.db.Model(&models.DutySchedule{})

	if r.Time != "" {
		db.Where("tenant_id = ? AND duty_id = ? AND time = ?", r.TenantId, r.DutyId, r.Time).Find(&dutyScheduleList)
		return dutyScheduleList, nil
	}

	yearMonth := fmt.Sprintf("%d-%d-", time.Now().Year(), time.Now().Month())
	db.Where("tenant_id = ? AND duty_id = ? AND time LIKE ?", r.TenantId, r.DutyId, yearMonth+"%")
	err := db.Find(&dutyScheduleList).Error
	if err != nil {
		return dutyScheduleList, err
	}

	return dutyScheduleList, nil
}

func (dc DutyCalendarRepo) GetCalendarUsers(r models.DutyScheduleQuery) ([]models.Users, error) {
	var users []models.Users
	db := dc.db.Model(&models.DutySchedule{})
	db.Where("tenant_id = ? AND duty_id = ?", r.TenantId, r.DutyId)
	db.Distinct("user_id", "username")
	if err := db.Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return users, nil
}
