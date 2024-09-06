package provider

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
	"watchAlert/internal/models"
)

func TestNewElasticSearchClient(t *testing.T) {
	client, err := NewElasticSearchClient(context.Background(), models.AlertDataSource{})
	if err != nil {
		logrus.Errorf("client -> %s", err.Error())
		return
	}

	client.Query(LogQueryOptions{})
}
