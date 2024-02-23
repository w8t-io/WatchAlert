package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"watchAlert/globals"
)

const CachePrefix = "cur-alert-"

type AlertCurEvent struct {
	RuleId                 string                 `json:"rule_id"`
	RuleName               string                 `json:"rule_name"`
	DatasourceType         string                 `json:"datasource_type"`
	DatasourceId           string                 `json:"-" gorm:"datasource_id"`
	DatasourceIdList       []string               `json:"datasource_id" gorm:"-"`
	Fingerprint            string                 `json:"fingerprint"`
	Severity               int64                  `json:"severity"`
	PromQl                 string                 `json:"prom_ql"`
	Instance               string                 `json:"instance"`
	Metric                 string                 `json:"-" gorm:"metric"`
	MetricMap              map[string]interface{} `json:"metric" gorm:"-"`
	LabelsMap              map[string]string      `json:"labels" gorm:"-"`
	Labels                 string                 `json:"-" gorm:"labels"`
	EvalInterval           int64                  `json:"eval_interval"`
	ForDuration            int64                  `json:"for_duration"`
	NoticeId               string                 `json:"notice_id" gorm:"-"` // 默认通知对象ID
	NoticeGroupList        NoticeGroup            `json:"noticeGroup" gorm:"-"`
	NoticeGroup            string                 `json:"-" gorm:"noticeGroup"`
	Annotations            string                 `json:"annotations" gorm:"-"`
	IsRecovered            bool                   `json:"is_recovered" gorm:"-"`
	FirstTriggerTime       int64                  `json:"first_trigger_time"` // 第一次触发时间
	FirstTriggerTimeFormat string                 `json:"-" gorm:"-"`
	RepeatNoticeInterval   int64                  `json:"repeat_notice_interval"`  // 重复通知间隔时间
	LastEvalTime           int64                  `json:"last_eval_time" gorm:"-"` // 上一次评估时间
	LastSendTime           int64                  `json:"last_send_time" gorm:"-"` // 上一次发送时间
	RecoverTime            int64                  `json:"recover_time" gorm:"-"`   // 恢复时间
	RecoverTimeFormat      string                 `json:"-" gorm:"-"`
	DutyUser               string                 `json:"duty_user" gorm:"-"`
}

func (ace *AlertCurEvent) CurAlertCacheKey(ruleId, fingerprint string) string {

	// cur-alert-xxx-xxx
	return CachePrefix + ruleId + "-" + fingerprint

}

func (ace *AlertCurEvent) GetCache(key string) AlertCurEvent {

	var alert AlertCurEvent

	d, err := globals.RedisCli.Get(key).Result()
	_ = json.Unmarshal([]byte(d), &alert)
	if err != nil {
		return AlertCurEvent{}
	}

	return alert

}

func (ace *AlertCurEvent) SetCache(alert AlertCurEvent, expiration time.Duration) {

	var alertRule AlertRule
	// 设置缓存前检查当前 Rule 是否存在，避免出现删除/禁用规则后依旧能添加缓存。
	globals.DBCli.Where("rule_id = ? and enabled = ?", alert.RuleId, "true").Find(&alertRule)
	if alertRule.RuleId == alert.RuleId {
		alertJson, _ := json.Marshal(alert)
		globals.RedisCli.Set(ace.CurAlertCacheKey(alert.RuleId, alert.Fingerprint), string(alertJson), expiration)
	}

}

func (ace *AlertCurEvent) DelCache(key string) {

	//globals.RedisCli.Del(key)

	// 使用Scan命令获取所有匹配指定模式的键
	iter := globals.RedisCli.Scan(0, key, 0).Iterator()
	keysToDelete := make([]string, 0)

	// 遍历匹配的键
	for iter.Next() {
		key := iter.Val()
		keysToDelete = append(keysToDelete, key)
	}

	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}

	// 批量删除键
	if len(keysToDelete) > 0 {
		err := globals.RedisCli.Del(keysToDelete...).Err()
		if err != nil {
			log.Fatal(err)
		}
		globals.Logger.Sugar().Infof("移除告警消息 -> %s\n", keysToDelete)
	}

}

// ParserAnnotation 处理变量形式的字符串，替换为对应的值
func (ace *AlertCurEvent) ParserAnnotation(annotations string) string {

	// 查找变量形式的字符串，并替换为对应的值
	result := annotations
	for strings.Contains(result, "${") && strings.Contains(result, "}") {
		startIndex := strings.Index(result, "${")
		endIndex := strings.Index(result, "}")
		if startIndex == -1 || endIndex == -1 || endIndex <= startIndex {
			break
		}

		// 获取变量名称
		variable := result[startIndex+2 : endIndex]

		// 获取对应的值
		value := getJSONValue(ace.MetricMap, variable)

		// 替换变量形式的字符串为对应的值
		result = strings.Replace(result, "${"+variable+"}", fmt.Sprintf("%v", value), 1)
	}

	return result

}

// 通过变量形式 ${key} 获取 JSON 数据中的值
func getJSONValue(data map[string]interface{}, variable string) interface{} {
	// 将变量形式的字符串分割为键名数组
	keys := strings.Split(variable, ".")

	// 逐级获取 JSON 数据中的值
	value := data
	for _, key := range keys {
		if v, ok := value[key]; ok {
			if nextValue, ok := v.(map[string]interface{}); ok {
				value = nextValue
			} else {
				return v
			}
		} else {
			return nil
		}
	}

	return nil
}

func (ace *AlertCurEvent) GetFirstTime() int64 {

	ft := ace.GetCache(ace.CurAlertCacheKey(ace.RuleId, ace.Fingerprint)).FirstTriggerTime
	if ft == 0 {
		return time.Now().Unix()
	}
	return ft

}

func (ace *AlertCurEvent) GetLastEvalTime() int64 {

	curTime := time.Now().Unix()
	let := ace.GetCache(ace.CurAlertCacheKey(ace.RuleId, ace.Fingerprint)).LastEvalTime
	if let == 0 || let < curTime {
		return curTime
	}

	return let

}

func (ace *AlertCurEvent) GetLastSendTime() int64 {

	return ace.GetCache(ace.CurAlertCacheKey(ace.RuleId, ace.Fingerprint)).LastSendTime

}