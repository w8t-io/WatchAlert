package models

import "time"

// Alerts 接收Prometheus推送的告警
type Alerts struct {
	AlertList []AlertInfo `json:"alerts"`
	Receiver  string      `json:"receiver"`
}

type AlertInfo struct {
	Annotations  Annotations       `json:"annotations"`
	EndsAt       time.Time         `json:"endsAt"`
	Fingerprint  string            `json:"fingerprint"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     time.Time         `json:"startsAt"`
	Status       string            `json:"status"`
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}
