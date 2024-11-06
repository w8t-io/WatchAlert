package services

import (
	"fmt"
	"watchAlert/alert"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/tools"
)

type monitorService struct {
	ctx *ctx.Context
}

type InterMonitorService interface {
	Create(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
}

func newInterMonitorService(ctx *ctx.Context) InterMonitorService {
	return &monitorService{
		ctx: ctx,
	}
}

func (m monitorService) Create(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MonitorSSLRule)
	r.ID = "m-" + tools.RandId()

	if *r.Enabled {
		alert.MonEvalTask.Submit(m.ctx, *r)
	}

	alert.MonConsumerTask.Add(*r)

	err := m.ctx.DB.MonitorSSL().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m monitorService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MonitorSSLRule)
	alert.MonConsumerTask.Stop(r.ID)
	alert.MonEvalTask.Stop(r.ID)
	if *r.Enabled {
		alert.MonEvalTask.Submit(m.ctx, *r)
		alert.MonConsumerTask.Add(*r)
	}

	err := m.ctx.DB.MonitorSSL().Update(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m monitorService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MonitorSSLRuleQuery)
	alert.MonEvalTask.Stop(r.ID)
	alert.MonConsumerTask.Stop(r.ID)
	key := fmt.Sprintf("%s:%s%s--", r.TenantId, models.FiringAlertCachePrefix, r.ID)
	m.ctx.Redis.Event().DelCache(key)

	err := m.ctx.DB.MonitorSSL().Delete(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m monitorService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MonitorSSLRuleQuery)
	data, err := m.ctx.DB.MonitorSSL().List(*r)
	if err != nil {
		return nil, err
	}

	//for k, v := range data {
	//	var object models.AlertCurEvent
	//	key := fmt.Sprintf("%s:%s%s--", r.TenantId, models.FiringAlertCachePrefix, v.ID)
	//	result, err := m.ctx.Redis.Redis().Get(key).Result()
	//	if err == nil {
	//		err = json.Unmarshal([]byte(result), &object)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	data[k].ResponseTime = object.ResponseTime
	//	data[k].TimeRemaining = object.TimeRemaining
	//}

	return data, nil
}

func (m monitorService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*models.MonitorSSLRuleQuery)
	data, err := m.ctx.DB.MonitorSSL().Get(*r)
	if err != nil {
		return nil, err
	}

	return data, nil
}
