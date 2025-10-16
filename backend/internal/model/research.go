package model

// Research summary content served to the public site.
type Research struct {
	ID        int64         `json:"id"`
	Title     LocalizedText `json:"title"`
	Summary   LocalizedText `json:"summary"`
	ContentMD LocalizedText `json:"contentMd"`
	Year      int           `json:"year"`
}
