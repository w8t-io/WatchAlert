package v1

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/controllers/alerts"
	"prometheus-manager/controllers/event"
)

var (
	aemc event.AlertEventMsgCollector
	amc  alerts.AlertManagerCollector
)

func AlertEventMsg(gin *gin.Engine) {

	api := gin.Group("api/v1/prom/")
	{
		// 接收 Alert
		api.POST("prometheusAlert", aemc.AlertEventMsg)
		// 接收飞书回调
		api.POST("feiShuEvent", amc.FeiShuEvent)
	}

}
