package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"watchAlert/internal/global"
	models "watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/response"
	"watchAlert/pkg/utils/cmd"
	"watchAlert/pkg/utils/jwt"
)

func Permission() gin.HandlerFunc {
	return func(context *gin.Context) {
		tid := context.Request.Header.Get(TenantIDHeaderKey)
		if tid == "null" || tid == "" {
			return
		}
		// 获取 Token
		tokenStr := context.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(context)
			context.Abort()
			return
		}

		// Bearer Token, 获取 Token 值
		tokenStr = tokenStr[len(utils.TokenType)+1:]

		userId := utils.GetUserID(tokenStr)

		c := ctx.DO()
		// 获取当前用户
		var user models.Member
		err := c.DB.DB().Model(&models.Member{}).Where("user_id = ?", userId).First(&user).Error
		if gorm.ErrRecordNotFound == err {
			global.Logger.Sugar().Errorf("用户不存在, uid: %s", userId)
		}
		if err != nil {
			response.PermissionFail(context)
			context.Abort()
			return
		}

		context.Set("UserId", user.UserId)
		context.Set("UserEmail", user.Email)

		// 获取租户用户角色
		tenantUserInfo, _ := c.DB.Tenant().GetTenantLinkedUserInfo(models.GetTenantLinkedUserInfo{ID: tid, UserID: userId})
		if err != nil {
			global.Logger.Sugar().Errorf("获取租户用户角色失败 %s", err.Error())
			response.TokenFail(context)
			context.Abort()
			return
		}

		var (
			role       models.UserRole
			permission []models.UserPermissions
		)
		// 根据用户角色获取权限
		err = c.DB.DB().Model(&models.UserRole{}).Where("id = ?", tenantUserInfo.UserRole).First(&role).Error
		if err != nil {
			response.Fail(context, fmt.Sprintf("获取用户 %s 的角色失败, %s %s", user.UserName, tenantUserInfo.UserRole, err.Error()), "failed")
			global.Logger.Sugar().Errorf("获取用户 %s 的角色失败 %s %s", user.UserName, tenantUserInfo.UserRole, err.Error())
			context.Abort()
			return
		}
		_ = json.Unmarshal([]byte(cmd.JsonMarshal(role.Permissions)), &permission)

		urlPath := context.Request.URL.Path

		var pass bool
		for _, v := range permission {
			if urlPath == v.API {
				pass = true
				break
			}
		}
		if !pass {
			response.PermissionFail(context)
			context.Abort()
			return
		}
	}
}
