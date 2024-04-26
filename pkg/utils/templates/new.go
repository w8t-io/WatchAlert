package templates

import (
	"watchAlert/internal/models"
)

type Template struct {
	CardContentMsg string
}

func NewTemplate(alert models.AlertCurEvent, notice models.AlertNotice) Template {
	switch notice.NoticeType {
	case "FeiShu":
		return Template{CardContentMsg: feishuTemplate(alert, notice)}
	case "DingDing":
		return Template{CardContentMsg: dingdingTemplate(alert, notice)}
	}

	return Template{}
}
