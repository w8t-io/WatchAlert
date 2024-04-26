package templates

import (
	"fmt"
	models2 "watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

func dingdingTemplate(alert models2.AlertCurEvent, notice models2.AlertNotice) string {

	Title := ParserTemplate("Title", alert, notice.Template)
	Footer := ParserTemplate("Footer", alert, notice.Template)

	userId := alert.DutyUser

	if alert.DutyUser != "暂无" {
		alert.DutyUser = fmt.Sprintf("@%s", alert.DutyUser)
	}

	t := models2.DingMsg{
		Msgtype: "markdown",
		Markdown: models2.Markdown{
			Title: Title,
			Text: "**" + Title + "**" +
				"\n" + "\n" +
				ParserTemplate("Event", alert, notice.Template) +
				"\n" +
				Footer,
		},
		At: models2.At{
			AtUserIds: []string{userId},
			IsAtAll:   false,
		},
	}

	cardContentString := cmd.JsonMarshal(t)

	return cardContentString

}
