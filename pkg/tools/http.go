package tools

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
	"watchAlert/internal/global"
)

func Get(headers map[string]string, url string) (*http.Response, error) {
	// 统一跳过证书检测，避免存在不安全的https
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	if err != nil {
		global.Logger.Sugar().Error("请求建立失败: ", err)
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		global.Logger.Sugar().Error("请求发送失败: ", err)
		return nil, err
	}

	return resp, nil
}

func Post(headers map[string]string, url string, bodyReader *bytes.Reader) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	if err != nil {
		global.Logger.Sugar().Error("请求建立失败: ", err)
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		global.Logger.Sugar().Error("请求发送失败: ", err)
		return nil, err
	}

	return resp, nil
}
