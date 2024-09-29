package services

import (
	"encoding/json"
	"fmt"
	"io"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
	"watchAlert/pkg/utils/http"
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
	ListFolder(req interface{}) (data interface{}, error interface{})
	GetFolder(req interface{}) (data interface{}, error interface{})
	CreateFolder(req interface{}) (data interface{}, error interface{})
	UpdateFolder(req interface{}) (data interface{}, error interface{})
	DeleteFolder(req interface{}) (data interface{}, error interface{})
	ListGrafanaDashboards(req interface{}) (data interface{}, error interface{})
	GetDashboardFullUrl(req interface{}) (data interface{}, error interface{})
}

func newInterDashboardService(ctx *ctx.Context) InterDashboardService {
	return &dashboardService{
		ctx: ctx,
	}
}

func (ds dashboardService) ListFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardFolders)
	var f []models.DashboardFolders
	var db = ds.ctx.DB.DB().Model(&models.DashboardFolders{})
	db.Where("tenant_id = ?", r.TenantId)
	err := db.Find(&f).Error
	if err != nil {
		return nil, err
	}

	return f, nil
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

func (ds dashboardService) GetFolder(req interface{}) (data interface{}, error interface{}) {
	var f models.DashboardFolders
	r := req.(*models.DashboardFolders)

	var db = ctx.DB.DB().Model(&models.DashboardFolders{})
	db.Where("id = ?", r.ID)
	err := db.First(&f).Error
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (ds dashboardService) CreateFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardFolders)
	r.ID = "f-" + cmd.RandId()
	err := ctx.DB.Dashboard().CreateDashboardFolder(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) UpdateFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardFolders)
	err := ctx.DB.Dashboard().UpdateDashboardFolder(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) DeleteFolder(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardFolders)
	err := ctx.DB.Dashboard().DeleteDashboardFolder(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds dashboardService) ListGrafanaDashboards(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardFolders)
	get, err := http.Get(nil, fmt.Sprintf("%s/api/search?folderIds=%d", r.GrafanaHost, r.GrafanaFolderId))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, err
	}

	var d []models.GrafanaDashboardInfo
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (ds dashboardService) GetDashboardFullUrl(req interface{}) (data interface{}, error interface{}) {
	r := req.(*models.DashboardFolders)
	get, err := http.Get(nil, fmt.Sprintf("%s/api/dashboards/uid/%s", r.GrafanaHost, r.GrafanaDashboardUid))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, err
	}

	var d models.GrafanaDashboardMeta
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}

	full := r.GrafanaHost + d.Meta.Url + "?theme=" + r.Theme
	return full, nil
}
