package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type AlertDataSourceController struct {
}

func (adsc *AlertDataSourceController) Create(ctx *gin.Context) {

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

func (adsc *AlertDataSourceController) List(ctx *gin.Context) {

	data, _ := dataSourceService.List(ctx)

	response.Success(ctx, data, "success")
}

func (adsc *AlertDataSourceController) Get(ctx *gin.Context) {

	tid, _ := ctx.Get("TenantID")
	id := ctx.Query("id")
	dsType := ctx.Query("dsType")

	data := dataSourceService.Get(tid.(string), id, dsType)
	response.Success(ctx, data, "success")

}

func (adsc *AlertDataSourceController) Search(ctx *gin.Context) {
	r := new(models.DatasourceQuery)
	BindQuery(ctx, r)
	tid, _ := ctx.Get("TenantID")
	r.TenantId = tid.(string)
	Service(ctx, func() (interface{}, interface{}) {
		return dataSourceService.Search(r)
	})
}

func (adsc *AlertDataSourceController) Update(ctx *gin.Context) {

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

func (adsc *AlertDataSourceController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	tid, _ := ctx.Get("TenantID")
	err := dataSourceService.Delete(tid.(string), id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}
