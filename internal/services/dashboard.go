package services

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
)

type dashboardService struct {
	ctx *ctx.Context
}

type InterDashboardService interface {
	List(req interface{}) (data interface{}, error interface{})
	Get(req interface{}) (data interface{}, error interface{})
	Create(req interface{}) (data interface{}, error interface{})
	Update(req interface{}) (data interface{}, error interface{})
	Delete(req interface{}) (data interface{}, error interface{})
	Search(req interface{}) (data interface{}, error interface{})
}

func newInterDashboardService(ctx *ctx.Context) InterDashboardService {
	return &dashboardService{
		ctx: ctx,
	}
}

func (ds dashboardService) List(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	var d []models.Dashboard
	var db = ds.ctx.DB.DB().Model(&models.Dashboard{})
	err := db.Where("tenant_id = ?", r.TenantId).Find(&d).Error
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (ds dashboardService) Get(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	var d models.Dashboard
	var db = ds.ctx.DB.DB().Model(&models.Dashboard{})
	err := db.Where("tenant_id = ? AND id = ?", r.TenantId, r.ID).First(&d).Error
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (ds dashboardService) Create(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.Dashboard)
	r.ID = "db" + cmd.RandId()
	err := ds.ctx.DB.Dashboard().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) Update(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.Dashboard)
	err := ds.ctx.DB.Dashboard().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) Delete(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	err := ds.ctx.DB.Dashboard().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) Search(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	data, err := ds.ctx.DB.Dashboard().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
