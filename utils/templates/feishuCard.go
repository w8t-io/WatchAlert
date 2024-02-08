package templates

import (
	"encoding/json"
	"strings"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type FeiShu struct{}

// Template 飞书消息卡片模版
func (f *FeiShu) Template(alert models.AlertCurEvent, notice models.AlertNotice) string {

	defaultTemplate := models.FeiShuMsg{
		MsgType: "interactive",
		Card: models.Cards{
			Config: models.Configs{
				WideScreenMode: true,
				EnableForward:  true,
			},
		},
	}

	var cardContentString string
	if notice.EnableCard == "true" {

		var tmplC models.Cards
		switch alert.IsRecovered {
		case false:
			_ = json.Unmarshal([]byte(notice.TemplateFiring), &tmplC)
		case true:
			_ = json.Unmarshal([]byte(notice.TemplateRecover), &tmplC)
		}
		defaultTemplate.Card.Elements = tmplC.Elements
		defaultTemplate.Card.Header = tmplC.Header
		cardContentString = cmd.JsonMarshal(defaultTemplate)
		cardContentString = ParserTemplate("", alert, cardContentString)

	} else {

		cardHeader := models.Headers{
			Template: ParserTemplate("TitleColor", alert, notice.Template),
			Title: models.Titles{
				Content: ParserTemplate("Title", alert, notice.Template),
				Tag:     "plain_text",
			},
		}
		cardElements := []models.Elements{
			{
				Tag:            "column_set",
				FlexMode:       "none",
				BackgroupStyle: "default",
				Columns: []models.Columns{
					{
						Tag:           "column",
						Width:         "weighted",
						Weight:        1,
						VerticalAlign: "top",
						Elements: []models.ColumnsElements{
							{
								Tag: "div",
								Text: models.Texts{
									Content: ParserTemplate("Event", alert, notice.Template),
									Tag:     "lark_md",
								},
							},
						},
					},
				},
			},
			{
				Tag: "hr",
			},
			{
				Tag: "note",
				Elements: []models.ElementsElements{
					{
						Tag:     "plain_text",
						Content: ParserTemplate("Footer", alert, notice.Template),
					},
				},
			},
		}

		defaultTemplate.Card.Elements = cardElements
		defaultTemplate.Card.Header = cardHeader
		cardContentString = cmd.JsonMarshal(defaultTemplate)

	}

	// 需要将所有换行符进行转义
	cardContentString = strings.Replace(cardContentString, "\n", "\\n", -1)

	return cardContentString

}
