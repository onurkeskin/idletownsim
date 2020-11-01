package effectview

type EffectView struct {
	EffectType      string     `json:"effecttype"`
	AreaArr         [][]string `json:"spaceeffectarea, omitempty"`
	Scheme          []string   `json:"effectscheme"`
	TargetResources []string   `json:"targetresources"`
}
