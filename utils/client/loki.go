package client

import (
	"encoding/json"
	"fmt"
	"io"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/http"

	"strconv"
	"strings"
	"time"
)

type LokiClient struct {
	BaseURL string
}

type QueryOptions struct {
	Query     string `json:"query,omitempty"`     // 查询语句
	Direction string `json:"direction,omitempty"` // 日志排序顺序，支持的值为forward或backward，默认为backward
	Limit     int64  `json:"limit,omitempty"`     // 要返回的最大条目数
	StartAt   string `json:"startAt,omitempty"`   // 查询的开始时间，以纳秒 Unix 纪元表示。默认为一小时前
	EndAt     string `json:"endAt,omitempty"`     // 查询的结束时间，以纳秒 Unix 纪元表示。默认为现在
}

func NewLokiClient(datasourceId string) LokiClient {

	var datasource models.AlertDataSource
	globals.DBCli.Model(&models.AlertDataSource{}).Where("id = ?", datasourceId).First(&datasource)

	return LokiClient{BaseURL: datasource.HTTP.URL}

}

type result struct {
	Data Data `json:"data"`
}

type Data struct {
	ResultType string   `json:"status"`
	Result     []Result `json:"result"`
}

type Result struct {
	Stream map[string]string `json:"stream"`
	Values []interface{}     `json:"values"`
}

func (lc LokiClient) QueryRange(options QueryOptions) ([]Result, error) {

	curTime := time.Now()

	if options.Query == "" {
		return nil, nil
	}

	if options.Direction == "" {
		options.Direction = "backward"
	}

	if options.Limit == 0 {
		options.Limit = 100
	}

	if options.StartAt == "" {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		options.StartAt = curTime.Add(-duration).Format(time.RFC3339Nano)
	}

	if options.EndAt == "" {
		options.EndAt = curTime.Format(time.RFC3339Nano)
	}

	args := fmt.Sprintf("/loki/api/v1/query_range?query=%s&direction=%s&limit=%d&start=%s&end=%s", options.Query, options.Direction, options.Limit, options.StartAt, options.EndAt)
	requestURL := lc.BaseURL + args
	requestURL = strings.ReplaceAll(requestURL, "{", "%7B")
	requestURL = strings.ReplaceAll(requestURL, "}", "%7D")
	requestURL = strings.ReplaceAll(requestURL, `"`, "%22")
	requestURL = strings.ReplaceAll(requestURL, " ", "%20")

	res, err := http.Get(requestURL)
	if err != nil {
		return nil, nil
	}

	body, _ := io.ReadAll(res.Body)
	var resultData result
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		return nil, err
	}

	return resultData.Data.Result, nil

}
