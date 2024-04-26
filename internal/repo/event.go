package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	EventRepo struct {
		entryRepo
	}

	InterEventRepo interface {
		GetHistoryEvent(r models.AlertHisEventQuery) (models.HistoryEventResponse, error)
		CreateHistoryEvent(r models.AlertHisEvent) error
	}
)

func newEventInterface(db *gorm.DB, g InterGormDBCli) InterEventRepo {
	return &EventRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (e EventRepo) GetHistoryEvent(r models.AlertHisEventQuery) (models.HistoryEventResponse, error) {
	var data []models.AlertHisEvent
	var count int64

	db := e.DB().Model(&models.AlertHisEvent{})
	db.Where("tenant_id = ?", r.TenantId)

	if r.DatasourceType != "" {
		db = db.Where("datasource_type = ?", r.DatasourceType)
	}

	if r.Severity != "" {
		db = db.Where("severity = ?", r.Severity)
	}

	if r.StartAt != 0 && r.EndAt != 0 {
		db = db.Where("first_trigger_time > ? and first_trigger_time < ?", r.StartAt, r.EndAt)
	}

	if err := db.Count(&count).Error; err != nil {
		return models.HistoryEventResponse{}, err
	}

	if err := db.Limit(int(r.PageSize)).Offset(int((r.PageIndex - 1) * r.PageSize)).Order("recover_time desc").Find(&data).Error; err != nil {
		return models.HistoryEventResponse{}, err
	}

	return models.HistoryEventResponse{
		List:       data,
		PageIndex:  r.PageIndex,
		PageSize:   r.PageSize,
		TotalCount: count,
	}, nil
}

func (e EventRepo) CreateHistoryEvent(r models.AlertHisEvent) error {
	err := e.g.Create(models.AlertHisEvent{}, r)
	if err != nil {
		return err
	}

	return nil
}
