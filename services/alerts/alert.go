package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"prometheus-manager/utils"
	"prometheus-manager/utils/sendAlertMessage"
	"strings"
)

type AlertManagerCollector struct{}

func (amc *AlertManagerCollector) ListAlerts() ([]models.GettableAlert, error) {

	req, err := http.NewRequest(http.MethodGet, "http://192.168.1.193:30111/api/v2/alerts", nil)
	if err != nil {
		log.Println("1 get failed", err)
		return []models.GettableAlert{}, err
	}

	body, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("2 get failed", err)
		return []models.GettableAlert{}, err
	}

	content, err := ioutil.ReadAll(body.Body)

	var gettableAlert []models.GettableAlert

	err = json.Unmarshal(content, &gettableAlert)
	if err != nil {
		log.Println("解析失败", err)
		return []models.GettableAlert{}, err
	}

	return gettableAlert, nil

	//for k, v := range gettableAlert {
	//
	//	fmt.Println("---")
	//	var labelsMap = make(map[string]string)
	//	for labelKey, labelValue := range v.Labels {
	//		if labelKey == "alertname" {
	//			continue
	//		}
	//		labelsMap[labelKey] = labelValue
	//	}
	//
	//	fmt.Printf("序列: %v\n名称: %s\n标签: %s\n描述: %s\n详情: %v\n状态: %v\n开始时间: %v\n结束时间: %v\n指纹: %v\nxx: %v\nxx: %v\nxx: %s\nxx: %s\n", k, v.Labels["alertname"], labelsMap, v.Annotations["description"], v.Annotations["summary"], v.Status, v.StartsAt, v.EndsAt, v.Fingerprint, v.GeneratorURL, v.Receivers, v.UpdatedAt, v.Alert)
	//}

}

func (amc *AlertManagerCollector) CreateAlertSilences(actionUser string, challengeInfo interface{}) {

	var action bool
	for _, user := range globals.Config.FeiShu.ActionUser {
		if actionUser == user {
			action = true
			break
		}
	}

	if !action {
		globals.Logger.Sugar().Error("「" + actionUser + "」你无权操作创建静默规则")
		return
	}

	rawDataJson, _ := json.Marshal(challengeInfo)

	var cardInfo models.CardInfo
	_ = json.Unmarshal(rawDataJson, &cardInfo)

	// To Json
	kLabel := cardInfo.Action.Value.(map[string]interface{})
	silencesValueJson, _ := json.Marshal(kLabel)
	bodyReader := bytes.NewReader(silencesValueJson)

	fingerprintID := kLabel["comment"]
	labelData := amc.GetAlertSilencesFingerprintID(fingerprintID.(string))
	if labelData == fingerprintID {
		globals.Logger.Sugar().Info("报警静默已存在, 无需重新创建, 报警ID ->", labelData)
		return
	}

	_, err := utils.Post(globals.Config.AlertManager.URL+"/api/v2/silences", bodyReader)
	if err != nil {
		globals.Logger.Sugar().Error("创建报警静默失败 ->", string(silencesValueJson))
		return
	}
	globals.Logger.Sugar().Info("创建报警静默成功 ->", string(silencesValueJson))

	var (
		promAlertManager = make(map[string]interface{})
	)
	err = json.Unmarshal(sendAlertMessage.RespBody, &promAlertManager)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return
	}
	promAlertManager["alerts"].([]interface{})[0].(map[string]interface{})["status"] = "silence"

	err = sendAlertMessage.SendMsg(actionUser, sendAlertMessage.DataSource, sendAlertMessage.AlertType, promAlertManager)
	if err != nil {
		globals.Logger.Sugar().Error("静默消息卡片发送失败 ->", err)
		return
	}

}

func (amc *AlertManagerCollector) GetAlertSilencesFingerprintID(fingerprintID string) (labelData interface{}) {

	resp, err := utils.Get(globals.Config.AlertManager.URL + "/api/v2/silences")
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)

	var (
		res             []models.SearchAlertManager
		activeLabelList []interface{}
	)
	err = json.Unmarshal(content, &res)

	for _, v := range res {
		if v.Status.State == "active" {
			labelJson, _ := json.Marshal(v.Comment)
			activeLabelList = append(activeLabelList, string(labelJson))
		}
	}

	for _, v := range activeLabelList {
		activeV := strings.Replace(v.(string), "\"", "", -1)
		if activeV == fingerprintID {
			return fingerprintID
		}
	}

	return nil
}
