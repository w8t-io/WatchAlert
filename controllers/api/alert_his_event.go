package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
)

type AlertHisEventController struct {
}

func (ahec *AlertHisEventController) List(ctx *gin.Context) {

	data, err := alertHisEventService.List(ctx)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}
