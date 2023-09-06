package models

import "github.com/go-openapi/strfmt"

type SearchAlertManager struct {
	Status    Status     `json:"status"`
	UpdatedAt string     `json:"updatedAt"`
	Comment   string     `json:"comment"`
	CreatedBy string     `json:"createdBy"`
	EndsAt    string     `json:"endsAt"`
	ID        string     `json:"id"`
	Matchers  []Matchers `json:"matchers"`
	StartsAt  string     `json:"startsAt"`
}

type Status struct {
	State string `json:"state"`
}

type CreateAlertSilence struct {
	Comment   string     `json:"comment"`
	CreatedBy string     `json:"createdBy"`
	EndsAt    string     `json:"endsAt"`
	ID        string     `json:"id"`
	Matchers  []Matchers `json:"matchers"`
	StartsAt  string     `json:"startsAt"`
}

type Matchers struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsEqual bool   `json:"isEqual"`
	IsRegex bool   `json:"isRegex"`
}

// swagger:model gettableAlert
type GettableAlert struct {

	// annotations
	// Required: true
	Annotations LabelSet `json:"annotations"`

	// ends at
	// Required: true
	// Format: date-time
	EndsAt *strfmt.DateTime `json:"endsAt"`

	// fingerprint
	// Required: true
	Fingerprint *string `json:"fingerprint"`

	// receivers
	// Required: true
	Receivers []*Receiver `json:"receivers"`

	// starts at
	// Required: true
	// Format: date-time
	StartsAt *strfmt.DateTime `json:"startsAt"`

	// status
	// Required: true
	Status *AlertStatus `json:"status"`

	// updated at
	// Required: true
	// Format: date-time
	UpdatedAt *strfmt.DateTime `json:"updatedAt"`

	Alert
}

// swagger:model labelSet
type LabelSet map[string]string

// swagger:model receiver
type Receiver struct {

	// name
	// Required: true
	Name *string `json:"name"`
}

// swagger:model alertStatus
type AlertStatus struct {

	// inhibited by
	// Required: true
	InhibitedBy []string `json:"inhibitedBy"`

	// silenced by
	// Required: true
	SilencedBy []string `json:"silencedBy"`

	// state
	// Required: true
	// Enum: [unprocessed active suppressed]
	State *string `json:"state"`
}

// swagger:model alert
//type Alert struct {
//
//	// generator URL
//	// Format: uri
//	GeneratorURL strfmt.URI `json:"generatorURL,omitempty"`
//
//	// labels
//	// Required: true
//	Labels LabelSet `json:"labels"`
//}
