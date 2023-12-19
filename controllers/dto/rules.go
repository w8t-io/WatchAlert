package dto

type AlertRules struct {
	Groups []Groups `json:"groups"`
}

type Groups struct {
	Name  string  `json:"name"`
	Rules []Rules `json:"rules"`
}

type Rules struct {
	Alert       string            `json:"alert"`
	Expr        string            `json:"expr"`
	For         string            `json:"for"`
	Labels      map[string]string `json:"labels"`
	Annotations Annotations       `json:"annotations"`
}

type RuleGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
