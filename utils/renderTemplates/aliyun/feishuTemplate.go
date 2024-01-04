package aliyun

import (
	"fmt"
	dto2 "watchAlert/controllers/dto"
	"watchAlert/globals"
	"strconv"
	"time"
)

type ALiYun struct{}

// FeiShuMsgTemplate 飞书消息卡片模版
func (a *ALiYun) FeiShuMsgTemplate(aliAlert renderALiYun) (msg dto2.FeiShuMsg) {

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

	switch aliAlert.AliAlert.Status {
	case "firing":
		return firingMsgTemplate(defaultTemplate, aliAlert.AliAlert, aliAlert.Env, aliAlert.currentDutyUser)
	}
	return

}

// firingMsgTemplate 告警模版
func firingMsgTemplate(template dto2.FeiShuMsg, aliAlert dto2.AliAlert, env string, dutyUser string) dto2.FeiShuMsg {

	alertTime, _ := strconv.ParseInt(aliAlert.AlertTime, 10, 64)

	GeneratorURL := globals.Config.Jaeger.URL + "/" + "trace/" + aliAlert.TraceID

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
								Content: "**🫧 报警指纹：**\n" + aliAlert.Fingerprint,
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
								Content: "**🤖 报警类型：**\n" + aliAlert.Name,
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
								Content: "**🕘 开始时间：**\n" + time.Unix(alertTime, 0).Format(globals.Layout),
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
								Content: "**📌 报警环境：**\n" + env,
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
								Content: "**🆔 TraceID：**\n" + aliAlert.TraceID,
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
								Content: "**🖥 报警主机：**\n" + aliAlert.Host,
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
								Content: "**📝 链路事件：**\n" + aliAlert.Attribute,
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
						Content: "⛓️ 链路查询",
						Tag:     "plain_text",
					},
					Type: "primary",
					MultiURL: &dto2.MultiURLs{
						URL: GeneratorURL,
					},
					Confirm: dto2.Confirms{
						Title: dto2.Titles{
							Content: "确认",
							Tag:     "plain_text",
						},
						Text: dto2.Texts{
							Content: fmt.Sprintf("查询当前链路 ID: %s 的详情", aliAlert.Fingerprint),
							Tag:     "plain_text",
						},
					},
				},
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
	template.Card.Header.Title.Content = "【报警中】链路报警 - 即时设计 🔥"
	template.Card.Elements = elements

	return template

}
