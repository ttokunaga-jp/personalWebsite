package model

// ContactRequest represents a submission from the contact form.
type ContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

// ContactSubmission is a stub for persistence/queueing, ready for expansion.
type ContactSubmission struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Comment string `json:"comment"`
}
