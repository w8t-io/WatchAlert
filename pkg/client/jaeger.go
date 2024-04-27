package client

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/http"
)

type JaegerClient struct {
	BaseURL string
}

type JaegerQueryOptions struct {
	Tags    string `json:"tags,omitempty"`    // 查询标签
	Service string `json:"service,omitempty"` // 服务名称
	Limit   int64  `json:"limit,omitempty"`   // 要返回的最大条目数
	StartAt int64  `json:"startAt,omitempty"` // 查询的开始时间，以微秒 Unix 表示。
	EndAt   int64  `json:"endAt,omitempty"`   // 查询的结束时间，以微秒 Unix 表示。
}

func NewJaegerClient(query models.AlertDataSource) JaegerClient {
	_, err := http.Get(query.HTTP.URL)
	if err != nil {
		global.Logger.Sugar().Errorf(err.Error())
		return JaegerClient{}
	}

	return JaegerClient{BaseURL: query.HTTP.URL}
}

type JaegerResult struct {
	Data []JaegerData `json:"data"`
}

type JaegerData struct {
	TraceID string `json:"traceID"`
}

func (jc JaegerClient) JaegerQuery(options JaegerQueryOptions) JaegerResult {

	curTime := time.Now()

	if options.Limit == 0 {
		options.Limit = 100
	}

	if options.StartAt == 0 {
		duration, _ := time.ParseDuration(strconv.Itoa(1) + "h")
		options.StartAt = curTime.Add(-duration).UnixNano()
	}

	if options.EndAt == 0 {
		options.EndAt = curTime.UnixNano()
	}

	args := fmt.Sprintf("/api/traces?service=%s&start=%d&end=%d&limit=%d&tags=%s", options.Service, options.StartAt, options.EndAt, options.Limit, options.Tags)
	requestURL := jc.BaseURL + args

	res, err := http.Get(requestURL)
	if err != nil {
		return JaegerResult{}
	}

	body, _ := io.ReadAll(res.Body)

	var jaegerResult JaegerResult
	err = json.Unmarshal(body, &jaegerResult)
	if err != nil {
		return JaegerResult{}
	}

	return jaegerResult

}

func (j JaegerData) GetFingerprint() string {
	return j.TraceID
}

func (j JaegerData) GetMetric(rule models.AlertRule) map[string]interface{} {
	return map[string]interface{}{
		"service": rule.JaegerConfig.Service,
	}
}

func (j JaegerData) GetAnnotations(rule models.AlertRule) string {
	return fmt.Sprintf("\n服务: %s 链路中存在异常状态码接口\nJaeger URL: %s/trace/%s\n\n详情查看 Jaeger Trace ⬆️", rule.JaegerConfig.Service, global.Config.Jaeger.URL, j.TraceID)
}

type JaegerServiceData struct {
	Data []string `json:"data"`
}

func (jc JaegerClient) GetJaegerService() (JaegerServiceData, error) {
	url := jc.BaseURL + "/api/services"
	res, err := http.Get(url)
	if err != nil {
		return JaegerServiceData{}, nil
	}

	if err != nil {
		return JaegerServiceData{}, err
	}

	if res.StatusCode != 200 {
		return JaegerServiceData{}, fmt.Errorf("后端服务请求异常, 上游返回 %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	var resData JaegerServiceData
	err = json.Unmarshal(body, &resData)
	if err != nil {
		return JaegerServiceData{}, err
	}
	return resData, nil
}
