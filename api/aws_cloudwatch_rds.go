package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/pkg/community/aws/cloudwatch/types"
)

type AWSCloudWatchRDSController struct{}

func (a AWSCloudWatchRDSController) API(gin *gin.RouterGroup) {
	community := gin.Group("community")
	community.Use(
		middleware.Cors(),
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		rds := community.Group("rds")
		{
			rds.GET("instances", a.GetRdsInstanceIdentifier)
			rds.GET("clusters", a.GetRdsClusterIdentifier)
		}
	}
}

func (a AWSCloudWatchRDSController) GetRdsInstanceIdentifier(ctx *gin.Context) {
	req := new(types.RdsInstanceReq)
	BindQuery(ctx, req)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchRdsService.GetDBInstanceIdentifier(req)
	})
}

func (a AWSCloudWatchRDSController) GetRdsClusterIdentifier(ctx *gin.Context) {
	req := new(types.RdsClusterReq)
	BindQuery(ctx, req)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchRdsService.GetDBClusterIdentifier(req)
	})
}
