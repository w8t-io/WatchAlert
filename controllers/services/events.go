package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"watchAlert/controllers/dao"
	"watchAlert/controllers/dto"
	"watchAlert/globals"
	"watchAlert/utils/renderTemplates/aliyun"
	prom "watchAlert/utils/renderTemplates/prometheus"
)

type EventService struct{}

type InterEventService interface {
	PushAlertToWebhook(actionUserID string, body []byte, uuid string) error
}

func NewInterEventService() InterEventService {
	return &EventService{}
}

func (es *EventService) PushAlertToWebhook(actionUserID string, body []byte, uuid string) error {

	alertNoticeObject := alertNotice.Get(uuid)
	_, dutyUser := dutySchedule.GetDutyScheduleInfo(alertNoticeObject.DutyId, strconv.Itoa(time.Now().Year())+
		"-"+strconv.Itoa(int(time.Now().Month()))+
		"-"+strconv.Itoa(time.Now().Day()))

	switch alertNoticeObject.DataSource {
	case "Prometheus":

		body = alertAggregation(body)
		err := prometheus(alertNoticeObject, actionUserID, body, dutyUser)
		if err != nil {
			return err
		}

	case "AliSls":

		err := aLiYun(alertNoticeObject, body, dutyUser)
		if err != nil {
			return err
		}

	}

	return nil

}

func prometheus(alertNotice dao.AlertNotice, actionUserID string, body []byte, dutyUser string) error {

	var (
		prometheusAlertInfo = make(map[string]interface{})
	)

	err := json.Unmarshal(body, &prometheusAlertInfo)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return err
	}

	switch alertNotice.NoticeType {
	case "FeiShu":
		err = prom.PrometheusTemplate(actionUserID, prometheusAlertInfo, dutyUser).
			FeiShu(alertNotice.FeishuChatId)
		if err != nil {
			return err
		}
	}
	return nil

}

func aLiYun(alertNotice dao.AlertNotice, body []byte, dutyUser string) error {

	switch alertNotice.NoticeType {
	case "FeiShu":
		err := aliyun.ALiYunSlSTemplate(body, alertNotice.Env, dutyUser).
			FeiShu(alertNotice.FeishuChatId)
		if err != nil {
			return err
		}
	}

	return nil

}

func alertAggregation(body []byte) []byte {

	var (
		alerts              dto.Alerts
		alertName           string
		filteredAlertInfo   []dto.AlertInfo
		countIdenticalAlert = make(map[string]int)
	)
	_ = json.Unmarshal(body, &alerts)
	for k, v := range alerts.AlertList {
		alertName = v.Labels["alertname"]
		countIdenticalAlert[alertName] = k + 1
	}
	if countIdenticalAlert[alertName] > 1 {
		filteredAlertInfo = append(filteredAlertInfo, alerts.AlertList[0])
		alerts.AlertList = filteredAlertInfo
		alerts.Aggregated = fmt.Sprintf("\n聚合 %d 条告警, 详情查看告警链接 ⬇️", countIdenticalAlert[alertName])
	}

	jsonData, _ := json.Marshal(alerts)
	return jsonData

}
