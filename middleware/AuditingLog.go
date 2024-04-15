package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"strings"
	"time"
	"watchAlert/controllers/response"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
	jwtUtils "watchAlert/public/utils/jwt"
)

func AuditingLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Operation user
		var username string
		createBy := jwtUtils.GetUser(context.Request.Header.Get("Authorization"))
		if createBy != "" {
			username = createBy
		} else {
			username = "用户未登录"
		}

		// Response log
		body := context.Request.Body
		readBody, err := io.ReadAll(body)
		if err != nil {
			globals.Logger.Sugar().Error(err)
			return
		}
		// 将 body 数据放回请求中
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(readBody))

		// 获取请求类型
		var reqTypeKey string
		// 获取 uri 的最后一位来定位审计类型
		splitAPI := strings.Split(context.Request.URL.Path, "/")
		if len(splitAPI) > 0 {
			reqTypeKey = splitAPI[len(splitAPI)-1]
		}

		tid := context.Request.Header.Get(TenantIDHeaderKey)
		if tid == "" {
			response.Fail(context, "租户ID不能为空", "failed")
			context.Abort()
			return
		}

		ps := models.PermissionsInfo()
		auditLog := models.AuditLog{
			TenantId:   tid,
			ID:         "Trace" + cmd.RandId(),
			Username:   username,
			IPAddress:  context.ClientIP(),
			Method:     context.Request.Method,
			Path:       context.Request.URL.Path,
			CreatedAt:  time.Now().Unix(),
			StatusCode: context.Writer.Status(),
			Body:       string(readBody),
			AuditType:  ps[reqTypeKey].Key,
		}

		err = auditLog.SaveAuditLog()
		if err != nil {
			response.Fail(context, "审计日志写入数据库失败, "+err.Error(), "failed")
			context.Abort()
			return
		}
	}
}
