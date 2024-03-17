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
	err := dataSourceService.Create(dataSource)

	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (adsc *AlertDataSourceController) List(ctx *gin.Context) {

	data, _ := dataSourceService.List()

	response.Success(ctx, data, "success")
}

func (adsc *AlertDataSourceController) Search(ctx *gin.Context) {

	id := ctx.Query("id")
	dsType := ctx.Query("dsType")

	data := dataSourceService.Get(id, dsType)
	response.Success(ctx, data, "success")

}

func (adsc *AlertDataSourceController) Update(ctx *gin.Context) {

	var datasource models.AlertDataSource
	_ = ctx.ShouldBindJSON(&datasource)

	err := dataSourceService.Update(datasource)

	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (adsc *AlertDataSourceController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := dataSourceService.Delete(id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}
