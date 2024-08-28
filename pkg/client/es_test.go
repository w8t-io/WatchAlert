package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
<<<<<<< HEAD
)

func TestNewElasticSearchClient(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), "http://192.168.1.190:9200", "", "")
=======
	"watchAlert/internal/models"
	"watchAlert/internal/types"
)

func TestNewElasticSearchClient(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), models.AlertDataSource{})
>>>>>>> Cairry-master
	if err != nil {
		logrus.Errorf("client -> %s", err.Error())
		return
	}

<<<<<<< HEAD
	q := []ESQueryFilter{
=======
	q := []types.ESQueryFilter{
>>>>>>> Cairry-master
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
