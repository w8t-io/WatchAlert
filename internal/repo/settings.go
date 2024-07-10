package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	settingRepo struct {
		entryRepo
	}

	InterSettingRepo interface {
		Create(r models.Settings) error
		Update(r models.Settings) error
		Get() (models.Settings, error)
		Check() bool
	}
)

func newSettingRepoInterface(db *gorm.DB, g InterGormDBCli) InterSettingRepo {
	return settingRepo{
		entryRepo{
			db: db,
			g:  g,
		},
	}
}

func (a settingRepo) Create(r models.Settings) error {
	r.IsInit = 1
	err := a.g.Create(models.Settings{}, r)
	if err != nil {
		return err
	}

	return nil
}

func (a settingRepo) Update(r models.Settings) error {
	err := a.g.Updates(
		Updates{
			Table: models.Settings{},
			Where: map[string]interface{}{
				"is_init = ?": 1,
			},
			Updates: r,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (a settingRepo) Get() (models.Settings, error) {
	var data models.Settings
	db := a.db.Model(models.Settings{})
	db.Where("is_init = ?", 1)
	db.First(&data)

	return data, nil
}

func (a settingRepo) Check() bool {
	var data models.Settings
	db := a.db.Model(models.Settings{})
	db.Where("is_init = ?", 1)
	err := db.First(&data).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return false
	}

	return true
}
