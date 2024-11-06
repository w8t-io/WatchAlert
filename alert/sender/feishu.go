package sender

import (
	"bytes"
	"errors"
	"fmt"
	"watchAlert/pkg/tools"
)

type FeiShuResponseMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendToFeiShu(hook, msg string) error {
	cardContentByte := bytes.NewReader([]byte(msg))
	res, err := tools.Post(nil, hook, cardContentByte)
	if err != nil {
		return err
	}

	var response FeiShuResponseMsg
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Feishu response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
