package repo

import (
	"watchAlert/models"
	"watchAlert/public/globals"
)

type NoticeTmplRepo struct{}

func (nr NoticeTmplRepo) Search(r models.NoticeTemplateExampleQuery) ([]models.NoticeTemplateExample, error) {
	var data []models.NoticeTemplateExample
	var db = globals.DBCli.Model(&models.NoticeTemplateExample{})
	db.Where("name LIKE ? OR description LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%")
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
