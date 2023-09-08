package sendAlertMessage

import (
	"prometheus-manager/globals"
	"strconv"
)

var (
	AlertType  string
	RespBody   []byte
	DataSource string
)

var (
	layout        string
	silenceTime   int64
	confirmPrompt string
)

func initBasic() {

	layout = "2006-01-02T15:04:05.000Z"
	silenceTime = globals.Config.AlertManager.SilenceTime
	confirmPrompt = "静默 " + strconv.FormatInt(globals.Config.AlertManager.SilenceTime, 10) + " 分钟"

}
