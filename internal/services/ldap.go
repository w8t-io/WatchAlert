package services

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/ldap.v2"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type ldapService struct {
	ctx *ctx.Context
}

type InterLdapService interface {
	ListUsers() ([]ldapUser, error)
	SyncUserToW8t()
	Login(username, password string) error
	SyncUsersCronjob()
}

func newInterLdapService(ctx *ctx.Context) InterLdapService {
	return &ldapService{
		ctx: ctx,
	}
}

func (l ldapService) getAdminAuth() (*ldap.Conn, error) {
	ls, err := ldap.Dial("tcp", global.Config.Ldap.Address)
	if err != nil {
		global.Logger.Sugar().Errorf("无法连接 LDAP 服务器, Address: %s, err: %s", global.Config.Ldap.Address, err.Error())
		return nil, err
	}

	err = ls.Bind(global.Config.Ldap.AdminUser, global.Config.Ldap.AdminPass)
	if err != nil {
		global.Logger.Sugar().Errorf("LDAP 管理员绑定失败 err: %s", err.Error())
		return nil, err
	}

	return ls, nil
}

type ldapUser struct {
	Uid    string `json:"uid"`
	Mobile string `json:"mobile"`
	Mail   string `json:"mail"`
}

func (l ldapService) ListUsers() ([]ldapUser, error) {
	lc := global.Config.Ldap
	searchRequest := ldap.NewSearchRequest(
		lc.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))",
		[]string{},
		nil,
	)

	auth, err := l.getAdminAuth()
	if err != nil {
		return nil, err
	}
	defer auth.Close()
	sr, err := auth.Search(searchRequest)
	if err != nil {
		global.Logger.Sugar().Errorf("LDAP 用户搜索失败, err: %s", err.Error())
		return nil, err
	}

	var users []ldapUser
	for _, entry := range sr.Entries {
		uid := entry.GetAttributeValue("uid")
		if uid == "" {
			continue
		}
		users = append(users, ldapUser{
			Uid:    entry.GetAttributeValue("uid"),
			Mobile: entry.GetAttributeValue("mobile"),
			Mail:   entry.GetAttributeValue("mail"),
		})
	}

	return users, nil
}

func (l ldapService) SyncUserToW8t() {
	users, err := l.ListUsers()
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return
	}

	for _, u := range users {
		_, b, _ := l.ctx.DB.User().Get(models.MemberQuery{Query: u.Mail})
		if b {
			continue
		}
		uid := tools.RandUid()
		m := models.Member{
			UserId:   uid,
			UserName: u.Uid,
			Email:    u.Mail,
			Phone:    u.Mobile,
			CreateBy: "LDAP",
			CreateAt: time.Now().Unix(),
			Tenants:  []string{"default"},
		}
		err = l.ctx.DB.User().Create(m)
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}

		err = l.ctx.DB.Tenant().AddTenantLinkedUsers(models.TenantLinkedUsers{
			ID:       "default",
			UserRole: global.Config.Ldap.DefaultUserRole,
			Users: []models.TenantUser{
				{
					UserID:   uid,
					UserName: u.Mail,
				},
			},
		})
		if err != nil {
			global.Logger.Sugar().Errorf(err.Error())
			return
		}
	}
}

func (l ldapService) Login(username, password string) error {
	auth, err := l.getAdminAuth()
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return err
	}

	userDn := fmt.Sprintf("%s=%s,%s", global.Config.Ldap.UserPrefix, username, global.Config.Ldap.UserDN)
	err = auth.Bind(userDn, password)
	if err != nil {
		global.Logger.Sugar().Errorf("LDAP 用户登陆失败, err: %s", err.Error())
		return err
	}

	return nil
}

func (l ldapService) SyncUsersCronjob() {
	c := cron.New()
	_, err := c.AddFunc(global.Config.Ldap.Cronjob, func() {
		l.SyncUserToW8t()
	})
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return
	}
	c.Start()
	defer c.Stop()

	select {}
}
