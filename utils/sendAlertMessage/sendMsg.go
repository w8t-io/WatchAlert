package sendAlertMessage

var (
	f FeiShu
)

func SendMsg(alertType string, alertMsg map[string]interface{}) {

	switch alertType {
	case "feishu":
		f.PushFeiShu(alertMsg)
	}

}
