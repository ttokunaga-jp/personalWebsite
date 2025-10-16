package model

// LocalizedText represents a translatable string with Japanese and English variants.
type LocalizedText struct {
	Ja string `json:"ja,omitempty"`
	En string `json:"en,omitempty"`
}

// NewLocalizedText is a helper to construct LocalizedText ensuring consistent ordering.
func NewLocalizedText(ja, en string) LocalizedText {
	return LocalizedText{
		Ja: ja,
		En: en,
	}
}
