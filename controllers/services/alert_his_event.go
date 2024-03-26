package services

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type AlertHisEventService struct {
	repo.Event
}

type InterAlertHisEventService interface {
	List(ctx *gin.Context) (response.HistoryEvent, error)
	Search() []models.AlertHisEvent
}

func NewInterAlertHisEventService() InterAlertHisEventService {
	return &AlertHisEventService{}
}

func (ahes *AlertHisEventService) List(ctx *gin.Context) (response.HistoryEvent, error) {

	datasourceType := ctx.Query("datasourceType")
	severity := ctx.Query("severity")
	startAt := ctx.Query("startAt")
	endAt := ctx.Query("endAt")
	pageIndex, _ := strconv.ParseInt(ctx.Query("pageIndex"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.Query("pageSize"), 10, 64)

	var (
		startAtInt64 int64
		endAtInt64   int64
		err          error
	)

	if startAt != "" && endAt != "" {
		startAtInt64, err = strconv.ParseInt(startAt, 10, 64)
		if err != nil {
			return response.HistoryEvent{}, err
		}

		endAtInt64, err = strconv.ParseInt(endAt, 10, 64)
		if err != nil {
			return response.HistoryEvent{}, err
		}
	}

	data, err := ahes.GetHistoryEvent(datasourceType, severity, startAtInt64, endAtInt64, pageIndex, pageSize)
	count, err := ahes.CountHistoryEvent()

	return response.HistoryEvent{
		List:       data,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalCount: count,
	}, err

}

func (ahes *AlertHisEventService) Search() []models.AlertHisEvent {
	return nil
}
