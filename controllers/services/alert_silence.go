package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"watchAlert/controllers/dto"
	"watchAlert/globals"
	"watchAlert/utils/cache"
	"watchAlert/utils/http"
)

type AlertSilenceService struct{}

type InterAlertSilenceService interface {
	CreateAlertSilence(uuid string, cardInfo dto.CardInfo) error
}

func NewInterAlertSilenceService() InterAlertSilenceService {
	return &AlertSilenceService{}
}

func (ass *AlertSilenceService) CreateAlertSilence(uuid string, cardInfo dto.CardInfo) error {

	var (
		cacheAlertInfo dto.AlertInfo
		AlertInfoList  []dto.AlertInfo
	)

	// 获取静默请求参数
	kLabel := cardInfo.Action.Value
	silencesValueJson, _ := json.Marshal(kLabel)
	bodyReader := bytes.NewReader(silencesValueJson)

	// 检查是否已创建
	getAlertSilencesFingerprintID(kLabel.Comment)

	// 发起静默请求
	_, err := http.Post(globals.Config.AlertManager.URL+"/api/v2/silences", bodyReader)
	if err != nil {
		globals.Logger.Sugar().Error("创建报警静默失败 ->", string(silencesValueJson))
		return fmt.Errorf("创建报警静默失败")
	}
	globals.Logger.Sugar().Info("创建报警静默成功 ->", string(silencesValueJson))

	// 将告警状态转换为静默状态
	cacheValue := globals.CacheCli.Get(kLabel.Comment)
	cacheAlertJson, _ := json.Marshal(cacheValue.(cache.CacheItem).Values)
	_ = json.Unmarshal(cacheAlertJson, &cacheAlertInfo)
	cacheAlertInfo.Status = "silence"
	AlertInfoList = append(AlertInfoList, cacheAlertInfo)
	alerts := dto.Alerts{
		AlertList: AlertInfoList,
	}
	body, _ := json.Marshal(alerts)

	// 发送消息卡片
	err = NewInterEventService().PushAlertToWebhook(cardInfo.UserID, body, uuid)
	if err != nil {
		log.Println("消息卡片发送失败", err)
		return err
	}

	return nil

}

func getAlertSilencesFingerprintID(fingerprintID string) {

	resp, err := http.Get(globals.Config.AlertManager.URL + "/api/v2/silences")
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)

	var (
		res             []dto.SearchAlertManager
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
			globals.Logger.Sugar().Info("报警静默已存在, 无需重新创建, 报警ID ->", fingerprintID)
			return
		}
	}

}
