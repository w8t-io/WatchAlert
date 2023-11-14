package renderTemplates

import (
	"encoding/json"
	"fmt"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"prometheus-manager/pkg/feishu"
	"strings"
	"time"
)

type renderPrometheus struct {
	MsgTemplate []msgTemplate
}

type msgTemplate struct {
	actionUser   string
	alerts       models.AlertInfo
	actionValues models.CreateAlertSilence
}

func Prometheus(actionUser string, alertMsg map[string]interface{}) *renderPrometheus {

	var (
		alerts       models.Alerts
		actionValues models.CreateAlertSilence
		rp           renderPrometheus
	)

	alertMsgJson, _ := json.Marshal(alertMsg)
	err := json.Unmarshal(alertMsgJson, &alerts)
	if err != nil {
		globals.Logger.Sugar().Error("数据解析失败, 无法进行渲染消息模版", err)
		return &renderPrometheus{}
	}
	globals.Logger.Sugar().Info("告警原数据 ->", string(alertMsgJson))

	for _, v := range alerts.AlertList {
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

		nowTimeUTC := time.Now().UTC().Add(time.Minute * time.Duration(globals.Config.AlertManager.SilenceTime)).Format(globals.Layout)
		actionValues = models.CreateAlertSilence{
			Comment:   v.Fingerprint,
			CreatedBy: "1",
			EndsAt:    nowTimeUTC,
			ID:        "",
			Matchers:  MatchersList,
			StartsAt:  v.StartsAt.UTC().Format(globals.Layout),
		}

		globals.CacheCli.Set(v.Fingerprint, v)

		mt := msgTemplate{
			actionUser:   actionUser,
			alerts:       v,
			actionValues: actionValues,
		}

		rp.MsgTemplate = append(rp.MsgTemplate, mt)

	}

	return &rp

}

func (r *renderPrometheus) FeiShu() error {

	var (
		cardContentMsg []string
		f              feishu.FeiShu
	)

	for _, i := range r.MsgTemplate {
		msg := f.FeiShuMsgTemplate(i.actionUser, i.alerts, i.actionValues)
		contentJson, _ := json.Marshal(msg.Card)
		// 需要将所有换行符进行转义
		cardContentJson := strings.Replace(string(contentJson), "\n", "\\n", -1)
		cardContentMsg = append(cardContentMsg, cardContentJson)
		err := f.PushFeiShu(cardContentMsg)
		if err != nil {
			return fmt.Errorf("飞书消息发送失败 -> %s", err)
		}
	}
	return nil
}

func (r *renderPrometheus) DingDing() error {

	//ToDo
	return nil

}
