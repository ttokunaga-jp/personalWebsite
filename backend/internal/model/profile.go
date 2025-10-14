package model

// Profile represents the public-facing profile summary.
type Profile struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Affiliation string   `json:"affiliation"`
	Lab         string   `json:"lab"`
	Summary     string   `json:"summary"`
	Skills      []string `json:"skills"`
}
