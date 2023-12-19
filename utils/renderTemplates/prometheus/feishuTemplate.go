package prometheus

import (
	"fmt"
	dto2 "prometheus-manager/controllers/dto"
	"prometheus-manager/globals"
	"prometheus-manager/utils/renderTemplates"
	"strconv"
	"strings"
	"time"
)

type FeiShu struct{}

// FeiShuMsgTemplate 飞书消息卡片模版
func (f *FeiShu) FeiShuMsgTemplate(prometheusAlert PrometheusAlert) (msg dto2.FeiShuMsg) {

	defaultTemplate := dto2.FeiShuMsg{
		MsgType: "interactive",
		Card: dto2.Cards{
			Config: dto2.Configs{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Header: dto2.Headers{
				Template: "",
				Title: dto2.Titles{
					Content: "",
					Tag:     "plain_text",
				},
			},
		},
	}

	switch prometheusAlert.alerts.Status {
	case "firing":
		return firingMsgTemplate(defaultTemplate, prometheusAlert.alerts, prometheusAlert.actionValues, prometheusAlert.currentDutyUser)
	case "resolved":
		return resolvedMsgTemplate(defaultTemplate, prometheusAlert.alerts)
	case "silence":
		return silenceMsgTemplate(defaultTemplate, prometheusAlert.alerts, prometheusAlert.actionValues, prometheusAlert.actionUser)
	}
	return

}

// firingMsgTemplate 告警模版
func firingMsgTemplate(template dto2.FeiShuMsg, v dto2.AlertInfo, ActionsValueStr dto2.CreateAlertSilence, dutyUser string) dto2.FeiShuMsg {

	var (
		confirmPrompt = "静默 " + strconv.Itoa(int(globals.Config.AlertManager.SilenceTime)) + " 分钟"
	)

	urlLine := strings.Split(v.GeneratorURL, "/")
	v.GeneratorURL = globals.Config.Prometheus.URL + "/" + urlLine[len(urlLine)-1]

	elements := []dto2.Elements{
		{
			Tag:            "column_set",
			FlexMode:       "none",
			BackgroupStyle: "default",
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Text: dto2.Texts{
				Content: " ",
				Tag:     "plain_text",
			},
		},
		{
			Tag: "hr",
		},
		{
			Tag: "div",
			Text: dto2.Texts{
				Content: dutyUser,
				Tag:     "lark_md",
			},
		},
		{
			Actions: []dto2.Actions{
				{
					Tag: "button",
					Text: dto2.ActionsText{
						Content: "🔕 告警静默",
						Tag:     "plain_text",
					},
					Type:  "primary",
					Value: ActionsValueStr,
					Confirm: dto2.Confirms{
						Title: dto2.Titles{
							Content: "确认",
							Tag:     "plain_text",
						},
						Text: dto2.Texts{
							Content: confirmPrompt,
							Tag:     "plain_text",
						},
					},
					MultiURL: nil,
				},
				{
					Tag: "button",
					Text: dto2.ActionsText{
						Content: "⛓️ 告警链接",
						Tag:     "plain_text",
					},
					Type: "primary",
					MultiURL: &dto2.MultiURLs{
						URL: v.GeneratorURL,
					},
					Confirm: dto2.Confirms{
						Title: dto2.Titles{
							Content: "确认",
							Tag:     "plain_text",
						},
						Text: dto2.Texts{
							Content: fmt.Sprintf("查询当前 ID: %s 的告警信息", v.Fingerprint),
							Tag:     "plain_text",
						},
					},
				},
				//{
				//	Tag: "button",
				//	Text: dto.ActionsText{
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
			Elements: []dto2.ElementsElements{
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
func resolvedMsgTemplate(template dto2.FeiShuMsg, v dto2.AlertInfo) dto2.FeiShuMsg {

	elements := []dto2.Elements{
		{
			Tag:            "column_set",
			FlexMode:       "none",
			BackgroupStyle: "default",
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Text: dto2.Texts{
				Content: " ",
				Tag:     "plain_text",
			},
		},
		{
			Tag: "hr",
		},
		{
			Tag: "note",
			Elements: []dto2.ElementsElements{
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
func silenceMsgTemplate(template dto2.FeiShuMsg, v dto2.AlertInfo, ActionsValueStr dto2.CreateAlertSilence, actionUserID string) dto2.FeiShuMsg {

	endsT, _ := time.Parse(time.RFC3339, ActionsValueStr.EndsAt)
	endsT = endsT.Add(8 * time.Hour)
	info := renderTemplates.GetFeiShuUserInfo(actionUserID)
	silenceMsgContent := fmt.Sprintf("操作人: %s\n静默时长: %v 分钟\n结束时间: %s\n", info.Data.User.Name, globals.Config.AlertManager.SilenceTime, endsT.Format(globals.Layout))

	elements := []dto2.Elements{
		{
			Tag:            "column_set",
			FlexMode:       "none",
			BackgroupStyle: "default",
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Columns: []dto2.Columns{
				{
					Tag:           "column",
					Width:         "weighted",
					Weight:        1,
					VerticalAlign: "top",
					Elements: []dto2.ColumnsElements{
						{
							Tag: "div",
							Text: dto2.Texts{
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
			Text: dto2.Texts{
				Content: " ",
				Tag:     "plain_text",
			},
		},
		{
			Tag: "hr",
		},
		{
			Tag: "div",
			Text: dto2.Texts{
				Content: silenceMsgContent,
				Tag:     "plain_text",
			},
		},
		{
			Tag: "hr",
		},
		{
			Tag: "note",
			Elements: []dto2.ElementsElements{
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
