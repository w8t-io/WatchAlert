package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/pkg/community/aws/cloudwatch/types"
)

type AWSCloudWatchController struct{}

func (cwc AWSCloudWatchController) API(gin *gin.RouterGroup) {
	community := gin.Group("community")
	community.Use(
		middleware.Cors(),
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		cloudwatch := community.Group("cloudwatch")
		{
			cloudwatch.GET("metricTypes", cwc.GetMetricTypes)
			cloudwatch.GET("metricNames", cwc.GetMetricNames)
			cloudwatch.GET("statistics", cwc.GetStatistics)
			cloudwatch.GET("dimensions", cwc.GetDimensions)
		}
	}
}

func (cwc AWSCloudWatchController) GetMetricTypes(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetMetricTypes()
	})
}

func (cwc AWSCloudWatchController) GetMetricNames(ctx *gin.Context) {
	q := new(types.MetricNamesQuery)
	BindQuery(ctx, q)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetMetricNames(q)
	})
}

func (cwc AWSCloudWatchController) GetStatistics(ctx *gin.Context) {
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetStatistics()
	})
}

func (cwc AWSCloudWatchController) GetDimensions(ctx *gin.Context) {
	q := new(types.RdsDimensionReq)
	BindQuery(ctx, q)
	Service(ctx, func() (interface{}, interface{}) {
		return services.AWSCloudWatchService.GetDimensions(q)
	})
}
