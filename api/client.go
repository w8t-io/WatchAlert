package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type ClientController struct{}

func (cc ClientController) API(gin *gin.RouterGroup) {
	c := gin.Group("c")
	c.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		c.GET("getJaegerService", cc.GetJaegerService)
	}
}

func (cc ClientController) GetJaegerService(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.ClientService.GetJaegerService(r)
	})
}
