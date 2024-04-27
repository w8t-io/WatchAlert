package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	UserPermissionsRepo struct {
		entryRepo
	}

	InterUserPermissionsRepo interface {
		List() ([]models.UserPermissions, error)
	}
)

func newInterUserPermissionsRepo(db *gorm.DB, g InterGormDBCli) InterUserPermissionsRepo {
	return &UserPermissionsRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (up UserPermissionsRepo) List() ([]models.UserPermissions, error) {
	var data []models.UserPermissions
	err := up.db.Model(&models.UserPermissions{}).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}
