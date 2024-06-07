package response

type RegionRes struct {
	List []Regions `json:"list"`
}

type Regions struct {
	Label *string `json:"label,omitempty"`
	Value *string `json:"value,omitempty"`
}

type MetricTypesRes struct {
	List []string `json:"list"`
}
