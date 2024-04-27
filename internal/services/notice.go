package services

import (
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/cmd"
)

type noticeService struct {
	ctx *ctx.Context
}

type InterNoticeService interface {
	List(req interface{}) (interface{}, interface{})
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	Check(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
}

func newInterAlertNoticeService(ctx *ctx.Context) InterNoticeService {
	return &noticeService{
		ctx,
	}
}

func (n noticeService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeQuery)
	data, err := n.ctx.DB.Notice().List(*r)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (n noticeService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertNotice)
	ok := n.ctx.DB.Notice().GetQuota(r.TenantId)
	if !ok {
		return models.AlertNotice{}, fmt.Errorf("创建失败, 配额不足")
	}

	r.Uuid = "n-" + cmd.RandId()

	err := n.ctx.DB.Notice().Create(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n noticeService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertNotice)
	err := n.ctx.DB.Notice().Update(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n noticeService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeQuery)
	err := n.ctx.DB.Notice().Delete(*r)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n noticeService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeQuery)
	data, err := n.ctx.DB.Notice().Get(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n noticeService) Search(req interface{}) (interface{}, interface{}) {
	r := req.(*models.NoticeQuery)
	data, err := n.ctx.DB.Notice().Search(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n noticeService) Check(req interface{}) (interface{}, interface{}) {

	// ToDo

	return nil, nil
}
