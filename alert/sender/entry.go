package sender

import (
	"bytes"
	"errors"
	"watchAlert/alert/mute"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/client"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/http"
	"watchAlert/pkg/utils/templates"
)

func Sender(ctx *ctx.Context, alert models.AlertCurEvent, notice models.AlertNotice) error {
	ok := mute.IsMuted(ctx, &alert)
	if ok {
		return nil
	}

	n := templates.NewTemplate(ctx, alert, notice)

	switch notice.NoticeType {
	case "Email":
		setting, err := ctx.DB.Setting().Get()
		if err != nil {
			return errors.New("获取系统配置失败: " + err.Error())
		}
		eCli := client.NewEmailClient(setting.EmailConfig.ServerAddress, setting.EmailConfig.Email, setting.EmailConfig.Token, setting.EmailConfig.Port)
		if alert.IsRecovered {
			notice.Email.Subject = notice.Email.Subject + "「已恢复」"
		} else {
			notice.Email.Subject = notice.Email.Subject + "「报警中」"
		}
		err = eCli.Send(notice.Email.To, notice.Email.CC, notice.Email.Subject, []byte(n.CardContentMsg))
		if err != nil {
			global.Logger.Sugar().Error("Email 类型报警发送失败: " + err.Error() + ", Content: " + n.CardContentMsg)
			return err
		}
	case "FeiShu", "DingDing":
		cardContentByte := bytes.NewReader([]byte(n.CardContentMsg))
		res, err := http.Post(nil, notice.Hook, cardContentByte)
		if err != nil || res.StatusCode != 200 {
			global.Logger.Sugar().Errorf("Hook 类型报警发送失败 code: %d data: %s", res.StatusCode, n.CardContentMsg)
			return err
		}
	default:
		return errors.New("无效的通知类型: " + notice.NoticeType)
	}

	global.Logger.Sugar().Info("报警发送成功: ", n.CardContentMsg)
	return nil
}
