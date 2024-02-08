package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/models"
)

type AlertDataSourceController struct {
}

func (adsc *AlertDataSourceController) Create(ctx *gin.Context) {

	var dataSource models.AlertDataSource
	_ = ctx.ShouldBindJSON(&dataSource)
	err := dataSourceService.Create(dataSource)

	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "",
		"msg":  "success",
	})

}

func (adsc *AlertDataSourceController) List(ctx *gin.Context) {

	data, _ := dataSourceService.List()

	ctx.JSON(200, gin.H{
		"code": "200",
		"data": data,
		"msg":  "success",
	})

}

func (adsc *AlertDataSourceController) Search(ctx *gin.Context) {

	id := ctx.Query("id")
	dsType := ctx.Query("dsType")

	data := dataSourceService.Get(id, dsType)
	ctx.JSON(200, gin.H{
		"code": "200",
		"data": data,
		"msg":  "success",
	})

}

func (adsc *AlertDataSourceController) Update(ctx *gin.Context) {

	var datasource models.AlertDataSource
	_ = ctx.ShouldBindJSON(&datasource)

	err := dataSourceService.Update(datasource)

	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
	}

	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "",
		"msg":  "success",
	})

}

func (adsc *AlertDataSourceController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := dataSourceService.Delete(id)
	if err != nil {
		ctx.JSON(401, gin.H{
			"code": "401",
			"data": err.Error(),
			"msg":  "failed",
		})
	}

	ctx.JSON(200, gin.H{
		"code": "200",
		"data": "",
		"msg":  "success",
	})
}
