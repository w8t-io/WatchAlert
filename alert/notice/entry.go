package notice

import (
	"bytes"
	"watchAlert/alert/mute"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/http"
)

type EntryNotice interface {
	NewTemplate(alert models.AlertCurEvent, notice models.AlertNotice) Prometheus
}

func NewEntryNotice(e EntryNotice, alert models.AlertCurEvent, notice models.AlertNotice) {

	ok := mute.IsMuted(&alert)
	if ok {
		return
	}

	go func() {
		n := e.NewTemplate(alert, notice)
		err := SendNotice(notice, n.CardContentMsg)
		if err != nil {
			return
		}
	}()

}

func SendNotice(notice models.AlertNotice, cardContentMsg string) error {

	cardContentByte := bytes.NewReader([]byte(cardContentMsg))
	_, err := http.Post(notice.Hook, cardContentByte)
	if err != nil {
		globals.Logger.Sugar().Error("Hook 发送失败 ->", err.Error())
		return err
	}

	globals.Logger.Sugar().Info("Hook 发送成功 ->", cardContentMsg)

	return nil
	
}
