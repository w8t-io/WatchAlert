package event

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"prometheus-manager/globals"
	"prometheus-manager/utils/sendAlertMessage"
)

type AlertEventMsgCollector struct{}

var (
	promAlertManager = make(map[string]interface{})
)

func (aemc *AlertEventMsgCollector) AlertEventMsg(ctx *gin.Context) {

	sendAlertMessage.AlertType = ctx.Query("type")
	sendAlertMessage.DataSource = ctx.Query("dataSource")

	resp, _ := ioutil.ReadAll(ctx.Request.Body)
	sendAlertMessage.RespBody = resp

	err := json.Unmarshal(resp, &promAlertManager)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return
	}

	err = sendAlertMessage.SendMsg("", sendAlertMessage.DataSource, sendAlertMessage.AlertType, promAlertManager)
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "data": err})
	}
	ctx.JSON(200, gin.H{"code": 200, "data": "消息发送成功"})

}
