package sendAlertMessage

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"prometheus-manager/globals"
	"prometheus-manager/models"
)

type FeiShu struct{}

func (f *FeiShu) PushFeiShu(cardContentJson []string) error {

	for _, v := range cardContentJson {
		client := lark.NewClient(globals.Config.FeiShu.AppID, globals.Config.FeiShu.AppSecret, lark.WithEnableTokenCache(true))

		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(`chat_id`).
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(globals.Config.FeiShu.ChatID).
				MsgType(`interactive`).
				Content(v).
				Build()).
			Build()

		resp, err := client.Im.Message.Create(context.Background(), req, larkcore.WithTenantAccessToken(globals.Config.FeiShu.Token))
		// 处理错误
		if err != nil {
			globals.Logger.Sugar().Error("消息卡片发送失败 ->", err)
			return fmt.Errorf("消息卡片发送失败 -> %s", err)
		}

		// 服务端错误处理
		if !resp.Success() {
			globals.Logger.Sugar().Error(resp.Code, resp.Msg, resp.RequestId())
			return fmt.Errorf("响应错误 -> %s", err)
		}

		globals.Logger.Sugar().Info("消息卡片发送成功 ->", string(resp.RawBody))
	}

	return nil
}

func feiShuMsgTemplate(v models.Alerts, ActionsValueStr models.CreateAlertSilence, confirmPrompt string) (msg models.FeiShuMsg) {

	firingMsg := models.FeiShuMsg{
		MsgType: "interactive",
		Card: models.Cards{
			Config: models.Configs{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Elements: []models.Elements{
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
										Content: "**🕘 开始时间：**\n" + v.StartsAt,
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
										Content: "**🕟 结束时间：**\n" + v.EndsAt,
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
						Content: "🐾 执行动作：",
						Tag:     "plain_text",
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
			},
			Header: models.Headers{
				Template: "red",
				Title: models.Titles{
					Content: "【报警中】一级报警 - 即时设计 🔥",
					Tag:     "plain_text",
				},
			},
		},
	}
	resolvedMsg := models.FeiShuMsg{
		MsgType: "interactive",
		Card: models.Cards{
			Config: models.Configs{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Elements: []models.Elements{
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
										Content: "**🕘 开始时间：**\n" + v.StartsAt,
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
										Content: "**🕟 结束时间：**\n" + v.EndsAt,
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
			},
			Header: models.Headers{
				Template: "green",
				Title: models.Titles{
					Content: "【已处理】一级报警 - 即时设计 ✨",
					Tag:     "plain_text",
				},
			},
		},
	}

	silenceMsgContent := fmt.Sprintf("静默时长: %v 分钟\n结束时间: %v\n", globals.Config.AlertManager.SilenceTime, ActionsValueStr.EndsAt)
	silenceMsg := models.FeiShuMsg{
		MsgType: "interactive",
		Card: models.Cards{
			Config: models.Configs{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Elements: []models.Elements{
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
										Content: "**🕘 开始时间：**\n" + v.StartsAt,
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
										Content: "**🕟 结束时间：**\n" + v.EndsAt,
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
			},
			Header: models.Headers{
				Template: "yellow",
				Title: models.Titles{
					Content: "【静默中】一级报警 - 即时设计 🧘",
					Tag:     "plain_text",
				},
			},
		},
	}

	switch v.Status {
	case "firing":
		return firingMsg
	case "resolved":
		return resolvedMsg
	case "silence":
		return silenceMsg
	}
	return

}
