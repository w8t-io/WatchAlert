package services

import (
	"fmt"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type AlertNoticeService struct{}

type InterAlertNoticeService interface {
	SearchNoticeObject() []models.AlertNotice
	CreateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error)
	UpdateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error)
	DeleteNoticeObject(uuid string) error
	GetNoticeObject(uuid string) models.AlertNotice
	CheckNoticeObjectStatus(uuid string) string
}

func NewInterAlertNoticeService() InterAlertNoticeService {
	return &AlertNoticeService{}
}

func (ans *AlertNoticeService) SearchNoticeObject() []models.AlertNotice {

	var alertNoticeObject []models.AlertNotice
	globals.DBCli.Find(&alertNoticeObject)
	return alertNoticeObject

}

func (ans *AlertNoticeService) CreateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	alertNotice.Uuid = "n-" + cmd.RandId()
	err := tx.Create(alertNotice).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("创建通知对象失败", err)
		return models.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("事务提交失败", err)
		return models.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) UpdateNoticeObject(alertNotice models.AlertNotice) (models.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	err := tx.Model(&models.AlertNotice{}).Where("uuid = ?", alertNotice.Uuid).Updates(&alertNotice).Error
	if err != nil {
		tx.Rollback()
		return models.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return models.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) DeleteNoticeObject(uuid string) error {

	var ruleNum1, ruleNum2 int64
	globals.DBCli.Model(&models.AlertRule{}).Where("notice_id = ?", uuid).Count(&ruleNum1)
	globals.DBCli.Model(&models.AlertRule{}).Where("notice_group LIKE ?", "%"+uuid+"%").Count(&ruleNum2)
	if ruleNum1 != 0 || ruleNum2 != 0 {
		return fmt.Errorf("无法删除通知对象 %s, 因为已有告警规则绑定", uuid)
	}

	tx := globals.DBCli.Begin()
	err := tx.Where("uuid = ?", uuid).Delete(&models.AlertNotice{}).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("删除通知对象失败", err)
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("事务提交失败", err)
		return err
	}
	return nil

}

func (ans *AlertNoticeService) GetNoticeObject(uuid string) models.AlertNotice {

	var alertNoticeObject models.AlertNotice
	globals.DBCli.Where("uuid = ?", uuid).Find(&alertNoticeObject)
	return alertNoticeObject

}

func (ans *AlertNoticeService) CheckNoticeObjectStatus(uuid string) string {

	// ToDo

	return ""
}

const PrometheusAlertTest = `{
    "alerts":[
        {
            "annotations":{
                "description":"test",
                "summary":"test"
            },
            "endsAt":"0001-01-01T08:05:43.000Z",
            "fingerprint":"8888888888",
            "generatorURL":"http://0425df9dd50d:9090/graph?g0.expr=up+%3D%3D+0\u0026g0.tab=1",
            "labels":{
                "alertname":"test",
                "instance":"test",
                "job":"prometheus",
                "severity":"serious"
            },
            "startsAt":"0001-01-01T08:05:43.000Z",
            "status":"firing"
        }
    ],
    "commonAnnotations":{
        "description":"test",
        "summary":"test"
    },
    "commonLabels":{
        "alertname":"test",
        "instance":"test",
        "job":"prometheus",
        "severity":"serious"
    },
    "externalURL":"http://test:9093",
    "groupLabels":{
        "alertname":"test"
    },
    "receiver":"web\\.hook",
    "status":"firing",
    "truncatedAlerts":0,
    "version":"4"
}`

const AliSlsAlertTest = `["{\"name\": \"test\",\"fingerprint\": \"88888888\",\"region\": \"cn-beijing\",\"status\": \"firing\",\"alert_time\": \"test\",\"fire_time\": \"test\",\"resolve_time\": \"test\",\"host\": \"test\",\"statusCode\": \"UNSET\",\"traceID\": \"test\",\"logs\": \"\"[]\"\",\"attribute\": \"\"test\"\"}"]`
