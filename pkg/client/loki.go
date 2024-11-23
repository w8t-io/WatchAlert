package client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
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

func NewLokiClient(query models.AlertDataSource) LokiClient {
	return LokiClient{BaseURL: query.HTTP.URL}
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

func (lc LokiClient) QueryRange(options QueryOptions) ([]Result, int, error) {

	curTime := time.Now()

	if options.Query == "" {
		return nil, 0, nil
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

	args := fmt.Sprintf("/loki/api/v1/query_range?query=%s&direction=%s&limit=%d&start=%s&end=%s", url.QueryEscape(options.Query), options.Direction, options.Limit, options.StartAt, options.EndAt)
	requestURL := lc.BaseURL + args
	res, err := tools.Get(nil, requestURL, 10)
	if err != nil {
		return nil, 0, err
	}

	var resultData result
	if err := tools.ParseReaderBody(res.Body, &resultData); err != nil {
		return nil, 0, err
	}

	// count 用于统计日志条数
	var count int
	for _, v := range resultData.Data.Result {
		count += len(v.Values)
	}

	return resultData.Data.Result, count, nil
}

func (r Result) GetFingerprint() string {
	// 使用 Loki 提供的 Stream label 进行 Hash 作为告警指纹.
	newMetric := map[string]interface{}{
		"namespace": r.Stream["namespace"],
		"container": r.Stream["container"],
	}
	h := md5.New()
	streamString := tools.JsonMarshal(newMetric)
	h.Write([]byte(streamString))
	fingerprint := hex.EncodeToString(h.Sum(nil))
	return fingerprint
}

func (r Result) GetMetric() map[string]interface{} {
	// 标签，用于推送告警消息时 获取相关 label 信息
	metricMap := make(map[string]interface{})
	for label, value := range r.Stream {
		metricMap[label] = value
	}

	delete(metricMap, "stream")
	delete(metricMap, "filename")
	return metricMap
}

func (r Result) GetAnnotations() interface{} {
	var logValue, annotations string
	if r.Values[0] != nil {
		if r.Values[0].([]interface{}) != nil {
			logValue = r.Values[0].([]interface{})[1].(string)
		}
	}
	annotations = tools.FormatJson(logValue)
	return annotations
}
