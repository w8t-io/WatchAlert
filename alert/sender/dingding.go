package sender

import (
	"bytes"
	"errors"
	"fmt"
	"watchAlert/pkg/tools"
)

type DingResponseMsg struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

func SendToDingDing(hook, msg string) error {
	cardContentByte := bytes.NewReader([]byte(msg))
	res, err := tools.Post(nil, hook, cardContentByte)
	if err != nil {
		return err
	}

	var response DingResponseMsg
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Dingding response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
