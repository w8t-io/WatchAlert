package services

import (
	"encoding/json"
	"log"
	"strings"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type eventService struct {
	ctx *ctx.Context
}

type InterEventService interface {
	ListCurrentEvent(req interface{}) (interface{}, interface{})
	ListHistoryEvent(req interface{}) (interface{}, interface{})
}

func newInterEventService(ctx *ctx.Context) InterEventService {
	return &eventService{
		ctx: ctx,
	}
}

func (e eventService) ListCurrentEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertCurEventQuery)

	iter := e.ctx.Redis.Redis().Scan(0, r.TenantId+":"+models.FiringAlertCachePrefix+"*", 0).Iterator()
	keys := make([]string, 0)

	// 遍历匹配的键
	for iter.Next() {
		key := iter.Val()
		keys = append(keys, key)
	}

	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}

	var dataList []models.AlertCurEvent
	for _, key := range keys {
		var data models.AlertCurEvent
		info, err := e.ctx.Redis.Redis().Get(key).Result()
		if err != nil {
			return nil, err
		}

		newInfo := info
		newInfo = strings.Replace(newInfo, "\"[\\", "[", 1)
		newInfo = strings.Replace(newInfo, "\\\"]\"", "\"]", 1)
		err = json.Unmarshal([]byte(newInfo), &data)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}

	if r.DatasourceType != "" {
		var dsTypeDataList []models.AlertCurEvent
		for _, v := range dataList {
			if v.DatasourceType == r.DatasourceType {
				dsTypeDataList = append(dsTypeDataList, v)
			}
		}
		dataList = dsTypeDataList
	}

	return dataList, nil

}

func (e eventService) ListHistoryEvent(req interface{}) (interface{}, interface{}) {
	r := req.(*models.AlertHisEventQuery)
	data, err := e.ctx.DB.Event().GetHistoryEvent(*r)
	if err != nil {
		return nil, err
	}

	return data, err

}
