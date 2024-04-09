package notice

import (
	"bytes"
	"watchAlert/alert/mute"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/http"
)

type EntryNotice interface {
	NewTemplate(alert models.AlertCurEvent, notice models.AlertNotice) Template
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
	res, err := http.Post(notice.Hook, cardContentByte)
	if err != nil || res.StatusCode != 200 {
		globals.Logger.Sugar().Errorf("Hook 发送失败 -> code: %d data: %s", res.StatusCode, cardContentMsg)
		return err
	}

	globals.Logger.Sugar().Info("Hook 发送成功 ->", cardContentMsg)

	return nil

}
