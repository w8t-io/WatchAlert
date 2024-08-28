package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

func TestNewElasticSearchClient(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), models.AlertDataSource{})
	if err != nil {
		logrus.Errorf("client -> %s", err.Error())
		return
	}

	q := []types.ESQueryFilter{
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
