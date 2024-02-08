package models

type PromQueryResponse struct {
	Data struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Value  []string          `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
