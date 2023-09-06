package event

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"prometheus-manager/globals"
	"prometheus-manager/utils/sendAlertMessage"
)

type AlertEventMsgCollector struct{}

func (aemc *AlertEventMsgCollector) AlertEventMsg(ctx *gin.Context) {

	alertType := ctx.Query("type")

	pAlertManagerJson := make(map[string]interface{})
	resp, _ := ioutil.ReadAll(ctx.Request.Body)

	err := json.Unmarshal(resp, &pAlertManagerJson)
	if err != nil {
		globals.Logger.Sugar().Error("告警信息解析失败 ->", err)
		return
	}

	sendAlertMessage.SendMsg(alertType, pAlertManagerJson)

}
