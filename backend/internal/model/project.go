package model

// Project encapsulates the metadata shown on the projects page.
type Project struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TechStack   []string `json:"techStack"`
	LinkURL     string   `json:"linkUrl"`
	Year        int      `json:"year"`
}
