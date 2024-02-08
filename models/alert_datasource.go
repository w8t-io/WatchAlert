package models

type AlertDataSource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	HTTPJson    HTTP   `json:"http" gorm:"-"`
	HTTP        string `json:"-" gorm:"http"`
	Description string `json:"description"`
	EnabledBool bool   `json:"enabled" gorm:"-"`
	Enabled     string `json:"-" gorm:"enabled"`
}

type HTTP struct {
	URL     string `json:"url"`
	Timeout string `json:"timeout"`
}
