package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
	jwtUtils "watchAlert/pkg/utils/jwt"
)

type userService struct {
	ctx *ctx.Context
}

type InterUserService interface {
	Search(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	Login(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Register(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	ChangePass(req interface{}) (interface{}, interface{})
}

func newInterUserService(ctx *ctx.Context) InterUserService {
	return &userService{
		ctx: ctx,
	}
}

func (us userService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	data, err := us.ctx.DB.User().Search(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (us userService) List(req interface{}) (interface{}, interface{}) {
	data, err := us.ctx.DB.User().List()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (us userService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	data, _, err := us.ctx.DB.User().Get(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (us userService) Login(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)

	// 校验 Password
	arr := md5.Sum([]byte(r.Password))
	hashPassword := hex.EncodeToString(arr[:])

	q := models.MemberQuery{
		UserName: r.UserName,
	}
	data, _, err := us.ctx.DB.User().Get(q)
	if err != nil {
		return nil, err
	}

	if data.Password != hashPassword {
		return nil, fmt.Errorf("密码错误")
	}

	r.UserId = data.UserId
	r.Password = hashPassword
	tokenData, err := jwtUtils.GenerateToken(*r)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(global.Config.Jwt.Expire) * time.Second
	us.ctx.Redis.Redis().Set("uid-"+data.UserId, cmd.JsonMarshal(r), duration)

	return tokenData, nil
}

func (us userService) Register(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)

	q := models.MemberQuery{UserName: r.UserName}
	_, ok, _ := us.ctx.DB.User().Get(q)
	if ok {
		return nil, fmt.Errorf("用户已存在")
	}

	arr := md5.Sum([]byte(r.Password))
	hashPassword := hex.EncodeToString(arr[:])
	// 在初始化admin用户时会固定一个userid，所以这里需要做一下判断；
	if r.UserId == "" {
		r.UserId = cmd.RandUid()
	}
	r.Password = hashPassword
	r.CreateAt = time.Now().Unix()

	if r.CreateBy == "" {
		r.CreateBy = "system"
	}

	err := us.ctx.DB.User().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (us userService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)
	var dbData models.Member

	db := us.ctx.DB.DB().Model(models.Member{})
	db.Where("user_id = ?", r.UserId).First(&dbData)

	r.Password = dbData.Password
	err := us.ctx.DB.User().Update(*r)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	err := us.ctx.DB.User().Delete(*r)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) ChangePass(req interface{}) (interface{}, interface{}) {
	r := req.(*models.Member)

	arr := md5.Sum([]byte(r.Password))
	hashPassword := hex.EncodeToString(arr[:])
	r.Password = hashPassword

	err := us.ctx.DB.User().ChangePass(*r)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}
