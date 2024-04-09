package services

import (
	"watchAlert/controllers/repo"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type DashboardService struct {
	repo.DashboardRepo
}

type InterDashboardService interface {
	ListDashboard(tid string) (data interface{}, error interface{})
	GetDashboard(req interface{}) (data interface{}, error interface{})
	CreateDashboard(req interface{}) (data interface{}, error interface{})
	UpdateDashboard(req interface{}) (data interface{}, error interface{})
	DeleteDashboard(req interface{}) (data interface{}, error interface{})
	SearchDashboard(req interface{}) (data interface{}, error interface{})
}

func NewInterDashboardService() InterDashboardService {
	return &DashboardService{}
}

func (ds DashboardService) ListDashboard(tid string) (data interface{}, error interface{}) {
	var d []models.Dashboard
	var db = globals.DBCli.Model(&models.Dashboard{})
	err := db.Where("tenant_id = ?", tid).Find(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (DashboardService) GetDashboard(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	var d models.Dashboard
	var db = globals.DBCli.Model(&models.Dashboard{})
	err := db.Where("tenant_id = ? AND id = ?", r.TenantId, r.ID).First(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (ds DashboardService) CreateDashboard(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.Dashboard)
	r.ID = "db" + cmd.RandId()
	err := ds.DashboardRepo.CreateDashboard(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ds DashboardService) UpdateDashboard(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.Dashboard)
	err := ds.DashboardRepo.UpdateDashboard(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ds DashboardService) DeleteDashboard(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	err := ds.DashboardRepo.DeleteDashboard(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ds DashboardService) SearchDashboard(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardQuery)
	data, err := ds.DashboardRepo.SearchDashboard(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}
