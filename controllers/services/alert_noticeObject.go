package services

import (
	"bytes"
	"prometheus-manager/controllers/dao"
	"prometheus-manager/globals"
	"prometheus-manager/utils/cmd"
	"prometheus-manager/utils/feishu"
	"prometheus-manager/utils/http"
)

type AlertNoticeService struct{}

type InterAlertNoticeService interface {
	SearchNoticeObject() []dao.AlertNotice
	CreateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error)
	UpdateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error)
	DeleteNoticeObject(uuid string) error
	GetNoticeObject(uuid string) dao.AlertNotice
	CheckNoticeObjectStatus(uuid string) string
}

func NewInterAlertNoticeService() InterAlertNoticeService {
	return &AlertNoticeService{}
}

func (ans *AlertNoticeService) SearchNoticeObject() []dao.AlertNotice {

	var alertNoticeObject []dao.AlertNotice
	globals.DBCli.Find(&alertNoticeObject)
	return alertNoticeObject

}

func (ans *AlertNoticeService) CreateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	alertNotice.Uuid = cmd.RandUuid()
	err := tx.Create(alertNotice).Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("创建通知对象失败", err)
		return dao.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		globals.Logger.Sugar().Error("事务提交失败", err)
		return dao.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) UpdateNoticeObject(alertNotice dao.AlertNotice) (dao.AlertNotice, error) {

	tx := globals.DBCli.Begin()
	err := tx.Model(&dao.AlertNotice{}).Where("uuid = ?", alertNotice.Uuid).Updates(&alertNotice).Error
	if err != nil {
		tx.Rollback()
		return dao.AlertNotice{}, err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return dao.AlertNotice{}, err
	}
	return alertNotice, nil

}

func (ans *AlertNoticeService) DeleteNoticeObject(uuid string) error {

	tx := globals.DBCli.Begin()
	err := tx.Where("uuid = ?", uuid).Delete(&dao.AlertNotice{}).Error
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

func (ans *AlertNoticeService) GetNoticeObject(uuid string) dao.AlertNotice {

	var alertNoticeObject dao.AlertNotice
	globals.DBCli.Where("uuid = ?", uuid).Find(&alertNoticeObject)
	return alertNoticeObject

}

func (ans *AlertNoticeService) CheckNoticeObjectStatus(uuid string) string {

	var alertNoticeData dao.AlertNotice

	globals.DBCli.Model(&dao.AlertNotice{}).Where("uuid = ?", uuid).Find(&alertNoticeData)

	noticeStatus := "正常"
	testBodyData := map[string]string{
		"Prometheus": PrometheusAlertTest,
		"AliSls":     AliSlsAlertTest,
	}

	switch alertNoticeData.NoticeType {
	case "FeiShu":
		if !feishu.CheckFeiShuChatId(alertNoticeData.FeishuChatId) {
			noticeStatus = "异常"
		} else {
			post, err := http.Post("http://localhost:9001/api/v1/prom/prometheusAlert?uuid="+alertNoticeData.Uuid, bytes.NewReader([]byte(testBodyData[alertNoticeData.DataSource])))
			if err != nil && post.StatusCode != 200 {
				noticeStatus = "异常"
			}
		}
	}

	globals.DBCli.Model(&dao.AlertNotice{}).Where("uuid = ?", uuid).Update("notice_status", noticeStatus)
	return noticeStatus
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
