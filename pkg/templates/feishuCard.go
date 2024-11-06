package templates

import (
	"encoding/json"
	"strings"
	models "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

// Template 飞书消息卡片模版
func feishuTemplate(alert models.AlertCurEvent, noticeTmpl models.NoticeTemplateExample) string {

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
	if *noticeTmpl.EnableFeiShuJsonCard {
		var tmplC models.Cards
		switch alert.IsRecovered {
		case false:
			_ = json.Unmarshal([]byte(noticeTmpl.TemplateFiring), &tmplC)
		case true:
			_ = json.Unmarshal([]byte(noticeTmpl.TemplateRecover), &tmplC)
		}
		defaultTemplate.Card.Elements = tmplC.Elements
		defaultTemplate.Card.Header = tmplC.Header
		cardContentString = tools.JsonMarshal(defaultTemplate)
		cardContentString = ParserTemplate("Card", alert, cardContentString)

	} else {
		cardHeader := models.Headers{
			Template: ParserTemplate("TitleColor", alert, noticeTmpl.Template),
			Title: models.Titles{
				Content: ParserTemplate("Title", alert, noticeTmpl.Template),
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
									Content: ParserTemplate("Event", alert, noticeTmpl.Template),
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
						Content: ParserTemplate("Footer", alert, noticeTmpl.Template),
					},
				},
			},
		}

		defaultTemplate.Card.Elements = cardElements
		defaultTemplate.Card.Header = cardHeader
		cardContentString = tools.JsonMarshal(defaultTemplate)

	}

	// 需要将所有换行符进行转义
	cardContentString = strings.Replace(cardContentString, "\n", "\\n", -1)

	return cardContentString

}
