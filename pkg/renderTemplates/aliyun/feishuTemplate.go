package aliyun

import (
	"fmt"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"prometheus-manager/pkg/schedule"
	"strconv"
	"time"
)

type ALiYun struct{}

// FeiShuMsgTemplate 飞书消息卡片模版
func (a *ALiYun) FeiShuMsgTemplate(aliAlert renderALiYun) (msg models.FeiShuMsg) {

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

	switch aliAlert.AliAlert.Status {
	case "firing":
		return firingMsgTemplate(defaultTemplate, aliAlert.AliAlert, aliAlert.Env)
	}
	return

}

// firingMsgTemplate 告警模版
func firingMsgTemplate(template models.FeiShuMsg, aliAlert models.AliAlert, env string) models.FeiShuMsg {

	var contentInfo string
	alertTime, _ := strconv.ParseInt(aliAlert.AlertTime, 10, 64)

	currentTime := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Day())

	_, userInfo := schedule.GetCurrentDutyInfo(currentTime)
	if len(userInfo.FeiShuUserID) == 0 {
		contentInfo = "暂无安排值班人员"
	} else {
		contentInfo = fmt.Sprintf("**👤 值班人员：**<at id=%s></at>", userInfo.FeiShuUserID)
	}

	GeneratorURL := globals.Config.Jaeger.URL + "/" + "trace/" + aliAlert.TraceID

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
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
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
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
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
					Elements: []models.ColumnsElements{
						{
							Tag: "div",
							Text: models.Texts{
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
						Content: "⛓️ 链路查询",
						Tag:     "plain_text",
					},
					Type: "primary",
					MultiURL: models.MultiURLs{
						URL: GeneratorURL,
					},
					Confirm: models.Confirms{
						Title: models.Titles{
							Content: "确认",
							Tag:     "plain_text",
						},
						Text: models.Texts{
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
			Elements: []models.ElementsElements{
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
