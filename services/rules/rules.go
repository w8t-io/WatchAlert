package rules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RuleStr struct{}

func (r *RuleStr) List() ([]AlertingRule, error) {

	req, err := http.NewRequest(http.MethodGet, "http://172.17.84.238:9090/api/v1/rules", nil)
	if err != nil {
		log.Println("err 1")
	}

	body, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("err 2")
	}

	content, err := ioutil.ReadAll(body.Body)

	var res info
	err = json.Unmarshal(content, &res)

	var alerts []interface{}

	for i := 0; i < len(res.Data.Groups); i++ {
		for _, rule := range res.Data.Groups[i].Rules {

			var alert = make(map[string]interface{})

			for key, value := range rule.(map[string]interface{}) {

				alert[key] = value

			}

			alerts = append(alerts, alert)

		}
	}

	var alertRule []AlertingRule

	arr, err := json.Marshal(alerts)
	if err != nil {
		fmt.Println(err)
		return []AlertingRule{}, err
	}

	err = json.Unmarshal(arr, &alertRule)
	if err != nil {
		fmt.Println(err)
		return []AlertingRule{}, err
	}

	return alertRule, err

}

//func (r *Rule) Create() {
//
//	test := RulesResult{
//
//		Groups: []RuleGroup{
//			{
//				Name: "x",
//				Rules: Rules{
//					AlertingRule{
//						Name:     "",
//						Query:    "",
//						Duration: 60,
//						Labels: map[model.LabelName]model.LabelValue{
//							"xx": "xx",
//						},
//						Annotations: map[model.LabelName]model.LabelValue{
//							"summary":     "xx",
//							"description": "xx",
//						},
//					},
//				},
//			},
//		},
//	}
//
//	jsonBody := []byte()
//	bodyReader := bytes.NewReader(jsonBody)
//
//	req, err := http.NewRequest(http.MethodPost, "http://172.17.84.238:9090/api/v1/rules", bodyReader)
//	req.Header.Set("Content-Type", "application/json")
//	if err != nil {
//		log.Println("请求发送失败", err)
//	}
//
//}
