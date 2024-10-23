package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"watchAlert/pkg/utils/http"
)

type DingResponseMsg struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

func SendToDingDing(hook, msg string) error {
	cardContentByte := bytes.NewReader([]byte(msg))
	res, err := http.Post(nil, hook, cardContentByte)
	if err != nil {
		return err
	}

	// 读取响应体内容
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading Dingding response body: %s", err.Error()))
	}

	var response DingResponseMsg
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Dingding response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
