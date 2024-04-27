package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	UserRoleRepo struct {
		entryRepo
	}

	InterUserRoleRepo interface {
		List(r models.UserRoleQuery) ([]models.UserRole, error)
		Create(r models.UserRole) error
		Update(r models.UserRole) error
		Delete(r models.UserRoleQuery) error
	}
)

func newUserRoleInterface(db *gorm.DB, g InterGormDBCli) InterUserRoleRepo {
	return &UserRoleRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (ur UserRoleRepo) List(r models.UserRoleQuery) ([]models.UserRole, error) {
	var data []models.UserRole
	err := ur.db.Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (ur UserRoleRepo) Create(r models.UserRole) error {
	err := ur.g.Create(models.UserRole{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRoleRepo) Update(r models.UserRole) error {
	u := Updates{
		Table: models.UserRole{},
		Where: map[string]interface{}{
			"id = ?": r.ID,
		},
		Updates: r,
	}

	err := ur.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRoleRepo) Delete(r models.UserRoleQuery) error {
	d := Delete{
		Table: models.UserRole{},
		Where: map[string]interface{}{
			"id = ?": r.ID,
		},
	}

	err := ur.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
