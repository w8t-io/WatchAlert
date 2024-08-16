package client

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/utils/cmd"
)

type ElasticSearchClient struct {
	cli *elastic.Client
}

func NewElasticSearchClient(ctx context.Context, ds models.AlertDataSource) (ElasticSearchClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(ds.ElasticSearch.Url),
		elastic.SetBasicAuth(ds.ElasticSearch.Username, ds.ElasticSearch.Password),
		elastic.SetSniff(false),
	)
	if err != nil {
		//logrus.Error(err.Error())
		global.Logger.Sugar().Errorf("ElasticSearch create client failed, %s", err.Error())
		return ElasticSearchClient{}, err
	}
	_, _, err = client.Ping(ds.ElasticSearch.Url).Do(ctx)
	if err != nil {
		//logrus.Error(err.Error())
		global.Logger.Sugar().Errorf("ElasticSearch ping test failed, %s", err.Error())
		return ElasticSearchClient{}, err
	}

	return ElasticSearchClient{
		client,
	}, nil
}

func (e ElasticSearchClient) Query(ctx context.Context, index string, query []types.ESQueryFilter, scope int64) ([]types.ESQueryResponse, error) {
	filter := elastic.NewBoolQuery()
	for _, f := range query {
		if f.Field != "" && f.Value != "" {
			filter.Must(elastic.NewMatchQuery(f.Field, f.Value))
		}
	}

	curTime := time.Now()
	from := cmd.ParserDuration(curTime, int(scope), "m")
	filter.Must(elastic.NewRangeQuery("@timestamp").Gte(cmd.FormatTimeToUTC(from.Unix())).Lte(cmd.FormatTimeToUTC(curTime.Unix())))

	res, err := e.cli.Search().
		Index(index).
		Query(filter).
		Pretty(true).
		Do(ctx)
	if err != nil {
		global.Logger.Sugar().Errorf("ElasticSearch query doc failed, %s", err.Error())
		return nil, err
	}

	var response []types.ESQueryResponse
	marshalHits, err := json.Marshal(res.Hits.Hits)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshalHits, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
