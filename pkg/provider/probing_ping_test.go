package provider

import (
	"fmt"
	"testing"
	"watchAlert/pkg/tools"
)

func TestPinger(t *testing.T) {
	buildOption := EndpointOption{
		Endpoint: "8.147.234.89",
		Timeout:  10,
		ICMP: Eicmp{
			Interval: 1,
			Count:    5,
		},
	}

	pinger, err := NewEndpointPinger().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(pinger))
}
