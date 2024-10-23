package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"watchAlert/pkg/utils/http"
)

type FeiShuResponseMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendToFeiShu(hook, msg string) error {
	cardContentByte := bytes.NewReader([]byte(msg))
	res, err := http.Post(nil, hook, cardContentByte)
	if err != nil {
		return err
	}

	// 读取响应体内容
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading Feishu response body: %s, msg: %s", string(body), err.Error()))
	}

	var response FeiShuResponseMsg
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling Feishu response: %s", err.Error()))
	}
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
