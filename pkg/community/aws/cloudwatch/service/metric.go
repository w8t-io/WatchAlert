package service

import types2 "watchAlert/pkg/community/aws/cloudwatch/types"

type (
	awsCloudWatchService struct{}

	InterAwsCloudWatchService interface {
		GetMetricTypes() (interface{}, interface{})
		GetMetricNames(req interface{}) (interface{}, interface{})
		GetStatistics() (interface{}, interface{})
		GetDimensions(req interface{}) (interface{}, interface{})
	}
)

func NewInterAwsCloudWatchService() InterAwsCloudWatchService {
	return awsCloudWatchService{}
}

func (a awsCloudWatchService) GetMetricTypes() (interface{}, interface{}) {
	var mt []string
	for k, _ := range types2.NamespaceMetricsMap {
		mt = append(mt, k)
	}

	return mt, nil
}

func (a awsCloudWatchService) GetMetricNames(req interface{}) (interface{}, interface{}) {
	r := req.(*types2.MetricNamesQuery)
	return types2.NamespaceMetricsMap[r.MetricType], nil
}

func (a awsCloudWatchService) GetStatistics() (interface{}, interface{}) {
	return []string{
		"Average",
		"Maximum",
		"Minimum",
		"Sum",
		"SampleCount",
		"IQM",
	}, nil
}

func (a awsCloudWatchService) GetDimensions(req interface{}) (interface{}, interface{}) {
	r := req.(*types2.RdsDimensionReq)
	return types2.NamespaceDimensionKeysMap[r.MetricType], nil
}
