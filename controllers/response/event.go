package response

import "watchAlert/models"

type HistoryEvent struct {
	List       []models.AlertHisEvent
	PageIndex  int64
	PageSize   int64
	TotalCount int64
}
