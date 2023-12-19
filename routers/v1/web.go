package v1

import "github.com/gin-gonic/gin"

func Web(r *gin.Engine) {

	r.LoadHTMLGlob("web/**/*")

	r.Static("/web/static", "web/static")

	r.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{})
	})

	r.GET("/ruleGroup", func(context *gin.Context) {
		context.HTML(200, "ruleGroup.html", gin.H{})
	})
	r.GET("/ruleGroup/:ruleGroup/rule", func(context *gin.Context) {
		context.HTML(200, "rule.html", gin.H{})
	})
	r.GET("/noticeObject", func(context *gin.Context) {
		context.HTML(200, "noticeObject.html", gin.H{})
	})
	r.GET("/dutyManage", func(context *gin.Context) {
		context.HTML(200, "duty_manage.html", gin.H{})
	})
	r.GET("/dutyManage/:dutyId/schedule", func(context *gin.Context) {
		context.HTML(200, "schedule.html", gin.H{})
	})

}
