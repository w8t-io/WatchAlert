package provider

import (
	"net"
	"time"
)

type Tcper struct{}

func NewEndpointTcper() EndpointFactoryProvider {
	return Tcper{}
}

func (p Tcper) Pilot(option EndpointOption) (EndpointValue, error) {
	startTime := time.Now()

	// 尝试拨测指定地址和端口
	conn, err := net.DialTimeout("tcp", option.Endpoint, time.Duration(option.Timeout)*time.Second)
	responseTime := time.Since(startTime)

	// 结果处理
	result := TcperInformation{
		Address:      option.Endpoint,
		ResponseTime: float64(responseTime),
		IsSuccessful: err == nil,
	}
	if err != nil {
		result.ErrorMessage = err.Error()
	} else {
		// 关闭连接
		conn.Close()
	}

	return convertTcperToEndpointValue(result), nil
}

func convertTcperToEndpointValue(detail TcperInformation) EndpointValue {
	return EndpointValue{
		"address":      detail.Address,
		"ResponseTime": detail.ResponseTime,
		"IsSuccessful": detail.IsSuccessful,
		"ErrorMessage": detail.ErrorMessage,
	}
}
