package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/middleware"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
	jwtUtils "watchAlert/public/utils/jwt"
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
		userB.GET("searchDutyUser", uc.SearchDutyUser)
		userB.GET("searchUser", uc.Search)
	}

}

func (uc UserController) Login(ctx *gin.Context) {

	var (
		use models.Member
		req models.Member
	)
	_ = ctx.ShouldBindJSON(&req)

	// 校验 Password
	arr := md5.Sum([]byte(req.Password))
	hashPassword := hex.EncodeToString(arr[:])

	// 查询用户信息
	err := globals.DBCli.Where("user_name = ?", req.UserName).First(&use).Error
	if err == gorm.ErrRecordNotFound || hashPassword != use.Password {
		response.Fail(ctx, "用户不存在或密码错误", "failed")
		return
	}

	req.UserId = use.UserId
	req.Password = hashPassword
	tokenData, err := jwtUtils.GenerateToken(req)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	duration := time.Duration(globals.Config.Jwt.Expire) * time.Second
	globals.RedisCli.Set("uid-"+use.UserId, cmd.JsonMarshal(use), duration)

	response.Success(ctx, tokenData, "success")

}

func (uc UserController) Register(ctx *gin.Context) {

	var (
		dataUser  models.Member
		parseUser models.Member
	)
	_ = ctx.ShouldBindJSON(&parseUser)

	globals.DBCli.Where("user_name = ?", parseUser.UserName).First(&dataUser)
	if dataUser.UserName != "" {
		response.Fail(ctx, "用户已存在", "failed")
		return
	}

	arr := md5.Sum([]byte(parseUser.Password))
	hashPassword := hex.EncodeToString(arr[:])
	parseUser.UserId = cmd.RandUid()
	parseUser.Password = hashPassword
	parseUser.CreateAt = time.Now().Unix()
	createUser := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	if createUser == "" {
		createUser = "system"
	}
	parseUser.CreateBy = createUser

	err := repo.DBCli.Create(models.Member{}, &parseUser)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, nil, "success")
}

func (uc UserController) Update(ctx *gin.Context) {

	var (
		user   models.Member
		dbUser models.Member
	)
	_ = ctx.ShouldBindJSON(&user)

	globals.DBCli.Model(&models.Member{}).Where("user_id = ?", user.UserId).First(&dbUser)

	user.Password = dbUser.Password
	err := repo.DBCli.Updates(repo.Updates{
		Table:   models.Member{},
		Where:   []interface{}{"user_id = ?", user.UserId},
		Updates: user,
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	uc.changeCache(user.UserId)

	response.Success(ctx, nil, "success")

}

func (uc UserController) CheckUser(ctx *gin.Context) {

	var member models.Member

	username := ctx.Query("username")

	err := globals.DBCli.Where("user_name = ?", username).First(&member).Error
	if err != nil {
		response.Fail(ctx, err, "failed")
		return
	}

	response.Success(ctx, member, "success")

}

func (uc UserController) List(ctx *gin.Context) {

	var userList []models.Member
	err := globals.DBCli.Find(&userList).Error
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, userList, "success")

}

func (uc UserController) ChangePass(ctx *gin.Context) {

	id := ctx.Query("userid")
	var data struct {
		Password string `json:"password"`
	}
	_ = ctx.ShouldBindJSON(&data)
	arr := md5.Sum([]byte(data.Password))
	hashPassword := hex.EncodeToString(arr[:])

	err := repo.DBCli.Update(repo.Update{
		Table:  models.Member{},
		Where:  []interface{}{"user_id = ?", id},
		Update: []string{"password", hashPassword},
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	uc.changeCache(id)

	response.Success(ctx, "", "success")

}

func (uc UserController) Delete(ctx *gin.Context) {

	id := ctx.Query("userid")

	err := repo.DBCli.Delete(repo.Delete{
		Table: models.Member{},
		Where: []interface{}{"user_id = ?", id},
	})
	if err != nil {
		return
	}

	response.Success(ctx, "userid: "+id, "success")

}

func (uc UserController) SearchDutyUser(ctx *gin.Context) {

	var data []models.Member
	err := globals.DBCli.Where("join_duty = ?", "true").Find(&data).Error
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}

func (uc UserController) GetUserInfo(ctx *gin.Context) {

	token := ctx.Request.Header.Get("Authorization")
	code, ok := jwtUtils.IsTokenValid(token)
	if !ok {
		if code == 401 {
			response.TokenFail(ctx)
			return
		}
	}

	username := jwtUtils.GetUser(token)

	userInfo := models.Member{}
	globals.DBCli.Model(&models.Member{}).Where("user_name = ?", username).Find(&userInfo)
	response.Success(ctx, userInfo, "success")

}

func (uc UserController) changeCache(userId string) {

	var dbUser models.Member
	globals.DBCli.Model(&models.Member{}).Where("user_id = ?", userId).First(&dbUser)

	var cacheUser models.Member
	result, err := globals.RedisCli.Get("uid-" + userId).Result()
	if err != nil {
		globals.Logger.Sugar().Error(err)
	}
	_ = json.Unmarshal([]byte(result), &cacheUser)

	duration, _ := globals.RedisCli.TTL("uid-" + userId).Result()
	globals.RedisCli.Set("uid-"+userId, cmd.JsonMarshal(dbUser), duration)

}

func (uc UserController) Search(ctx *gin.Context) {
	r := new(models.MemberQuery)
	BindQuery(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return userService.Search(r)
	})
}
