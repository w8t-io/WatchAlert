package models

type App struct {
	Server       Server       `json:"server"`
	FeiShu       FeiShu       `json:"feishu"`
	AlertManager AlertManager `json:"AlertManager"`
	DutyUser     []string     `json:"dutyUser"`
}

type Server struct {
	Port string `json:"port"`
}

type FeiShu struct {
	AppID      string   `json:"appId"`
	AppSecret  string   `json:"sppSecret"`
	ChatID     string   `json:"chatId"`
	Token      string   `json:"token"`
	ActionUser []string `json:"actionUser"`
}

type AlertManager struct {
	URL         string `json:"url"`
	SilenceTime int64  `json:"silenceTime"`
}
