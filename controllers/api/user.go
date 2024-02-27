package api

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
	jwtUtils "watchAlert/utils/jwt"
)

type UserController struct{}

func (u *UserController) Login(ctx *gin.Context) {

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

	tokenData, err := jwtUtils.GenerateToken(use.UserId, use.UserName)
	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}
	response.Success(ctx, tokenData, "success")

}

func (u *UserController) Register(ctx *gin.Context) {

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

func (u *UserController) Update(ctx *gin.Context) {

	var user models.Member
	_ = ctx.ShouldBindJSON(&user)

	err := repo.DBCli.Updates(repo.Updates{
		Table:   models.Member{},
		Where:   []string{"user_id = ?", user.UserId},
		Updates: user,
	})
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, nil, "success")

}

func (u *UserController) CheckUser(ctx *gin.Context) {

	var member models.Member

	username := ctx.Query("username")

	err := globals.DBCli.Where("user_name = ?", username).First(&member).Error
	if err != nil {
		response.Fail(ctx, err, "failed")
		return
	}

	response.Success(ctx, member, "success")

}

func (u *UserController) List(ctx *gin.Context)  {

	var userList []models.Member
	err := globals.DBCli.Find(&userList).Error
	if err != nil{
		response.Fail(ctx,err.Error(),"failed")
		return
	}

	response.Success(ctx,userList,"success")
	
}

func (u *UserController) ChangePass(ctx *gin.Context) {

	id := ctx.Query("userid")
	var data struct{
		Password string `json:"password"`
	}
    _ = ctx.ShouldBindJSON(&data)
	arr := md5.Sum([]byte(data.Password))
	hashPassword := hex.EncodeToString(arr[:])

	err := repo.DBCli.Update(repo.Update{
		Table:  models.Member{},
		Where:  []string{"user_id = ?", id},
		Update: []string{"password",hashPassword},
	})
	if err != nil {
		response.Fail(ctx,err.Error(),"failed")
		return
	}

	response.Success(ctx,"","success")

}

func (u *UserController) Delete(ctx *gin.Context)  {

	id := ctx.Query("userid")

	err := repo.DBCli.Delete(repo.Delete{
		Table: models.Member{},
		Where: []string{"user_id = ?", id},
	})
	if err != nil {
		return
	}

	response.Success(ctx, "userid: "+id, "success")

}

func (u *UserController) SearchDutyUser(ctx *gin.Context) {

	var data []models.Member
	err := globals.DBCli.Where("join_duty = ?", "true").Find(&data).Error
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, data, "success")

}

func (u *UserController) GetUserInfo(ctx *gin.Context) {

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