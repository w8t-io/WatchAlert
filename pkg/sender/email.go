package sender

import (
	"errors"
	"fmt"
	"watchAlert/pkg/client"
	"watchAlert/pkg/ctx"
)

func SendToEmail(IsRecovered bool, subject string, to, cc []string, msg string) error {
	setting, err := ctx.DB.Setting().Get()
	if err != nil {
		return errors.New("获取系统配置失败: " + err.Error())
	}
	eCli := client.NewEmailClient(setting.EmailConfig.ServerAddress, setting.EmailConfig.Email, setting.EmailConfig.Token, setting.EmailConfig.Port)
	if IsRecovered {
		subject = subject + "「已恢复」"
	} else {
		subject = subject + "「报警中」"
	}
	err = eCli.Send(to, cc, subject, []byte(msg))
	if err != nil {
		return fmt.Errorf("%s, %s", err.Error(), "Content: "+msg)
	}

	return nil
}
