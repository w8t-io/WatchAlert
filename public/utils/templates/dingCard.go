package templates

import (
	"fmt"
	"watchAlert/models"
	"watchAlert/public/utils/cmd"
)

type DingDing struct{}

func (d DingDing) Template(alert models.AlertCurEvent, notice models.AlertNotice) string {

	Title := ParserTemplate("Title", alert, notice.Template)
	Footer := ParserTemplate("Footer", alert, notice.Template)

	userId := alert.DutyUser

	if alert.DutyUser != "暂无" {
		alert.DutyUser = fmt.Sprintf("@%s", alert.DutyUser)
	}

	t := models.DingMsg{
		Msgtype: "markdown",
		Markdown: models.Markdown{
			Title: Title,
			Text: "**" + Title + "**" +
				"\n" + "\n" +
				ParserTemplate("Event", alert, notice.Template) +
				"\n" +
				Footer,
		},
		At: models.At{
			AtUserIds: []string{userId},
			IsAtAll:   false,
		},
	}

	cardContentString := cmd.JsonMarshal(t)

	return cardContentString

}
