package v1

import "github.com/gin-gonic/gin"

func Web(r *gin.Engine) {

	r.LoadHTMLGlob("web/**/*")

	r.Static("/web/static", "web/static")

	r.GET("/ruleGroup", func(context *gin.Context) {
		context.HTML(200, "ruleGroup.html", gin.H{})
	})
	r.GET("/ruleGroup/:ruleGroup/rule", func(context *gin.Context) {
		context.HTML(200, "rule.html", gin.H{})
	})

}
