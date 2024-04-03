package services

import (
	"fmt"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type TenantService struct {
	repo repo.TenantRepo
}

type InterTenantService interface {
	CreateTenant(req interface{}) (data interface{}, err interface{})
	UpdateTenant(req interface{}) (data interface{}, err interface{})
	DeleteTenant(req interface{}) (data interface{}, err interface{})
	ListTenant() (data interface{}, err interface{})
}

func NewInterTenantService() InterTenantService {
	return TenantService{}
}

func (ts TenantService) CreateTenant(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.Tenant)
	nt := models.Tenant{
		ID:           "tid-" + cmd.RandId(),
		Name:         r.Name,
		CreateAt:     time.Now().Unix(),
		CreateBy:     r.CreateBy,
		Manager:      r.Manager,
		Description:  r.Description,
		RuleNumber:   r.RuleNumber,
		UserNumber:   r.UserNumber,
		DutyNumber:   r.DutyNumber,
		NoticeNumber: r.NoticeNumber,
	}

	err = ts.repo.CreateTenant(nt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts TenantService) UpdateTenant(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.Tenant)
	nt := *r

	err = ts.repo.UpdateTenant(nt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts TenantService) DeleteTenant(req interface{}) (data interface{}, err interface{}) {
	r := req.(*models.Tenant)

	var t models.Tenant
	globals.DBCli.Model(&models.Tenant{}).Where("id = ?", r.ID).Find(&t)
	if *t.RemoveProtection {
		return nil, fmt.Errorf("删除失败, 删除保护已开启 关闭后再删除")
	}

	err = ts.repo.DeleteTenant(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ts TenantService) ListTenant() (data interface{}, err interface{}) {
	data, err = ts.repo.ListTenant()
	if err != nil {
		return nil, err
	}
	return data, err
}
