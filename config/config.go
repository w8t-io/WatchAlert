package config

type App struct {
	Server Server `json:"Server"`
	MySQL  MySQL  `json:"MySQL"`
	Redis  Redis  `json:"Redis"`
	Jwt    Jwt    `json:"Jwt"`
}

type Server struct {
	Port          string `json:"port"`
	GroupInterval int    `json:"groupInterval"`
}

type MySQL struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	DBName  string `json:"dbName"`
	Timeout string `json:"timeout"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Pass string `json:"pass"`
}

type Jwt struct {
	Expire int64 `json:"expire"`
}
