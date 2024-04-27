package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
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
	db.Where("name LIKE ? OR description LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%")
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (nr NoticeTmplRepo) Create(r models.NoticeTemplateExample) error {
	nt := r
	nt.Id = "nt-" + cmd.RandId()
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
