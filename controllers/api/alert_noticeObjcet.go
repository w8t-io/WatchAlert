package api

import (
	"github.com/gin-gonic/gin"
	"prometheus-manager/controllers/dao"
	"prometheus-manager/utils/feishu"
)

type AlertNoticeObjectController struct{}

func (ano *AlertNoticeObjectController) List(ctx *gin.Context) {

	object := alertNoticeService.SearchNoticeObject()
	ctx.JSON(200, gin.H{
		"code": "4000",
		"data": object,
		"msg":  "查询成功",
	})

}

func (ano *AlertNoticeObjectController) Create(ctx *gin.Context) {

	var alertNotice dao.AlertNotice
	_ = ctx.ShouldBindJSON(&alertNotice)

	object, err := alertNoticeService.CreateNoticeObject(alertNotice)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": "4001",
			"data": err,
			"msg":  "创建失败",
		})
	}
	ctx.JSON(200, gin.H{
		"code": "4000",
		"data": object,
		"msg":  "创建成功",
	})

}

func (ano *AlertNoticeObjectController) Update(ctx *gin.Context) {

	var alertNotice dao.AlertNotice
	_ = ctx.ShouldBindJSON(&alertNotice)

	object, err := alertNoticeService.UpdateNoticeObject(alertNotice)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": "4002",
			"data": err,
			"msg":  "创建失败",
		})
	}
	ctx.JSON(200, gin.H{
		"code": "4000",
		"data": object,
		"msg":  "创建成功",
	})
}

func (ano *AlertNoticeObjectController) Delete(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	err := alertNoticeService.DeleteNoticeObject(uuid)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": "4003",
			"data": err,
			"msg":  "删除失败",
		})
	}
	ctx.JSON(200, gin.H{
		"code": "4000",
		"data": "",
		"msg":  "删除成功",
	})
}

func (ano *AlertNoticeObjectController) Get(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	object := alertNoticeService.GetNoticeObject(uuid)
	ctx.JSON(200, gin.H{
		"code": "4000",
		"data": object,
		"msg":  "查询成功",
	})

}

func (ano *AlertNoticeObjectController) GetFeishuChats(ctx *gin.Context) {

	object := feishu.GetFeiShuChatsID()
	ctx.JSON(200, gin.H{
		"code": "4000",
		"data": object,
		"msg":  "查询成功",
	})

}
