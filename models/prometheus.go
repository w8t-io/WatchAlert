package models

// Alert 接收Prometheus推送的告警
type Alert struct {
	Alerts   []Alerts `json:"alerts"`
	Receiver string   `json:"receiver"`
}

type Alerts struct {
	Annotations  Annotations       `json:"annotations"`
	EndsAt       string            `json:"endsAt"`
	Fingerprint  string            `json:"fingerprint"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     string            `json:"startsAt"`
	Status       string            `json:"status"`
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
