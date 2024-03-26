package repo

import (
	"watchAlert/globals"
	"watchAlert/models"
)

type Event struct{}

func (e Event) GetHistoryEvent(datasourceType, severity string, startAt, endAt, pageIndex, pageSize int64) ([]models.AlertHisEvent, error) {
	var data []models.AlertHisEvent
	db := globals.DBCli.Model(&models.AlertHisEvent{})

	if datasourceType != "" {
		db = db.Where("datasource_type = ?", datasourceType)
	}

	if severity != "" {
		db = db.Where("severity = ?", severity)
	}

	if startAt != 0 && endAt != 0 {
		db = db.Where("first_trigger_time > ? and first_trigger_time < ?", startAt, endAt)
	}

	if err := db.Limit(int(pageSize)).Offset(int((pageIndex - 1) * pageSize)).Order("recover_time desc").Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (e Event) CountHistoryEvent() (int64, error) {
	var count int64
	db := globals.DBCli.Model(&models.AlertHisEvent{})

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
