package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
)

type (
	TenantRepo struct {
		entryRepo
	}

	InterTenantRepo interface {
		Create(t models.Tenant) error
		Update(t models.Tenant) error
		Delete(t models.TenantQuery) error
		List() (data []models.Tenant, err error)
	}
)

func newTenantInterface(db *gorm.DB, g InterGormDBCli) InterTenantRepo {
	return &TenantRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (tr TenantRepo) Create(t models.Tenant) error {
	err := tr.g.Create(&models.Tenant{}, t)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (tr TenantRepo) Update(t models.Tenant) error {
	u := Updates{
		Table: &models.Tenant{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
		Updates: t,
	}
	err := tr.g.Updates(u)
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (tr TenantRepo) Delete(t models.TenantQuery) error {
	err := tr.g.Delete(Delete{
		Table: &models.Tenant{},
		Where: map[string]interface{}{
			"id = ?": t.ID,
		},
	})
	if err != nil {
		global.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (tr TenantRepo) List() (data []models.Tenant, err error) {
	var d []models.Tenant
	err = tr.db.Model(&models.Tenant{}).Find(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}
