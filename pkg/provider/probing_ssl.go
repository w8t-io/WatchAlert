package provider

import (
	"fmt"
	"time"
	"watchAlert/pkg/tools"
)

type Ssler struct{}

func NewEndpointSSLer() EndpointFactoryProvider {
	return Ssler{}
}

func (p Ssler) Pilot(option EndpointOption) (EndpointValue, error) {
	var (
		detail SslInformation
		ev     EndpointValue
	)
	startTime := time.Now()
	// 发起 HTTPS 请求
	resp, err := tools.Get(nil, "https://"+option.Endpoint, option.Timeout)
	if err != nil {
		return ev, err
	}
	defer resp.Body.Close()

	// 证书为空, 跳过检测
	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return ev, fmt.Errorf("证书为空, 跳过检测")
	}

	// 获取证书信息
	cert := resp.TLS.PeerCertificates[0]
	notBefore := cert.NotBefore // 证书开始时间
	notAfter := cert.NotAfter   // 证书过期时间
	currentTime := time.Now()

	// 计算剩余有效期（单位：天）
	timeRemaining := int64(notAfter.Sub(currentTime).Hours() / 24)

	detail = SslInformation{
		Address:       option.Endpoint,
		StartTime:     notBefore.Format("2006-01-02"), // 格式化开始时间
		ExpireTime:    notAfter.Format("2006-01-02"),  // 格式化过期时间
		TimeRemaining: float64(timeRemaining),
		ResponseTime:  float64(time.Since(startTime).Milliseconds()),
	}

	return convertSslerToEndpointValues(detail), nil
}

func convertSslerToEndpointValues(detail SslInformation) EndpointValue {
	return EndpointValue{
		"address":       detail.Address,
		"StartTime":     detail.StartTime,
		"ExpireTime":    detail.ExpireTime,
		"TimeRemaining": detail.TimeRemaining,
		"ResponseTime":  detail.ResponseTime,
	}
}
