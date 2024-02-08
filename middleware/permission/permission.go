package permission

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
	utils "watchAlert/utils/jwt"
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

		// 获取当前用户
		var user models.Member
		err := globals.DBCli.Model(&models.Member{}).Where("user_id = ?", userId).First(&user).Error
		if err != nil {
			response.PermissionFail(context)
			globals.Logger.Sugar().Error("获取当前用户失败 ->", err.Error())
			context.Abort()
			return
		}

		var (
			role       models.UserRole
			permission []models.UserPermissions
		)
		// 根据用户角色获取权限
		err = globals.DBCli.Model(&models.UserRole{}).Where("name = ?", user.Role).First(&role).Error
		if err != nil {
			response.PermissionFail(context)
			globals.Logger.Sugar().Errorf("获取用户 %s 的角色失败 -> %s", user.UserName, err.Error())
			context.Abort()
			return
		}
		_ = json.Unmarshal([]byte(role.Permissions), &permission)

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
