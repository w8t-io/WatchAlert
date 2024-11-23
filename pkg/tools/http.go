package tools

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"
	"time"
)

func Get(headers map[string]string, url string, timeout int) (*http.Response, error) {
	// 统一跳过证书检测，避免存在不安全的https
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools get 请求建立失败, err: %s", err.Error()))
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools get 请求发送失败, err: %s", err.Error()))
		return nil, err
	}

	return resp, nil
}

func Post(headers map[string]string, url string, bodyReader *bytes.Reader, timeout int) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}

	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools post 请求建立失败, err: %s", err.Error()))
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Tools post 请求发送失败, err: %s", err.Error()))
		return nil, err
	}

	return resp, nil
}
