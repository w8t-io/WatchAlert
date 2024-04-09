package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"watchAlert/public/globals"
	"watchAlert/public/utils/http"
)

type EventController struct{}

// FeiShuEvent 飞书回调
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

	jsonData, _ := json.Marshal(challengeInfo)
	body := bytes.NewReader(jsonData)
	_, err = http.Post("http://127.0.0.1:"+globals.Config.Server.Port+"/api/v1/alert/createSilence?uuid="+uuid, body)
	if err != nil {
		log.Println(err)
		return
	}

}
