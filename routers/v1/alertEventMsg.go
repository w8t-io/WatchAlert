package v1

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/services/event"
)

var (
	aemc event.AlertEventMsgCollector
)

func AlertEventMsg(gin *gin.Engine) {

	api := gin.Group("api/v1/prom/")
	{
		api.POST("prometheusAlert", aemc.AlertEventMsg)
		api.POST("feiShuEvent", aemc.FeiShuEvent)
	}

}
