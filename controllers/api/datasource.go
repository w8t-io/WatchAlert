package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
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
	}

}

func (dc DatasourceController) Create(ctx *gin.Context) {

	var dataSource models.AlertDataSource
	_ = ctx.ShouldBindJSON(&dataSource)
	tid, _ := ctx.Get("TenantID")
	dataSource.TenantId = tid.(string)
	err := dataSourceService.Create(dataSource)

	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (dc DatasourceController) List(ctx *gin.Context) {

	data, _ := dataSourceService.List(ctx)

	response.Success(ctx, data, "success")
}

func (dc DatasourceController) Get(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	id := ctx.Query("id")
	dsType := ctx.Query("dsType")

	data := dataSourceService.Get(tid.(string), id, dsType)
	response.Success(ctx, data, "success")

}

func (dc DatasourceController) Search(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dataSourceService.Search(r)
	})
}

func (dc DatasourceController) Update(ctx *gin.Context) {

	var datasource models.AlertDataSource
	_ = ctx.ShouldBindJSON(&datasource)
	tid, _ := ctx.Get("TenantID")
	datasource.TenantId = tid.(string)
	err := dataSourceService.Update(datasource)

	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (dc DatasourceController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	tid, _ := ctx.Get("TenantID")
	err := dataSourceService.Delete(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}
