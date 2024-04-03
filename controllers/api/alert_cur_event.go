package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
)

type AlertCurEventController struct {
}

func (acec *AlertCurEventController) List(ctx *gin.Context) {

	dsType := ctx.Query("dsType")
	tid, _ := ctx.Get("TenantID")
	data, err := alertCurEventService.List(tid.(string), dsType)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}
