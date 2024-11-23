package provider

import (
	"fmt"
	"github.com/go-ping/ping"
	"time"
)

type Pinger struct{}

func NewEndpointPinger() EndpointFactoryProvider {
	return Pinger{}
}

func (p Pinger) Pilot(option EndpointOption) (EndpointValue, error) {
	var (
		detail PingerInformation
		ev     EndpointValue
	)
	pinger, err := ping.NewPinger(option.Endpoint)
	if err != nil {
		return ev, fmt.Errorf("New pinger error: %s", err.Error())
	}
	pinger.SetPrivileged(true)

	// 请求次数
	pinger.Count = option.ICMP.Count
	// 请求间隔
	pinger.Interval = time.Second * time.Duration(option.ICMP.Interval)
	// 超时时间
	pinger.Timeout = time.Second * time.Duration(option.Timeout)

	pinger.OnFinish = func(stats *ping.Statistics) {
		detail = PingerInformation{
			Address:     stats.Addr,
			PacketsSent: stats.PacketsSent,
			PacketsRecv: stats.PacketsRecv,
			PacketLoss:  stats.PacketLoss,
			Addr:        stats.Addr,
			IPAddr:      stats.IPAddr.String(),
			MinRtt:      float64(stats.MinRtt.Milliseconds()),
			MaxRtt:      float64(stats.MaxRtt.Milliseconds()),
			AvgRtt:      float64(stats.AvgRtt.Milliseconds()),
		}
	}

	err = pinger.Run()
	if err != nil {
		return ev, fmt.Errorf("Ping error: %s", err.Error())
	}

	return convertPingerToEndpointValues(detail), nil
}

func convertPingerToEndpointValues(detail PingerInformation) EndpointValue {
	return EndpointValue{
		"address":     detail.Address,
		"PacketsSent": detail.PacketsSent,
		"PacketsRecv": detail.PacketsRecv,
		"PacketLoss":  detail.PacketLoss,
		"Addr":        detail.Addr,
		"IPAddr":      detail.IPAddr,
		"MinRtt":      detail.MinRtt,
		"MaxRtt":      detail.MaxRtt,
		"AvgRtt":      detail.AvgRtt,
	}
}
