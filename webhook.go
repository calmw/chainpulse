package chainpulse

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type WebhookConfig struct {
	URL     string   `json:"url,omitempty"`
	Events  []string `json:"events,omitempty"`
	ChainID *int     `json:"chainId,omitempty"`
}

type UpdateWebhookConfig struct {
	URL      *string  `json:"url,omitempty"`
	Events   []string `json:"events,omitempty"`
	ChainID  *int     `json:"chainId,omitempty"`
	IsActive *bool    `json:"isActive,omitempty"`
}

func (c *Client) CreateWebhook(ctx context.Context, config WebhookConfig) (CreateWebhookResponse, error) {
	var out CreateWebhookResponse
	err := c.doServer(ctx, http.MethodPost, "/webhooks", nil, config, &out)
	return out, err
}

func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	var out []Webhook
	err := c.doServer(ctx, http.MethodGet, "/webhooks", nil, nil, &out)
	return out, err
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, config UpdateWebhookConfig) (Webhook, error) {
	var out Webhook
	err := c.doServer(ctx, http.MethodPut, "/webhooks/"+pathEscape(webhookID), nil, config, &out)
	return out, err
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) (MessageResponse, error) {
	var out MessageResponse
	err := c.doServer(ctx, http.MethodDelete, "/webhooks/"+pathEscape(webhookID), nil, nil, &out)
	return out, err
}

func ParseWebhookEvent(body []byte) (WebhookEvent, error) {
	var event WebhookEvent
	err := json.Unmarshal(body, &event)
	return event, err
}

func VerifyWebhookSignature(secret string, timestamp string, body []byte, signature string, tolerance time.Duration) bool {
	if tolerance > 0 {
		parsed, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return false
		}
		if time.Since(parsed) > tolerance || time.Until(parsed) > tolerance {
			return false
		}
	}
	expected := SignWebhookPayload(secret, timestamp, body)
	return hmac.Equal([]byte(expected), []byte(strings.TrimSpace(signature)))
}

func SignWebhookPayload(secret string, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
