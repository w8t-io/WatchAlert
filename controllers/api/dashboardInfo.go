package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
)

type DashboardInfoController struct {
	models.AlertCurEvent
}

type ResponseDashboardInfo struct {
	CountAlertRules   int64                    `json:"countAlertRules"`
	CurAlerts         int                      `json:"curAlerts"`
	CurAlertList      []string                 `json:"curAlertList"`
	AlarmDistribution AlarmDistribution        `json:"alarmDistribution"`
	ServiceResource   []models.ServiceResource `json:"serviceResource"`
}

type AlarmDistribution struct {
	P0 int `json:"P0"`
	P1 int `json:"P1"`
	P2 int `json:"P2"`
}

func (di DashboardInfoController) GetDashboardInfo(ctx *gin.Context) {

	var (
		// 规则总数
		countAlertRules int64
		// 当前告警
		keys []string
	)
	// 告警分布
	alarmDistribution := make(map[string]int)
	globals.DBCli.Model(&models.AlertRule{}).Count(&countAlertRules)

	cursor := uint64(0)
	pattern := models.FiringAlertCachePrefix + "*"
	// 每次获取的键数量
	count := int64(100)

	for {
		var curKeys []string
		var err error

		curKeys, cursor, err = globals.RedisCli.Scan(cursor, pattern, count).Result()
		if err != nil {
			break
		}

		keys = append(keys, curKeys...)

		if cursor == 0 {
			break
		}
	}

	var curAlertList []string
	for _, v := range keys {
		if len(curAlertList) >= 5 {
			continue
		}
		alarmDistribution[di.GetCache(v).Severity] += 1
		curAlertList = append(curAlertList, di.GetCache(v).Annotations)
	}

	var resource []models.ServiceResource
	globals.DBCli.Model(&models.ServiceResource{}).Find(&resource)

	response.Success(ctx, ResponseDashboardInfo{
		CountAlertRules: countAlertRules,
		CurAlerts:       len(keys),
		CurAlertList:    curAlertList,
		AlarmDistribution: AlarmDistribution{
			P0: alarmDistribution["P0"],
			P1: alarmDistribution["P1"],
			P2: alarmDistribution["P2"],
		},
		ServiceResource: resource,
	}, "success")

}
