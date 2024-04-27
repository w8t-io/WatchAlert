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

		var (
			role       models.UserRole
			permission []models.UserPermissions
		)
		// 根据用户角色获取权限
		err = c.DB.DB().Model(&models.UserRole{}).Where("name = ?", user.Role).First(&role).Error
		if err != nil {
			response.Fail(context, fmt.Sprintf("获取用户 %s 的角色失败, 角色名称: %s", user.UserName, user.Role), "failed")
			global.Logger.Sugar().Errorf("获取用户 %s 的角色失败 -> %s", user.UserName, err.Error())
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
