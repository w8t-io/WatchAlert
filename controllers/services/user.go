package services

import (
	"watchAlert/controllers/repo"
	"watchAlert/models"
)

type UserService struct {
	repo.UserRepo
}

type InterUserService interface {
	Search(req interface{}) (interface{}, interface{})
}

func NewInterUserService() InterUserService {
	return &UserService{}
}

func (us UserService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MemberQuery)
	data, err := us.UserRepo.SearchUser(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}
