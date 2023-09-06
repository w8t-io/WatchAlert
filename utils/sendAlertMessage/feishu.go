package sendAlertMessage

import (
	"context"
	"encoding/json"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"strconv"
	"strings"
	"time"
)

type FeiShu struct{}

func (f *FeiShu) PushFeiShu(alertMsg map[string]interface{}) {

	var (
		alerts       models.Alert
		actionValues models.CreateAlertSilence
	)

	alertMsgJson, _ := json.Marshal(alertMsg)
	err := json.Unmarshal(alertMsgJson, &alerts)
	if err != nil {
		return
	}
	globals.Logger.Sugar().Info("告警原数据 ->", string(alertMsgJson))

	layout := "2006-01-02T15:04:05.000Z"
	silenceTime := globals.Config.AlertManager.SilenceTime

	for _, v := range alerts.Alerts {
		var MatchersList []models.Matchers
		for kk, vv := range v.Labels {
			Matchers := models.Matchers{
				Name:    kk,
				Value:   vv,
				IsEqual: true,
				IsRegex: false,
			}
			MatchersList = append(MatchersList, Matchers)
		}

		actionValues = models.CreateAlertSilence{
			Comment:   v.Fingerprint,
			CreatedBy: "1",
			EndsAt:    time.Now().Add(time.Minute * time.Duration(silenceTime)).Format(layout),
			ID:        "",
			Matchers:  MatchersList,
			StartsAt:  v.StartsAt,
		}
	}

	actionValueJson, _ := json.Marshal(actionValues)
	var ActionsValueStr models.CreateAlertSilence
	_ = json.Unmarshal(actionValueJson, &ActionsValueStr)

	// 消息提示
	msgContent := "静默 " + strconv.FormatInt(globals.Config.AlertManager.SilenceTime, 10) + " 分钟"
	fmt.Println(msgContent)

	for _, v := range alerts.Alerts {

		msg := feiShuMsgTemplate(v, ActionsValueStr, msgContent)

		contentJson, _ := json.Marshal(msg.Card)

		// 需要将所有换行符进行转义
		strContentJson := strings.Replace(string(contentJson), "\n", "\\n", -1)

		client := lark.NewClient(globals.Config.FeiShu.AppID, globals.Config.FeiShu.AppSecret, lark.WithEnableTokenCache(true))

		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(`chat_id`).
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(globals.Config.FeiShu.ChatID).
				MsgType(`interactive`).
				Content(strContentJson).
				Build()).
			Build()

		resp, err := client.Im.Message.Create(context.Background(), req, larkcore.WithTenantAccessToken(globals.Config.FeiShu.Token))
		// 处理错误
		if err != nil {
			globals.Logger.Sugar().Error("消息卡片发送失败 ->", err)
			return
		}

		// 服务端错误处理
		if !resp.Success() {
			globals.Logger.Sugar().Error(resp.Code, resp.Msg, resp.RequestId())
			return
		}

		globals.Logger.Sugar().Info("消息卡片发送成功 ->", string(resp.RawBody))

	}

}

func feiShuMsgTemplate(v models.Alerts, ActionsValueStr models.CreateAlertSilence, msgContent string) (msg models.FeiShuMsg) {

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
									Content: msgContent,
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
					Tag: "div",
					Text: models.Texts{
						Content: "🧑‍💻 即时设计 - 运维团队",
						Tag:     "plain_text",
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
					Tag: "div",
					Text: models.Texts{
						Content: "🧑‍💻 即时设计 - 运维团队",
						Tag:     "plain_text",
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

	switch v.Status {
	case "firing":
		return firingMsg
	case "resolved":
		return resolvedMsg
	}
	return

}
