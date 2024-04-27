package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(gin *gin.Engine) {

	gin.GET("hello", health)

}

func health(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"success": "true",
	})

}
