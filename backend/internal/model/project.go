package model

// Project encapsulates the metadata shown on the projects page.
type Project struct {
	ID          int64         `json:"id"`
	Title       LocalizedText `json:"title"`
	Description LocalizedText `json:"description"`
	TechStack   []string      `json:"techStack"`
	LinkURL     string        `json:"linkUrl"`
	Year        int           `json:"year"`
}
