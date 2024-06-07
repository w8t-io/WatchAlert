package service

import (
	"watchAlert/pkg/community/aws/cloudwatch/response"
	"watchAlert/pkg/community/aws/types"
)

type (
	awsRegionService struct{}

	InterAwsRegionService interface {
		GetRegion() ([]response.Regions, error)
	}
)

func NewInterAwsRegionService() InterAwsRegionService {
	return awsRegionService{}
}

func (a awsRegionService) GetRegion() ([]response.Regions, error) {
	var rs []response.Regions
	for _, r := range types.Regions {
		rs = append(rs, response.Regions{
			Label: &r,
			Value: &r,
		})
	}

	return rs, nil
}
