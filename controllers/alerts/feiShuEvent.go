package alerts

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"prometheus-manager/globals"
)

func (amc *AlertManagerCollector) FeiShuEvent(ctx *gin.Context) {

	var challengeInfo map[string]interface{}

	rawDataIO := ctx.Request.Body
	rawData, _ := ioutil.ReadAll(rawDataIO)

	err := json.Unmarshal(rawData, &challengeInfo)
	if err != nil {
		globals.Logger.Sugar().Error("飞书回调参数解析失败 ->", err)
		return
	}

	ctx.JSON(200, gin.H{"challenge": challengeInfo["challenge"]})

	//info := f.GetFeiShuUserInfo(challengeInfo["user_id"].(string))

	/*
		回传数据示例:
		map[action:map[tag:button value:map[comment:f7f8ad553ebd24d2 createdBy:1 endsAt:2023-09-28T02:48:13.308Z id: matchers:[map[isEqual:true isRegex:false name:severity value:serious] map[isEqual:true isRegex:false name:alertname value:Exporter Componen is Down] map[isEqual:true isRegex:false name:instance value:localhost:9090] map[isEqual:true isRegex:false name:job value:prometheus]] startsAt:2023-09-28T02:27:43.113Z]] open_chat_id:oc_17518f8653e322021f71e4f0ad3ac08b open_id:ou_67a11fb06dceca64fec95dbc8c7133b2 open_message_id:om_e918cce173baf23f04ba48efbe5dc252 tenant_key:17217576210a575f token:c-0da3146133d9d1a9ea7a31b940527cbc7243283c user_id:886341fg]
	*/

	amc.CreateAlertSilences(challengeInfo["user_id"].(string), challengeInfo)

}
