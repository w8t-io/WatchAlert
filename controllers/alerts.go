package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"prometheus-manager/models"
	"prometheus-manager/pkg/alerts"
)

type AlertController struct{}

func (ac *AlertController) CreateSilence(ctx *gin.Context) {

	var challengeInfo map[string]interface{}

	body := ctx.Request.Body
	bodyByte, _ := ioutil.ReadAll(body)
	_ = json.Unmarshal(bodyByte, &challengeInfo)

	err := alerts.CreateAlertSilence(challengeInfo)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 2001,
			"data": err.Error(),
			"msg":  "创建失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 2000,
		"data": nil,
		"msg":  "创建成功",
	})

}

func (amc *AlertController) ListAlerts() ([]models.GettableAlert, error) {

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
