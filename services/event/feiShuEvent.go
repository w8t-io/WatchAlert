package event

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"prometheus-manager/globals"
	"prometheus-manager/services/alerts"
	"prometheus-manager/utils/sendAlertMessage"
)

var (
	amc *alerts.AlertManagerCollector
	f   sendAlertMessage.FeiShu
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

	resp := f.GetFeiShuUserInfo(challengeInfo["user_id"].(string))

	fmt.Println("=== 用户信息 ->", resp)

	amc.CreateAlertSilences(*resp.Data.User.Name, challengeInfo)

}
