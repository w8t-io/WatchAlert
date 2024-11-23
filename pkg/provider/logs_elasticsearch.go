package provider

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"watchAlert/internal/models"
	utilsHttp "watchAlert/pkg/tools"
)

type ElasticSearchDsProvider struct {
	cli *elastic.Client
	url string
}

func NewElasticSearchClient(ctx context.Context, ds models.AlertDataSource) (LogsFactoryProvider, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(ds.ElasticSearch.Url),
		elastic.SetBasicAuth(ds.ElasticSearch.Username, ds.ElasticSearch.Password),
		elastic.SetSniff(false),
	)
	if err != nil {
		return ElasticSearchDsProvider{}, err
	}

	return ElasticSearchDsProvider{
		client,
		ds.ElasticSearch.Url,
	}, nil
}

type esQueryResponse struct {
	Source map[string]interface{} `json:"_source"`
}

func (e ElasticSearchDsProvider) Query(options LogQueryOptions) ([]Logs, int, error) {
	filter := elastic.NewBoolQuery()
	for _, f := range options.ElasticSearch.QueryFilter {
		if f.Field != "" && f.Value != "" {
			filter.Must(elastic.NewMatchQuery(f.Field, f.Value))
		}
	}

	filter.Must(elastic.NewRangeQuery("@timestamp").Gte(options.StartAt.(string)).Lte(options.EndAt.(string)))

	res, err := e.cli.Search().
		Index(options.ElasticSearch.Index).
		Query(filter).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, 0, err
	}

	var response []esQueryResponse
	marshalHits, err := json.Marshal(res.Hits.Hits)
	if err != nil {
		return nil, 0, err
	}
	err = json.Unmarshal(marshalHits, &response)
	if err != nil {
		return nil, 0, err
	}

	var (
		data      []Logs
		msg       []interface{}
		kvMapList []map[string]interface{}
	)
	for _, v := range response {
		kvMapList = append(kvMapList, v.Source)
	}

	for _, m := range kvMapList {
		msg = append(msg, m["message"])
	}

	data = append(data, Logs{
		ProviderName: ElasticSearchDsProviderName,
		Metric:       commonKeyValuePairs(kvMapList),
		Message:      msg,
	})

	return data, len(msg), nil
}

func (e ElasticSearchDsProvider) Check() (bool, error) {
	res, err := utilsHttp.Get(nil, e.url+"/_cat/health", 10)
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, err
	}
	return true, nil
}
