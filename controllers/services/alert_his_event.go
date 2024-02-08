package services

import (
	"encoding/json"
	"watchAlert/globals"
	"watchAlert/models"
)

type AlertHisEventService struct {
}

type InterAlertHisEventService interface {
	List() ([]models.AlertHisEvent, error)
	Search() []models.AlertHisEvent
}

func NewInterAlertHisEventService() InterAlertHisEventService {
	return &AlertHisEventService{}
}

func (ahes *AlertHisEventService) List() ([]models.AlertHisEvent, error) {

	var data   []models.AlertHisEvent

	err := globals.DBCli.Find(&data).Error
	if err != nil {
		return nil, err
	}

	for k,v:=range data{
		var metric map[string]string
		_ = json.Unmarshal([]byte(v.Metric),&metric)
		data[k].MetricMap = metric
	}

	return data, err

}

func (ahes *AlertHisEventService) Search() []models.AlertHisEvent {
	return nil
}
