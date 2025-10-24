package mail

import "context"

// Message represents an email to be delivered.
type Message struct {
	From    string
	To      []string
	CC      []string
	Subject string
	Body    string
}

// Client delivers email messages.
type Client interface {
	Send(ctx context.Context, message Message) error
}
