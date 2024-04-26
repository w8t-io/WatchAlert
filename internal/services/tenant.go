package services

import (
	"fmt"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
)

type tenantService struct {
	ctx *ctx.Context
}

type InterTenantService interface {
	Create(req interface{}) (data interface{}, err interface{})
	Update(req interface{}) (data interface{}, err interface{})
	Delete(req interface{}) (data interface{}, err interface{})
	List() (data interface{}, err interface{})
}

func newInterTenantService(ctx *ctx.Context) InterTenantService {
	return &tenantService{
		ctx: ctx,
	}
}

func (ts tenantService) Create(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.Tenant)
	nt := models.Tenant{
		ID:               "tid-" + cmd.RandId(),
		Name:             r.Name,
		CreateAt:         time.Now().Unix(),
		CreateBy:         r.CreateBy,
		Manager:          r.Manager,
		Description:      r.Description,
		RuleNumber:       r.RuleNumber,
		UserNumber:       r.UserNumber,
		DutyNumber:       r.DutyNumber,
		NoticeNumber:     r.NoticeNumber,
		RemoveProtection: r.RemoveProtection,
	}

	err = ts.ctx.DB.Tenant().Create(nt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts tenantService) Update(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.Tenant)
	nt := *r

	err = ts.ctx.DB.Tenant().Update(nt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts tenantService) Delete(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.TenantQuery)

	var t models.Tenant
	ts.ctx.DB.DB().Model(&models.Tenant{}).Where("id = ?", r.ID).Find(&t)

	if *t.RemoveProtection {
		return nil, fmt.Errorf("删除失败, 删除保护已开启 关闭后再删除")
	}

	err = ts.ctx.DB.Tenant().Delete(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts tenantService) List() (data interface{}, err interface{}) {
	data, err = ts.ctx.DB.Tenant().List()
	if err != nil {
		return nil, err
	}
	return data, err
}
