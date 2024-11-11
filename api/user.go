package api

import (
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	jwtUtils "watchAlert/pkg/tools"
)

type UserController struct{}

/*
	用户 API
	/api/w8t/user
*/
func (uc UserController) API(gin *gin.RouterGroup) {

	userA := gin.Group("user")
	userA.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		userA.POST("userUpdate", uc.Update)
		userA.POST("userDelete", uc.Delete)
		userA.POST("userChangePass", uc.ChangePass)
	}

	userB := gin.Group("user")
	userB.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		userB.GET("userList", uc.List)
		userB.GET("searchDutyUser", uc.Search)
		userB.GET("searchUser", uc.Search)
	}

}

func (uc UserController) List(ctx *gin.Context) {
	r := new(models.MemberQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.List(r)
	})
}

func (uc UserController) Search(ctx *gin.Context) {
	r := new(models.MemberQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Search(r)
	})
}

func (uc UserController) Get(ctx *gin.Context) {
	r := new(models.MemberQuery)
	token := ctx.Request.Header.Get("Authorization")
	username := jwtUtils.GetUser(token)
	r.UserName = username

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Get(r)
	})
}

func (uc UserController) Login(ctx *gin.Context) {
	r := new(models.Member)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Login(r)
	})
}

func (uc UserController) Register(ctx *gin.Context) {
	r := new(models.Member)
	BindJson(ctx, r)

	createUser := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	r.CreateBy = createUser

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Register(r)
	})
}

func (uc UserController) Update(ctx *gin.Context) {
	r := new(models.Member)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Update(r)
	})
}

func (uc UserController) Delete(ctx *gin.Context) {
	r := new(models.MemberQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Delete(r)
	})
}

func (uc UserController) CheckUser(ctx *gin.Context) {
	r := new(models.MemberQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Get(r)
	})
}

func (uc UserController) ChangePass(ctx *gin.Context) {
	r := new(models.Member)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.ChangePass(r)
	})
}
