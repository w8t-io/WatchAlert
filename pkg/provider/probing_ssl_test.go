package provider

import (
	"fmt"
	"testing"
	"watchAlert/pkg/tools"
)

func TestNewEndpointSSLer(t *testing.T) {
	buildOption := EndpointOption{
		Endpoint: "www.baidu.com",
		Timeout:  10,
	}
	pilot, err := NewEndpointSSLer().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(pilot))
}
