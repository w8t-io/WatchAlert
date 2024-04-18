package services

import (
	"watchAlert/controllers/repo"
	"watchAlert/models"
)

type NoticeTmplService struct {
	repo.NoticeRepo
	repo.NoticeTmplRepo
}

type InterNoticeTmplService interface {
	Search(req interface{}) (interface{}, interface{})
}

func NewInterNoticeTmplService() InterNoticeTmplService {
	return &NoticeTmplService{}
}

func (nts NoticeTmplService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeTemplateExampleQuery)
	data, err := nts.NoticeTmplRepo.Search(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}
