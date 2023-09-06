package event

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"prometheus-manager/globals"
	"prometheus-manager/services/alerts"
)

var (
	amc *alerts.AlertManagerCollector
)

func (aemc *AlertEventMsgCollector) FeiShuEvent(ctx *gin.Context) {

	var challengeInfo map[string]interface{}

	rawDataIO := ctx.Request.Body
	rawData, _ := ioutil.ReadAll(rawDataIO)

	err := json.Unmarshal(rawData, &challengeInfo)
	if err != nil {
		globals.Logger.Sugar().Error("飞书回调参数解析失败 ->", err)
		return
	}

	ctx.JSON(200, gin.H{"challenge": challengeInfo["challenge"]})

	amc.CreateAlertSilences(challengeInfo)

}
