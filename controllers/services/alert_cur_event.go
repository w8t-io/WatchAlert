package services

import (
	"encoding/json"
	"log"
	"strings"
	"watchAlert/globals"
	"watchAlert/models"
)

type AlertCurEventService struct{}

type InterAlertCurEventService interface {
	List(dsType string) ([]models.AlertCurEvent, error)
}

func NewInterAlertCurEventService() InterAlertCurEventService {
	return &AlertCurEventService{}
}

func (aces *AlertCurEventService) List(dsType string) ([]models.AlertCurEvent, error) {

	iter := globals.RedisCli.Scan(0, models.CachePrefix+"*", 0).Iterator()
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
		info, err := globals.RedisCli.Get(key).Result()
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

	if dsType != "" {
		var dsTypeDataList []models.AlertCurEvent
		for _, v := range dataList {
			if v.DatasourceType == dsType {
				dsTypeDataList = append(dsTypeDataList, v)
			}
		}
		dataList = dsTypeDataList
	}

	return dataList, nil

}
