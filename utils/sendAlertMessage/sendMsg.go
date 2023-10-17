package sendAlertMessage

import (
	"encoding/json"
	"fmt"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"strings"
	"time"
)

var (
	f FeiShu
)

func SendMsg(actionUser string, dataSource, alertType string, resp map[string]interface{}) error {

	initBasic()

	var cardContentMsg []string
	switch dataSource {
	case "prometheus":
		Msg, err := renderPrometheusMsgTemplate(actionUser, resp)
		if err != nil {
			return fmt.Errorf("消息渲染失败 -> %s", err)
		}
		cardContentMsg = Msg
	default:
		globals.Logger.Sugar().Error("无效的消息卡片模版")
		return fmt.Errorf("无效的消息卡片模版")
	}

	switch alertType {
	case "feishu":
		err := f.PushFeiShu(cardContentMsg)
		if err != nil {
			return fmt.Errorf("飞书消息发送失败 -> %s", err)
		}
	default:
		globals.Logger.Sugar().Error("无效的告警推送类型")
		return fmt.Errorf("无效的告警推送类型")
	}

	return nil

}

func renderPrometheusMsgTemplate(actionUser string, alertMsg map[string]interface{}) ([]string, error) {

	var (
		alerts         models.Alert
		actionValues   models.CreateAlertSilence
		cardContentMsg []string
	)

	alertMsgJson, _ := json.Marshal(alertMsg)
	err := json.Unmarshal(alertMsgJson, &alerts)
	if err != nil {
		globals.Logger.Sugar().Error("数据解析失败, 无法进行渲染消息模版", err)
		return nil, err
	}
	globals.Logger.Sugar().Info("告警原数据 ->", string(alertMsgJson))

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

		nowTimeUTC := time.Now().UTC().Add(time.Minute * time.Duration(silenceTime)).Format(layout)
		actionValues = models.CreateAlertSilence{
			Comment:   v.Fingerprint,
			CreatedBy: "1",
			EndsAt:    nowTimeUTC,
			ID:        "",
			Matchers:  MatchersList,
			StartsAt:  v.StartsAt.UTC().Format(layout),
		}

		globals.CacheCli.Set(v.Fingerprint, v)

		msg := feiShuMsgTemplate(actionUser, v, actionValues, confirmPrompt)

		contentJson, _ := json.Marshal(msg.Card)

		// 需要将所有换行符进行转义
		cardContentJson := strings.Replace(string(contentJson), "\n", "\\n", -1)

		cardContentMsg = append(cardContentMsg, cardContentJson)

	}

	return cardContentMsg, nil

}
