package sendAlertMessage

var (
	f FeiShu
)

func SendMsg(alertType string, resp map[string]interface{}) {

	switch alertType {
	case "feishu":
		f.PushFeiShu(resp)
	}

}
