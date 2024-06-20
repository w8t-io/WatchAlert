package services

import (
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type datasourceService struct {
	ctx *ctx.Context
}

type InterDatasourceService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	Search(req interface{}) (interface{}, interface{})
}

func newInterDatasourceService(ctx *ctx.Context) InterDatasourceService {
	return &datasourceService{
		ctx: ctx,
	}
}

func (ds datasourceService) Create(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.AlertDataSource)
	ok, err := dataSource.CheckHealth()
	if !ok {
		return nil, err
	}

	err = ds.ctx.DB.Datasource().Create(*dataSource)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) Update(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.AlertDataSource)
	err := ds.ctx.DB.Datasource().Update(*dataSource)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) Delete(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.DatasourceQuery)
	err := ds.ctx.DB.Datasource().Delete(*dataSource)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ds datasourceService) List(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.DatasourceQuery)
	data, err := ds.ctx.DB.Datasource().List(*dataSource)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ds datasourceService) Get(req interface{}) (interface{}, interface{}) {
	dataSource := req.(*models.DatasourceQuery)
	data, err := ds.ctx.DB.Datasource().Get(*dataSource)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ds datasourceService) Search(req interface{}) (interface{}, interface{}) {
	var newData []models.AlertDataSource
	r := req.(*models.DatasourceQuery)
	data, err := ds.ctx.DB.Datasource().Search(*r)
	if err != nil {
		return nil, err
	}
	newData = data

	return newData, nil
}
