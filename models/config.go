package models

type App struct {
	Server       Server       `json:"server"`
	FeiShu       FeiShu       `json:"feishu"`
	AlertManager AlertManager `json:"AlertManager"`
	Prometheus   Prometheus   `json:"Prometheus"`
	Jaeger       Jaeger       `json:"Jaeger"`
	MySQL        MySQL        `json:"MySQL"`
}

type Server struct {
	Port string `json:"port"`
}

type FeiShu struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"sppSecret"`
	ChatID    string `json:"chatId"`
	Token     string `json:"token"`
}

type AlertManager struct {
	URL         string `json:"url"`
	SilenceTime int64  `json:"silenceTime"`
}

type Prometheus struct {
	URL      string `json:"url"`
	RulePath string `yaml:"rulePath"`
}

type Jaeger struct {
	URL string `json:"url"`
}

type MySQL struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	DBName  string `json:"dbName"`
	Timeout string `json:"timeout"`
}
