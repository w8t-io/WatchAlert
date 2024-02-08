package aliyun

type AliAlert struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Region      string `json:"region"`
	Status      string `json:"status"`
	AlertTime   string `json:"alert_time"`
	FireTime    string `json:"fire_time"`
	ResolveTime string `json:"resolve_time"`
	Host        string `json:"host"`
	TraceID     string `json:"traceID"`
	StatusCode  string `json:"statusCode"`
	Logs        string `json:"logs"`
	Attribute   string `json:"attribute"`
}
