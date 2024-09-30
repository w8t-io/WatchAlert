package sender

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"watchAlert/alert/mute"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
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
		err := SendToEmail(alert, notice.Email.Subject, notice.Email.To, notice.Email.CC, n.CardContentMsg)
		if err != nil {
			return fmt.Errorf("邮件发送失败, err: %s", err.Error())
		}
	case "FeiShu", "DingDing":
		var msg string
		cardContentByte := bytes.NewReader([]byte(n.CardContentMsg))
		res, err := http.Post(nil, notice.Hook, cardContentByte)
		if err != nil {
			msg = err.Error()
		} else {
			if res.StatusCode != 200 {
				all, err := io.ReadAll(res.Body)
				if err != nil {
					global.Logger.Sugar().Error(err.Error())
					return err
				}
				msg = string(all)
			}
		}

		if msg != "" {
			global.Logger.Sugar().Errorf("Hook 类型报警发送失败 data: %s", n.CardContentMsg)
			return errors.New(msg)
		}
	default:
		return errors.New("无效的通知类型: " + notice.NoticeType)
	}

	global.Logger.Sugar().Info("报警发送成功: ", n.CardContentMsg)
	return nil
}
