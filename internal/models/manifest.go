package models

type Manifest struct {
	Version    string                  `json:"version"`
	Libs       map[string]LibDef       `json:"libs"`
	Frameworks map[string]FrameworkDef `json:"frameworks"`
}

type LibDef struct {
	Imports       []string `json:"imports"`
	ConfigSection string   `json:"config_section"`
	Templates     []string `json:"templates"`
	Category      string   `json:"category,omitempty"`     // e.g., "database", "caching", "utilities"
	DisplayName   string   `json:"display_name,omitempty"` // e.g., "PostgreSQL", "Redis"
	Icon          string   `json:"icon,omitempty"`         // e.g., "🐘", "🔴"
	IsRadio       bool     `json:"is_radio,omitempty"`     // true for radio (mutually exclusive), false for checkbox
}

type FrameworkDef struct {
	Imports       []string `json:"imports"`
	ConfigSection string   `json:"config_section,omitempty"`
	Templates     []string `json:"templates"`
	DisplayName   string   `json:"display_name,omitempty"` // e.g., "Gin", "Echo"
	Icon          string   `json:"icon,omitempty"`         // e.g., "🍸", "🔊"
}
