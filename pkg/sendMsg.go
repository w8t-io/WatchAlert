package pkg

import (
	"encoding/json"
	"prometheus-manager/globals"
	"prometheus-manager/pkg/renderTemplates/aliyun"
	prom "prometheus-manager/pkg/renderTemplates/prometheus"
)

func SendMessageToWebhook(actionUserID string, body []byte, env string) error {

	switch globals.DataSource {
	case "prometheus":

		err := prometheus(actionUserID, body)
		if err != nil {
			return err
		}

	case "ali":

		err := aLiYun(body, env)
		if err != nil {
			return err
		}

	}

	return nil

}

func prometheus(actionUserID string, body []byte) error {

	var (
		prometheusAlertInfo = make(map[string]interface{})
	)

	err := json.Unmarshal(body, &prometheusAlertInfo)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return err
	}

	switch globals.AlertType {
	case "feishu":
		err = prom.PrometheusTemplate(actionUserID, prometheusAlertInfo).FeiShu()
		if err != nil {
			return err
		}
	}
	return nil

}

func aLiYun(body []byte, env string) error {

	switch globals.AlertType {
	case "feishu":
		err := aliyun.ALiYunSlSTemplate(body, env).FeiShu()
		if err != nil {
			return err
		}
	}

	return nil

}
