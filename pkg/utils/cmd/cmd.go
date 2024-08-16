package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"watchAlert/internal/global"
)

func RandId() string {

	return xid.New().String()

}

func RandUid() string {

	limit := 8
	gid := xid.New().String()

	var xx []string
	for _, v := range gid {
		xx = append(xx, string(v))
	}

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(gid))

	var id string
	for i := 0; i < limit; i++ {
		id += xx[perm[i]]
	}

	return id

}

func RandUuid() string {

	return uuid.NewString()

}

func JsonMarshal(v interface{}) string {

	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)

}

// ParserVariables 处理告警内容中变量形式的字符串，替换为对应的值
func ParserVariables(annotations string, data map[string]interface{}) string {
	// 使用正则表达式匹配变量形式的字符串
	re := regexp.MustCompile(`\$\{(.*?)\}`)

	// 使用正则表达式替换所有匹配的变量
	return re.ReplaceAllStringFunc(annotations, func(match string) string {
		variable := match[2 : len(match)-1] // 提取变量名
		return fmt.Sprintf("%v", getJSONValue(data, variable))
	})
}

// 通过变量形式 ${key} 获取 JSON 数据中的值
func getJSONValue(data map[string]interface{}, variable string) interface{} {
	// 使用正则表达式分割键名数组
	keys := regexp.MustCompile(`\.`).Split(variable, -1)

	// 逐级获取 JSON 数据中的值
	for _, key := range keys {
		if v, ok := data[key]; ok {
			data, ok = v.(map[string]interface{})
			if !ok {
				return v
			}
		} else {
			return nil
		}
	}

	return nil
}

func IsJSON(str string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

func FormatJson(s string) string {
	var ns string
	if IsJSON(s) {
		// 将字符串解析为map类型
		var data map[string]interface{}
		err := json.Unmarshal([]byte(s), &data)
		if err != nil {
			global.Logger.Sugar().Errorf("Error parsing JSON: %s", err.Error())
		} else {
			// 格式化JSON并输出
			formattedJson, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				global.Logger.Sugar().Errorf("Error marshalling JSON: %s", err.Error())
			} else {
				ns = string(formattedJson)
			}
		}
	} else {
		ns = s
	}
	return ns
}

func FormatTimeToUTC(t int64) string {
	utcTime := time.Unix(t, 0).UTC()
	utcTimeString := utcTime.Format("2006-01-02T15:04:05.999Z")
	return utcTimeString
}

func ParserDuration(curTime time.Time, logScope int, timeType string) time.Time {
	duration, err := time.ParseDuration(strconv.Itoa(logScope) + timeType)
	if err != nil {
		logrus.Error(err.Error())
		return time.Time{}
	}
	startsAt := curTime.Add(-duration)
	return startsAt
}
