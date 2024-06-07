package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"watchAlert/internal/repo"
	"watchAlert/pkg/client"
	"watchAlert/pkg/community/aws/cloudwatch/types"
)

type (
	awsRdsService struct{}

	InterAwsRdsService interface {
		GetDBInstanceIdentifier(req interface{}) (interface{}, interface{})
		GetDBClusterIdentifier(req interface{}) (interface{}, interface{})
	}
)

func NewInterAWSRdsService() InterAwsRdsService {
	return awsRdsService{}
}

func (a awsRdsService) GetDBInstanceIdentifier(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RdsInstanceReq)
	var dr repo.DatasourceRepo
	datasourceObj, err := dr.GetInstance(r.DatasourceId)
	if err != nil {
		return nil, err
	}

	cfg, err := client.NewAWSCredentialCfg(datasourceObj.AWSCloudWatch.AccessKey, datasourceObj.AWSCloudWatch.SecretKey)
	if err != nil {
		return nil, err
	}

	cli := cfg.RdsCli()
	input := &rds.DescribeDBInstancesInput{}
	result, err := cli.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, instance := range result.DBInstances {
		instances = append(instances, *instance.DBInstanceIdentifier)
	}

	return instances, nil
}

func (a awsRdsService) GetDBClusterIdentifier(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RdsClusterReq)
	var dr repo.DatasourceRepo
	datasourceObj, err := dr.GetInstance(r.DatasourceId)
	if err != nil {
		return nil, err
	}

	cfg, err := client.NewAWSCredentialCfg(datasourceObj.AWSCloudWatch.AccessKey, datasourceObj.AWSCloudWatch.SecretKey)
	if err != nil {
		return nil, err
	}

	cli := cfg.RdsCli()
	input := &rds.DescribeDBClustersInput{}
	result, err := cli.DescribeDBClusters(context.TODO(), input)

	var clusters []string
	for _, cluster := range result.DBClusters {
		clusters = append(clusters, *cluster.DBClusterIdentifier)
	}

	return clusters, nil
}
