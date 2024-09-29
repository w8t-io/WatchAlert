package repo

import (
	"fmt"
	"gorm.io/gorm"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
)

type (
	DashboardRepo struct {
		entryRepo
	}

	InterDashboardRepo interface {
		Create(d models.Dashboard) error
		Update(d models.Dashboard) error
		Delete(d models.DashboardQuery) error
		Search(d models.DashboardQuery) ([]models.Dashboard, error)
		CreateDashboardFolder(fd models.DashboardFolders) error
		UpdateDashboardFolder(fd models.DashboardFolders) error
		DeleteDashboardFolder(fd models.DashboardFolders) error
	}
)

func newDashboardInterface(db *gorm.DB, g InterGormDBCli) InterDashboardRepo {
	return &DashboardRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (dr DashboardRepo) Create(d models.Dashboard) error {
	err := dr.g.Create(&models.Dashboard{}, d)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) Update(d models.Dashboard) error {
	u := Updates{
		Table: &models.Dashboard{},
		Where: map[string]interface{}{
			"tenant_id = ?": d.TenantId,
			"id = ?":        d.ID,
		},
		Updates: d,
	}
	err := dr.g.Updates(u)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) Delete(d models.DashboardQuery) error {
	del := Delete{
		Table: &models.Dashboard{},
		Where: map[string]interface{}{
			"tenant_id = ?": d.TenantId,
			"id = ?":        d.ID,
		},
	}
	err := dr.g.Delete(del)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) Search(d models.DashboardQuery) ([]models.Dashboard, error) {
	var db = dr.db.Model(&models.Dashboard{})
	var data []models.Dashboard
	if d.Query != "" {
		db.Where("tenant_id = ? AND name LIKE ? OR description LIKE ? OR url LIKE ?", d.TenantId, "%"+d.Query+"%", "%"+d.Query+"%", "%"+d.Query+"%")
	} else {
		db.Where("tenant_id = ?", d.TenantId).Find(&data)
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (dr DashboardRepo) CreateDashboardFolder(fd models.DashboardFolders) error {
	fmt.Println("--->", fd)
	err := dr.g.Create(&models.DashboardFolders{}, fd)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) UpdateDashboardFolder(fd models.DashboardFolders) error {
	fmt.Println("--->", fd)
	u := Updates{
		Table: &models.DashboardFolders{},
		Where: map[string]interface{}{
			"tenant_id = ?": fd.TenantId,
			"id = ?":        fd.ID,
		},
		Updates: fd,
	}
	err := dr.g.Updates(u)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (dr DashboardRepo) DeleteDashboardFolder(fd models.DashboardFolders) error {
	d := Delete{
		Table: &models.DashboardFolders{},
		Where: map[string]interface{}{
			"tenant_id = ?": fd.TenantId,
			"id = ?":        fd.ID,
		},
	}
	err := dr.g.Delete(d)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}
