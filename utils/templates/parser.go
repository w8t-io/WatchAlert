package templates

import (
	"bytes"
	"text/template"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
)

var tmpl *template.Template

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
	if defineName == "" {
		err = tmpl.Execute(&buf, alert)
		return buf.String()
	}
	err = tmpl.ExecuteTemplate(&buf, defineName, alert)
	if err != nil {
		globals.Logger.Sugar().Error("告警模版执行失败 ->", err.Error())
		return ""
	}

	return buf.String()

}
