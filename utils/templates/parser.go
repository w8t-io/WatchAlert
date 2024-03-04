package templates

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

var tmpl *template.Template

// ParserTemplate 处理告警推送的消息模版
func ParserTemplate(defineName string, alert models.AlertCurEvent, templateStr string) string {

	firstTriggerTime := time.Unix(alert.FirstTriggerTime, 0).Format(globals.Layout)
	recoverTime := time.Unix(alert.RecoverTime, 0).Format(globals.Layout)
	alert.FirstTriggerTimeFormat = firstTriggerTime
	alert.RecoverTimeFormat = recoverTime

	tmpl = template.Must(template.New("tmpl").Parse(templateStr))

	var (
		buf bytes.Buffer
		err error
	)

	if defineName == "Card" {
		err = tmpl.Execute(&buf, alert)
		data := parserEvent(alert)
		return cmd.ParserVariables(buf.String(), data)
	}

	err = tmpl.ExecuteTemplate(&buf, defineName, alert)
	if err != nil {
		globals.Logger.Sugar().Error("告警模版执行失败 ->", err.Error())
		return ""
	}

	// 前面只会渲染出模版框架, 下面来渲染告警数据内容
	if defineName == "Event" {
		data := parserEvent(alert)
		return cmd.ParserVariables(buf.String(), data)
	}

	return buf.String()

}

func parserEvent(alert models.AlertCurEvent) map[string]interface{} {

	data := make(map[string]interface{})

	if alert.DatasourceType == "AliCloudSLS" {
		eventJson := cmd.JsonMarshal(alert)
		eventJson = strings.ReplaceAll(eventJson, "\"{", "{")
		eventJson = strings.ReplaceAll(eventJson, "\\\\\"", "\"")
		eventJson = strings.ReplaceAll(eventJson, "\\\"", "\"")
		eventJson = strings.ReplaceAll(eventJson, "}\"", "}")
		eventJson = strings.ReplaceAll(eventJson, "}\\n\"", "}")
		eventJson = strings.ReplaceAll(eventJson, "\\{", "{")
		eventJson = strings.ReplaceAll(eventJson, "\\", "")
		eventJson = strings.ReplaceAll(eventJson, "\\\\\\\\", "")
		_ = json.Unmarshal([]byte(eventJson), &data)

		annotations, _ := data["annotations"].(map[string]interface{})
		// 将content进行转义, 在 ${annotations.content} 获取日志信息时用到.
		contentString := strconv.Quote(cmd.JsonMarshal(annotations["content"]))
		annotations["content"] = contentString
	} else {
		eventJson := cmd.JsonMarshal(alert)
		_ = json.Unmarshal([]byte(eventJson), &data)
	}

	return data

}
