package sender

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
)

type SendParmas struct {
	// 基础
	TenantId string
	RuleName string
	Severity string
	// 通知
	NoticeType string
	NoticeId   string
	NoticeName string
	// 恢复通知
	IsRecovered bool
	// hook 地址
	Hook string
	// 邮件
	Email models.Email
	// 消息
	Content string
	// 事件
	Event interface{}
}

func Sender(ctx *ctx.Context, sendParmas SendParmas) error {
	NoticeType := sendParmas.NoticeType
	var sendFunc func() error
	switch NoticeType {
	case "Email":
		sendFunc = func() error {
			return SendToEmail(sendParmas.IsRecovered, sendParmas.Email.Subject, sendParmas.Email.To, sendParmas.Email.CC, sendParmas.Content)
		}
	case "FeiShu":
		sendFunc = func() error {
			return SendToFeiShu(sendParmas.Hook, sendParmas.Content)
		}
	case "DingDing":
		sendFunc = func() error {
			return SendToDingDing(sendParmas.Hook, sendParmas.Content)
		}
	default:
		return fmt.Errorf("Send alarm failed, exist 无效的通知类型: %s", sendParmas.NoticeType)
	}

	if err := sendFunc(); err != nil {
		addRecord(ctx, sendParmas, 1, sendParmas.Content, err.Error())
		return fmt.Errorf("Send alarm failed to %s, err: %s", sendParmas.NoticeType, err.Error())
	}

	addRecord(ctx, sendParmas, 0, sendParmas.Content, "")
	logc.Info(ctx.Ctx, fmt.Sprintf("Send alarm ok, msg: %s", sendParmas.Content))
	return nil
}

func addRecord(ctx *ctx.Context, sendParmas SendParmas, status int, msg, errMsg string) {
	err := ctx.DB.Notice().AddRecord(models.NoticeRecord{
		Date:     time.Now().Format("2006-01-02"),
		CreateAt: time.Now().Unix(),
		TenantId: sendParmas.TenantId,
		RuleName: sendParmas.RuleName,
		NType:    sendParmas.NoticeType,
		NObj:     sendParmas.NoticeName + " (" + sendParmas.NoticeId + ")",
		Severity: sendParmas.Severity,
		Status:   status,
		AlarmMsg: msg,
		ErrMsg:   errMsg,
	})
	if err != nil {
		logc.Errorf(ctx.Ctx, fmt.Sprintf("Add notice record failed, err: %s", err.Error()))
	}
}
