package sender

import (
	"errors"
	"fmt"
	"watchAlert/alert/mute"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	"watchAlert/pkg/utils/templates"
)

func Sender(ctx *ctx.Context, alert models.AlertCurEvent, notice models.AlertNotice) error {
	ok := mute.IsMuted(ctx, &alert)
	if ok {
		return nil
	}

	n := templates.NewTemplate(ctx, alert, notice)
	NoticeType := notice.NoticeType
	switch NoticeType {
	case "Email":
		err := SendToEmail(alert, notice.Email.Subject, notice.Email.To, notice.Email.CC, n.CardContentMsg)
		if err != nil {
			return fmt.Errorf("Send alarm failed to email , err: %s", err.Error())
		}
	case "FeiShu":
		err := SendToFeiShu(notice.Hook, n.CardContentMsg)
		if err != nil {
			return fmt.Errorf("Send alarm failed to FeiShu, err: %s", err.Error())
		}
	case "DingDing":
		err := SendToDingDing(notice.Hook, n.CardContentMsg)
		if err != nil {
			return fmt.Errorf("Send alarm failed to DingDing, err: %s", err.Error())
		}
	default:
		return errors.New(fmt.Sprintf("Send alarm failed, exist 无效的通知类型: %s, NoticeId: %s, NoticeName: %s", notice.NoticeType, notice.Uuid, notice.Name))
	}

	global.Logger.Sugar().Info("Send alarm ok, msg: ", n.CardContentMsg)
	return nil
}
