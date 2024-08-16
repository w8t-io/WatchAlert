package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestNewElasticSearchClient(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), "http://192.168.1.190:9200", "", "")
	if err != nil {
		logrus.Errorf("client -> %s", err.Error())
		return
	}

	q := []ESQueryFilter{
		{
			"message",
			"docker",
		},
		{
			"message",
			"ready",
		},
	}
	client.Query(context.Background(), "test-2024-05.20", q, 1000000)
}
