package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"time"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/pkg/provider"
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
		datasourceB.POST("dataSourcePing", dc.Ping)
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
		var res provider.QueryResponse
		path := "/api/v1/query"
		params := url.Values{}
		params.Add("query", r.Query)
		params.Add("time", strconv.FormatInt(time.Now().Unix(), 10))
		fullURL := fmt.Sprintf("%s%s?%s", r.Addr, path, params.Encode())
		get, err := tools.Get(nil, fullURL, 10)
		if err != nil {
			return nil, err
		}

		if err := tools.ParseReaderBody(get.Body, &res); err != nil {
			return nil, err
		}

		return res, nil
	})
}

func (dc DatasourceController) Ping(ctx *gin.Context) {
	r := new(models.AlertDataSource)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		ok := provider.CheckDatasourceHealth(*r)
		if !ok {
			return "", fmt.Errorf("数据源不可达!")
		}
		return "", nil
	})
}
