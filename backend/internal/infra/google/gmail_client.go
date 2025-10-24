package google

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/takumi/personal-website/internal/mail"
)

const gmailSendURL = "https://gmail.googleapis.com/gmail/v1/users/me/messages/send"

// GmailAPIClient implements mail.Client using the Gmail REST API.
type GmailAPIClient struct {
	client        *http.Client
	tokenProvider TokenProvider
}

// NewGmailAPIClient constructs a Gmail client backed by the provided HTTP client and token provider.
func NewGmailAPIClient(httpClient *http.Client, tokenProvider TokenProvider) mail.Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &GmailAPIClient{
		client:        httpClient,
		tokenProvider: tokenProvider,
	}
}

func (g *GmailAPIClient) Send(ctx context.Context, message mail.Message) error {
	token, err := g.tokenProvider.AccessToken(ctx)
	if err != nil {
		return err
	}

	rawMessage := buildMIMEMessage(message)
	payload := map[string]string{
		"raw": base64.URLEncoding.EncodeToString([]byte(rawMessage)),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gmail send marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, gmailSendURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gmail send request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("gmail send call: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return fmt.Errorf("gmail send error: status=%d body=%s", resp.StatusCode, string(payload))
	}

	return nil
}

func buildMIMEMessage(message mail.Message) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("From: %s\r\n", message.From))
	if len(message.To) > 0 {
		builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))
	}
	if len(message.CC) > 0 {
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.CC, ", ")))
	}
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	builder.WriteString(message.Body)
	return builder.String()
}
