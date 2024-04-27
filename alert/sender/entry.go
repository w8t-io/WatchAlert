package sender

import (
	"bytes"
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

	n := templates.NewTemplate(alert, notice)

	cardContentByte := bytes.NewReader([]byte(n.CardContentMsg))
	res, err := http.Post(notice.Hook, cardContentByte)
	if err != nil || res.StatusCode != 200 {
		global.Logger.Sugar().Errorf("Hook 发送失败 -> code: %d data: %s", res.StatusCode, n.CardContentMsg)
		return err
	}

	global.Logger.Sugar().Info("Hook 发送成功 ->", n.CardContentMsg)
	return nil
}
