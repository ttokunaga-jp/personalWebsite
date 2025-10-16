package model

// Profile represents the public-facing profile summary enriched with localisation metadata.
type Profile struct {
	Name        LocalizedText   `json:"name"`
	Title       LocalizedText   `json:"title"`
	Affiliation LocalizedText   `json:"affiliation"`
	Lab         LocalizedText   `json:"lab"`
	Summary     LocalizedText   `json:"summary"`
	Skills      []LocalizedText `json:"skills"`
}
