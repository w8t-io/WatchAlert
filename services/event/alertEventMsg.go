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

	globals.AlertType = ctx.Query("type")

	resp, _ := ioutil.ReadAll(ctx.Request.Body)
	globals.RespBody = resp

	err := json.Unmarshal(resp, &promAlertManager)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return
	}

	sendAlertMessage.SendMsg(globals.AlertType, promAlertManager)

}
