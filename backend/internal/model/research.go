package model

// Research summary content served to the public site.
type Research struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	ContentMD string `json:"contentMd"`
	Year      int    `json:"year"`
}
