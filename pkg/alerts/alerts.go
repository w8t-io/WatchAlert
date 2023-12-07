package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"prometheus-manager/pkg"
	"prometheus-manager/pkg/cache"
	"prometheus-manager/utils"
	"strings"
)

func CreateAlertSilence(challengeInfo map[string]interface{}) error {

	var (
		cardInfo models.CardInfo
	)

	rawDataJson, _ := json.Marshal(challengeInfo)
	_ = json.Unmarshal(rawDataJson, &cardInfo)

	// To Json
	kLabel := cardInfo.Action.Value.(map[string]interface{})
	silencesValueJson, _ := json.Marshal(kLabel)
	bodyReader := bytes.NewReader(silencesValueJson)

	fingerprintID := kLabel["comment"]
	labelData := GetAlertSilencesFingerprintID(fingerprintID.(string))
	if labelData == fingerprintID {
		globals.Logger.Sugar().Info("报警静默已存在, 无需重新创建, 报警ID ->", labelData)
		return fmt.Errorf("报警静默已存在, 无需重新创建")
	}

	_, err := utils.Post(globals.Config.AlertManager.URL+"/api/v2/silences", bodyReader)
	if err != nil {
		globals.Logger.Sugar().Error("创建报警静默失败 ->", string(silencesValueJson))
		return fmt.Errorf("创建报警静默失败")
	}
	globals.Logger.Sugar().Info("创建报警静默成功 ->", string(silencesValueJson))

	var (
		promAlertManager = make(map[string]interface{})
		cacheAlertInfo   interface{}
		AlertInfo        []interface{}
	)

	cacheValue := globals.CacheCli.Get(fingerprintID.(string))
	cacheAlertJson, _ := json.Marshal(cacheValue.(cache.CacheItem).Values)
	_ = json.Unmarshal(cacheAlertJson, &cacheAlertInfo)
	AlertInfo = append(AlertInfo, cacheAlertInfo)

	// 将告警状态转换为静默状态
	promAlertManager["alerts"] = AlertInfo
	promAlertManager["alerts"].([]interface{})[0].(map[string]interface{})["status"] = "silence"
	prometheusAlertBody, _ := json.Marshal(promAlertManager)

	// 发送消息卡片
	actionUserID := challengeInfo["user_id"].(string)
	err = pkg.SendMessageToWebhook(actionUserID, prometheusAlertBody, "")
	if err != nil {
		log.Println("消息卡片发送失败", err)
		return err
	}

	return nil

}

func GetAlertSilencesFingerprintID(fingerprintID string) (labelData interface{}) {

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
