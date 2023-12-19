package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"prometheus-manager/globals"
	"prometheus-manager/utils/http"
)

type EventController struct{}

func (ec *EventController) AlertEventMsg(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	resp, _ := ioutil.ReadAll(ctx.Request.Body)

	err := eventService.PushAlertToWebhook("", resp, uuid)
	if err != nil {
		ctx.JSON(500, gin.H{"code": 500, "data": err})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "data": "消息发送成功"})

}

func (ec *EventController) FeiShuEvent(ctx *gin.Context) {

	var challengeInfo map[string]interface{}

	uuid := ctx.Query("uuid")
	rawDataIO := ctx.Request.Body
	rawData, _ := ioutil.ReadAll(rawDataIO)

	err := json.Unmarshal(rawData, &challengeInfo)
	if err != nil {
		globals.Logger.Sugar().Error("飞书回调参数解析失败 ->", err)
		return
	}

	ctx.JSON(200, gin.H{"challenge": challengeInfo["challenge"]})

	/*
		回传数据示例:
		map[action:map[tag:button value:map[comment:f7f8ad553ebd24d2 createdBy:1 endsAt:2023-09-28T02:48:13.308Z id: matchers:[map[isEqual:true isRegex:false name:severity value:serious] map[isEqual:true isRegex:false name:alertname value:Exporter Componen is Down] map[isEqual:true isRegex:false name:instance value:localhost:9090] map[isEqual:true isRegex:false name:job value:prometheus]] startsAt:2023-09-28T02:27:43.113Z]] open_chat_id:oc_17518f8653e322021f71e4f0ad3ac08b open_id:ou_67a11fb06dceca64fec95dbc8c7133b2 open_message_id:om_e918cce173baf23f04ba48efbe5dc252 tenant_key:17217576210a575f token:c-0da3146133d9d1a9ea7a31b940527cbc7243283c user_id:886341fg]
	*/

	jsonData, _ := json.Marshal(challengeInfo)
	body := bytes.NewReader(jsonData)
	_, err = http.Post("http://127.0.0.1:"+globals.Config.Server.Port+"/api/v1/alert/createSilence?uuid="+uuid, body)
	if err != nil {
		log.Println(err)
		return
	}

	//amc.CreateAlertSilences(challengeInfo["user_id"].(string), challengeInfo)

}
