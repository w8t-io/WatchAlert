package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"watchAlert/internal/global"
	"watchAlert/pkg/tools"
)

type CallbackController struct{}

// FeiShuEvent 飞书回调
func (cc *CallbackController) FeiShuEvent(ctx *gin.Context) {
	var challengeInfo map[string]interface{}
	uuid := ctx.Query("uuid")

	if err := tools.ParseReaderBody(ctx.Request.Body, &challengeInfo); err != nil {
		return
	}

	ctx.JSON(200, gin.H{"challenge": challengeInfo["challenge"]})

	jsonData, _ := json.Marshal(challengeInfo)
	body := bytes.NewReader(jsonData)
	_, err := tools.Post(nil, "http://127.0.0.1:"+global.Config.Server.Port+"/api/v1/alert/createSilence?uuid="+uuid, body)
	if err != nil {
		log.Println(err)
		return
	}

}
