package prometheus

import (
	"fmt"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"prometheus-manager/utils"
	"strconv"
	"strings"
	"time"
)

type FeiShu struct{}

var (
	confirmPrompt = "静默 " + strconv.FormatInt(globals.Config.AlertManager.SilenceTime, 10) + " 分钟"
)

// FeiShuMsgTemplate 飞书消息卡片模版
func (f *FeiShu) FeiShuMsgTemplate(actionUser string, v models.AlertInfo, ActionsValueStr models.CreateAlertSilence) (msg models.FeiShuMsg) {

	defaultTemplate := models.FeiShuMsg{
		MsgType: "interactive",
		Card: models.Cards{
			Config: models.Configs{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Header: models.Headers{
				Template: "",
				Title: models.Titles{
					Content: "",
					Tag:     "plain_text",
				},
			},
		},
	}

	switch v.Status {
	case "firing":
		return firingMsgTemplate(defaultTemplate, v, ActionsValueStr)
	case "resolved":
		return resolvedMsgTemplate(defaultTemplate, v)
	case "silence":
		return silenceMsgTemplate(defaultTemplate, v, ActionsValueStr, actionUser)
	}
	return

}

// firingMsgTemplate 告警模版
func firingMsgTemplate(template models.FeiShuMsg, v models.AlertInfo, ActionsValueStr models.CreateAlertSilence) models.FeiShuMsg {

	var contentInfo string

	user := utils.GetCurrentDutyUser()
	if len(user) == 0 {
		contentInfo = "暂无安排值班人员"
	} else {
		contentInfo = fmt.Sprintf("**👤 值班人员：**<at id=%s></at>", user)
	}

	urlLine := strings.Split(v.GeneratorURL, "/")
	v.GeneratorURL = globals.Config.Prometheus.URL + "/" + urlLine[len(urlLine)-1]

	elements := []models.Elements{
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
								Content: "",
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**🫧 报警指纹：**\n" + v.Fingerprint,
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🤖 报警类型：**\n" + v.Labels["alertname"],
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**📌 报警等级：**\n" + v.Labels["severity"],
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🕘 开始时间：**\n" + v.StartsAt.Local().Format(globals.Layout),
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**🕟 结束时间：**\n" + v.EndsAt.Local().Format(globals.Layout),
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🖥 报警主机：**\n" + v.Labels["instance"],
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**📝 报警事件：**\n" + v.Annotations.Description,
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
		{
			Tag: "div",
			Text: models.Texts{
				Content: " ",
				Tag:     "plain_text",
			},
		},
		{
			Tag: "hr",
		},
		{
			Tag: "div",
			Text: models.Texts{
				Content: contentInfo,
				Tag:     "lark_md",
			},
		},
		{
			Actions: []models.Actions{
				{
					Tag: "button",
					Text: models.ActionsText{
						Content: "🔕 告警静默",
						Tag:     "plain_text",
					},
					Type:  "primary",
					Value: ActionsValueStr,
					Confirm: models.Confirms{
						Title: models.Titles{
							Content: "确认",
							Tag:     "plain_text",
						},
						Text: models.Texts{
							Content: confirmPrompt,
							Tag:     "plain_text",
						},
					},
				},
				{
					Tag: "button",
					Text: models.ActionsText{
						Content: "⛓️ 告警链接",
						Tag:     "plain_text",
					},
					Type: "primary",
					MultiURL: models.MultiURLs{
						URL: v.GeneratorURL,
					},
					Confirm: models.Confirms{
						Title: models.Titles{
							Content: "确认",
							Tag:     "plain_text",
						},
						Text: models.Texts{
							Content: fmt.Sprintf("查询当前 ID: %s 的告警信息", v.Fingerprint),
							Tag:     "plain_text",
						},
					},
				},
				//{
				//	Tag: "button",
				//	Text: models.ActionsText{
				//		Content: "👤 告警认领",
				//		Tag:     "plain_text",
				//	},
				//	Type:  "primary",
				//	Value: ActionsValueStr,
				//},
			},
			Tag: "action",
		},
		{
			Tag: "hr",
		},
		{
			Tag: "note",
			Elements: []models.ElementsElements{
				{
					Tag:     "plain_text",
					Content: "🧑‍💻 即时设计 - 运维团队",
				},
			},
		},
	}

	template.Card.Header.Template = "red"
	template.Card.Header.Title.Content = "【报警中】一级报警 - 即时设计 🔥"
	template.Card.Elements = elements

	return template

}

// resolvedMsgTemplate 恢复模版
func resolvedMsgTemplate(template models.FeiShuMsg, v models.AlertInfo) models.FeiShuMsg {

	elements := []models.Elements{
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
								Content: "",
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**🫧 报警指纹：**\n" + v.Fingerprint,
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🤖 报警类型：**\n" + v.Labels["alertname"],
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**📌 报警等级：**\n" + v.Labels["severity"],
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🕘 开始时间：**\n" + v.StartsAt.Local().Format(globals.Layout),
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**🕟 结束时间：**\n" + v.EndsAt.Local().Format(globals.Layout),
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🖥 报警主机：**\n" + v.Labels["instance"],
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**📝 报警事件：**\n" + v.Annotations.Description,
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
		{
			Tag: "div",
			Text: models.Texts{
				Content: " ",
				Tag:     "plain_text",
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
					Content: "🧑‍💻 即时设计 - 运维团队",
				},
			},
		},
	}

	template.Card.Header.Template = "green"
	template.Card.Header.Title.Content = "【已处理】一级报警 - 即时设计 ✨"
	template.Card.Elements = elements

	return template

}

// silenceMsgTemplate 静默模版
func silenceMsgTemplate(template models.FeiShuMsg, v models.AlertInfo, ActionsValueStr models.CreateAlertSilence, actionUserID string) models.FeiShuMsg {

	endsT, _ := time.Parse(time.RFC3339, ActionsValueStr.EndsAt)
	endsT = endsT.Add(8 * time.Hour)
	info := utils.GetFeiShuUserInfo(actionUserID)
	silenceMsgContent := fmt.Sprintf("操作人: %s\n静默时长: %v 分钟\n结束时间: %s\n", info.Data.User.Name, globals.Config.AlertManager.SilenceTime, endsT.Format(globals.Layout))

	elements := []models.Elements{
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
								Content: "",
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**🫧 报警指纹：**\n" + v.Fingerprint,
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🤖 报警类型：**\n" + v.Labels["alertname"],
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**📌 报警等级：**\n" + v.Labels["severity"],
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🕘 开始时间：**\n" + v.StartsAt.Local().Format(globals.Layout),
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**🕟 结束时间：**\n" + v.EndsAt.Local().Format(globals.Layout),
								Tag:     "lark_md",
							},
						},
					},
				},
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
								Content: "**🖥 报警主机：**\n" + v.Labels["instance"],
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
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
								Content: "**📝 报警事件：**\n" + v.Annotations.Description,
								Tag:     "lark_md",
							},
						},
					},
				},
			},
		},
		{
			Tag: "div",
			Text: models.Texts{
				Content: " ",
				Tag:     "plain_text",
			},
		},
		{
			Tag: "hr",
		},
		{
			Tag: "div",
			Text: models.Texts{
				Content: silenceMsgContent,
				Tag:     "plain_text",
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
					Content: "🧑‍💻 即时设计 - 运维团队",
				},
			},
		},
	}

	template.Card.Header.Template = "yellow"
	template.Card.Header.Title.Content = "【静默中】一级报警 - 即时设计 🧘"
	template.Card.Elements = elements

	return template

}
