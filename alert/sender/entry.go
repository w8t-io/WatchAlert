package sender

import (
	"fmt"
	"time"
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
	var sendFunc func() error
	switch NoticeType {
	case "Email":
		sendFunc = func() error {
			return SendToEmail(alert, notice.Email.Subject, notice.Email.To, notice.Email.CC, n.CardContentMsg)
		}
	case "FeiShu":
		sendFunc = func() error {
			return SendToFeiShu(notice.Hook, n.CardContentMsg)
		}
	case "DingDing":
		sendFunc = func() error {
			return SendToDingDing(notice.Hook, n.CardContentMsg)
		}
	default:
		return fmt.Errorf("Send alarm failed, exist 无效的通知类型: %s, NoticeId: %s, NoticeName: %s", notice.NoticeType, notice.Uuid, notice.Name)
	}

	if err := sendFunc(); err != nil {
		addRecord(ctx, alert, notice, 1, n.CardContentMsg, err.Error())
		return fmt.Errorf("Send alarm failed to %s, err: %s", notice.NoticeType, err.Error())
	}

	addRecord(ctx, alert, notice, 0, n.CardContentMsg, "")
	global.Logger.Sugar().Info("Send alarm ok, msg: ", n.CardContentMsg)
	return nil
}

func addRecord(ctx *ctx.Context, alert models.AlertCurEvent, notice models.AlertNotice, status int, msg, errMsg string) {
	err := ctx.DB.Notice().AddRecord(models.NoticeRecord{
		Date:     time.Now().Format("2006-01-02"),
		CreateAt: time.Now().Unix(),
		TenantId: alert.TenantId,
		RuleName: alert.RuleName,
		NType:    notice.NoticeType,
		NObj:     notice.Name + " (" + notice.Uuid + ")",
		Severity: alert.Severity,
		Status:   status,
		AlarmMsg: msg,
		ErrMsg:   errMsg,
	})
	if err != nil {
		global.Logger.Sugar().Errorf("Add notice record failed, err: %s", err.Error())
	}
}
