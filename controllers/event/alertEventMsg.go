package event

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"prometheus-manager/globals"
	"prometheus-manager/pkg"
)

type AlertEventMsgCollector struct{}

var (
	promAlertManager = make(map[string]interface{})
)

func (aemc *AlertEventMsgCollector) AlertEventMsg(ctx *gin.Context) {

	globals.AlertType = ctx.Query("type")
	globals.DataSource = ctx.Query("dataSource")

	resp, _ := ioutil.ReadAll(ctx.Request.Body)
	globals.RespBody = resp
	fmt.Println("======= body ->", string(resp))

	err := json.Unmarshal(resp, &promAlertManager)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return
	}

	err = pkg.SendMessageToWebhook("", promAlertManager)
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "data": err})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "data": "消息发送成功"})

}
