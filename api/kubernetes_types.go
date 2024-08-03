package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

type KubernetesTypesController struct{}

func (ktc KubernetesTypesController) API(gin *gin.RouterGroup) {
	k8s := gin.Group("kubernetes")
	k8s.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		k8s.GET("getResourceList", ktc.getResourceList)
		k8s.GET("getReasonList", ktc.getReasonList)
	}
}

func (ktc KubernetesTypesController) getResourceList(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return types.EventResourceTypeList, nil
	})
}

func (ktc KubernetesTypesController) getReasonList(ctx *gin.Context) {
	r := new(models.RequestEventTypes)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return types.EventReasonLMapping[r.Resource], nil
	})
}
