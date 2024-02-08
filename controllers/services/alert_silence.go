package services

import (
	"fmt"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	models2 "watchAlert/models"
	"watchAlert/utils/cmd"
)

type AlertSilenceService struct {
	alertEvent models2.AlertCurEvent
}

type InterAlertSilenceService interface {
	CreateAlertSilence(silence models2.AlertSilences) error
	UpdateAlertSilence(silence models2.AlertSilences) (models2.AlertSilences, error)
	DeleteAlertSilence(id string) error
	ListAlertSilence() ([]models2.AlertSilences, error)
}

func NewInterAlertSilenceService() InterAlertSilenceService {
	return &AlertSilenceService{}
}

func (ass *AlertSilenceService) CreateAlertSilence(silence models2.AlertSilences) error {

	createAt := time.Now().Unix()
	silenceEvent := models2.AlertSilences{
		Id:             "s-" + cmd.RandId(),
		Fingerprint:    silence.Fingerprint,
		Datasource:     silence.Datasource,
		DatasourceType: silence.DatasourceType,
		StartsAt:       silence.StartsAt,
		EndsAt:         silence.EndsAt,
		CreateBy:       silence.CreateBy,
		CreateAt:       createAt,
		UpdateAt:       createAt,
		Comment:        silence.Comment,
	}

	event, ok := silenceEvent.GetCache(silence.Fingerprint)
	if ok && event != "" {
		return fmt.Errorf("静默消息已存在, ID:%s", silenceEvent.Id)
	}

	muteAt := silence.EndsAt - createAt
	duration := time.Duration(muteAt) * time.Second
	silenceEvent.SetCache(duration)

	err := repo.DBCli.Create(models2.AlertSilences{}, silenceEvent)
	if err != nil {
		return err
	}

	return nil

}

func (ass *AlertSilenceService) UpdateAlertSilence(silence models2.AlertSilences) (models2.AlertSilences, error) {

	updateAt := time.Now().Unix()

	silence.UpdateAt = updateAt
	muteAt := silence.EndsAt - silence.StartsAt
	duration := time.Duration(muteAt) * time.Second
	silence.SetCache(duration)

	err := repo.DBCli.Updates(repo.Updates{
		Table:   models2.AlertSilences{},
		Where:   []string{"id = ?", silence.Id},
		Updates: silence,
	})

	if err != nil {
		return models2.AlertSilences{}, err
	}

	return silence, nil

}

func (ass *AlertSilenceService) DeleteAlertSilence(id string) error {

	var silence models2.AlertSilences
	globals.DBCli.Where("id = ?", id).Find(&silence)

	del := repo.Delete{
		Table: models2.AlertSilences{},
		Where: []string{"id = ?", id},
	}
	repo.DBCli.Delete(del)

	_, err := globals.RedisCli.Del(models2.SilenceCachePrefix + silence.Fingerprint).Result()
	if err != nil {
		return err
	}
	return nil

}

func (ass *AlertSilenceService) ListAlertSilence() ([]models2.AlertSilences, error) {

	var silenceList []models2.AlertSilences
	err := globals.DBCli.Find(&silenceList).Error

	if err != nil {
		return []models2.AlertSilences{}, err
	}

	return silenceList, nil

}
