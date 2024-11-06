package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

type (
	NoticeTmplRepo struct {
		entryRepo
	}

	InterNoticeTmplRepo interface {
		List(r models.NoticeTemplateExampleQuery) ([]models.NoticeTemplateExample, error)
		Search(r models.NoticeTemplateExampleQuery) ([]models.NoticeTemplateExample, error)
		Create(r models.NoticeTemplateExample) error
		Update(r models.NoticeTemplateExample) error
		Delete(r models.NoticeTemplateExampleQuery) error
		Get(r models.NoticeTemplateExampleQuery) models.NoticeTemplateExample
	}
)

func newNoticeTmplInterface(db *gorm.DB, g InterGormDBCli) InterNoticeTmplRepo {
	return &NoticeTmplRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (nr NoticeTmplRepo) List(r models.NoticeTemplateExampleQuery) ([]models.NoticeTemplateExample, error) {
	var (
		data []models.NoticeTemplateExample
		db   = nr.db.Model(&models.NoticeTemplateExample{})
	)
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (nr NoticeTmplRepo) Search(r models.NoticeTemplateExampleQuery) ([]models.NoticeTemplateExample, error) {
	var (
		data []models.NoticeTemplateExample
		db   = nr.db.Model(&models.NoticeTemplateExample{})
	)
	if r.Id != "" {
		db.Where("id = ?", r.Id)
	}

	if r.NoticeType != "" {
		db.Where("notice_type = ?", r.NoticeType)
	}

	if r.Query != "" {
		db.Where("name LIKE ? OR description LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%")
	}
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (nr NoticeTmplRepo) Create(r models.NoticeTemplateExample) error {
	nt := r
	nt.Id = "nt-" + tools.RandId()
	err := nr.g.Create(models.NoticeTemplateExample{}, nt)
	if err != nil {
		return err
	}

	return nil
}

func (nr NoticeTmplRepo) Update(r models.NoticeTemplateExample) error {
	u := Updates{
		Table: models.NoticeTemplateExample{},
		Where: map[string]interface{}{
			"id = ?": r.Id,
		},
		Updates: r,
	}

	err := nr.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (nr NoticeTmplRepo) Delete(r models.NoticeTemplateExampleQuery) error {
	d := Delete{
		Table: models.NoticeTemplateExample{},
		Where: map[string]interface{}{
			"id = ?": r.Id,
		},
	}

	err := nr.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}

func (nr NoticeTmplRepo) Get(r models.NoticeTemplateExampleQuery) models.NoticeTemplateExample {
	var (
		data models.NoticeTemplateExample
		db   = nr.db.Model(&models.NoticeTemplateExample{})
	)
	if r.Id != "" {
		db.Where("id = ?", r.Id)
	}

	err := db.First(&data).Error
	if err != nil {
		return data
	}
	return data
}
