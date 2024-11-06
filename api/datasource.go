package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/pkg/tools"
)

type DatasourceController struct{}

/*
数据源 API
/api/w8t/datasource
*/
func (dc DatasourceController) API(gin *gin.RouterGroup) {
	datasourceA := gin.Group("datasource")
	datasourceA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		datasourceA.POST("dataSourceCreate", dc.Create)
		datasourceA.POST("dataSourceUpdate", dc.Update)
		datasourceA.POST("dataSourceDelete", dc.Delete)
	}

	datasourceB := gin.Group("datasource")
	datasourceB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		datasourceB.GET("dataSourceList", dc.List)
		datasourceB.GET("dataSourceGet", dc.Get)
		datasourceB.GET("dataSourceSearch", dc.Search)
		datasourceB.GET("promQuery", dc.PromQuery)
	}

}

func (dc DatasourceController) Create(ctx *gin.Context) {
	d := new(models.AlertDataSource)
	BindJson(ctx, d)

	tid, _ := ctx.Get("TenantID")
	d.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Create(d)
	})
}

func (dc DatasourceController) List(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.List(r)
	})
}

func (dc DatasourceController) Get(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Get(r)
	})
}

func (dc DatasourceController) Search(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Search(r)
	})
}

func (dc DatasourceController) Update(ctx *gin.Context) {
	r := new(models.AlertDataSource)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Update(r)
	})
}

func (dc DatasourceController) Delete(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindJson(ctx, r)

	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.DatasourceService.Delete(r)
	})
}

func (dc DatasourceController) PromQuery(ctx *gin.Context) {
	r := new(models.PromQueryReq)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		var res models.PromQueryRes
		path := "/api/v1/query"
		if r.DatasourceType == "VictoriaMetrics" {
			path = "/prometheus" + path
		}

		encodedQuery := url.QueryEscape(r.Query)
		get, err := tools.Get(nil, fmt.Sprintf("%s%s?query=%s", r.Addr, path, encodedQuery))
		if err != nil {
			return nil, err
		}

		if err := tools.ParseReaderBody(get.Body, &res); err != nil {
			return nil, err
		}

		return res, nil
	})
}
