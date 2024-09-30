package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type SubscribeController struct{}

func (sc SubscribeController) API(gin *gin.RouterGroup) {
	subscribeA := gin.Group("subscribe")
	subscribeA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.AuditingLog(),
		middleware.ParseTenant(),
	)
	{
		subscribeA.POST("createSubscribe", sc.Create)
		subscribeA.POST("deleteSubscribe", sc.Delete)
	}

	subscribeB := gin.Group("subscribe")
	subscribeB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		subscribeB.GET("listSubscribe", sc.List)
		subscribeB.GET("getSubscribe", sc.Get)
	}
}

func (sc SubscribeController) List(ctx *gin.Context) {
	r := new(models.AlertSubscribeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.List(r)
	})
}

func (sc SubscribeController) Get(ctx *gin.Context) {
	r := new(models.AlertSubscribeQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.Get(r)
	})
}

func (sc SubscribeController) Create(ctx *gin.Context) {
	r := new(models.AlertSubscribe)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)
	uid, _ := ctx.Get("UserId")
	r.SUserId = uid.(string)
	ue, _ := ctx.Get("UserEmail")
	r.SUserEmail = ue.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.Create(r)
	})
}

func (sc SubscribeController) Delete(ctx *gin.Context) {
	r := new(models.AlertSubscribeQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.STenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.SubscribeService.Delete(r)
	})
}
