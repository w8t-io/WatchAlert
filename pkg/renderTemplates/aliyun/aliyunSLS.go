package aliyun

import (
	"encoding/json"
	"fmt"
	"prometheus-manager/globals"
	"prometheus-manager/models"
	"prometheus-manager/utils"
	"strings"
)

type renderALiYun struct {
	AliAlert models.AliAlert
	Env      string
}

func ALiYunSlSTemplate(body []byte, env string) *renderALiYun {

	var (
		jsonArray []string
		AliAlert  models.AliAlert
	)

	err := json.Unmarshal(body, &jsonArray)
	if err != nil {
		globals.Logger.Sugar().Error("jsonArray 告警信息解析失败 ->", err)
		return &renderALiYun{}
	}

	jsonArray[0] = strings.ReplaceAll(jsonArray[0], `""`, `"`)

	err = json.Unmarshal([]byte(jsonArray[0]), &AliAlert)
	if err != nil {
		globals.Logger.Sugar().Error("AliAlert 告警信息解析失败 ->", err)
		return &renderALiYun{}
	}

	return &renderALiYun{
		AliAlert: AliAlert,
		Env:      env,
	}

}

func (ra *renderALiYun) FeiShu() error {

	var (
		cardContentMsg []string
		f              ALiYun
	)

	msg := f.FeiShuMsgTemplate(*ra)
	contentJson, _ := json.Marshal(msg.Card)
	// 需要将所有换行符进行转义
	cardContentJson := strings.Replace(string(contentJson), "\n", "\\n", -1)
	cardContentMsg = append(cardContentMsg, cardContentJson)
	err := utils.PushFeiShu(cardContentMsg)
	if err != nil {
		return fmt.Errorf("飞书消息发送失败 -> %s", err)
	}

	return nil

}
