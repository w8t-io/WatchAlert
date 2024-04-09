package repo

import (
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
)

type Event struct{}

func (e Event) GetHistoryEvent(tid, datasourceType, severity string, startAt, endAt, pageIndex, pageSize int64) (response.HistoryEvent, error) {
	var data []models.AlertHisEvent
	var count int64

	db := globals.DBCli.Model(&models.AlertHisEvent{})
	db.Where("tenant_id = ?", tid)

	if datasourceType != "" {
		db = db.Where("datasource_type = ?", datasourceType)
	}

	if severity != "" {
		db = db.Where("severity = ?", severity)
	}

	if startAt != 0 && endAt != 0 {
		db = db.Where("first_trigger_time > ? and first_trigger_time < ?", startAt, endAt)
	}

	if err := db.Count(&count).Error; err != nil {
		return response.HistoryEvent{}, err
	}

	if err := db.Limit(int(pageSize)).Offset(int((pageIndex - 1) * pageSize)).Order("recover_time desc").Find(&data).Error; err != nil {
		return response.HistoryEvent{}, err
	}

	return response.HistoryEvent{
		List:       data,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: count,
	}, nil
}
