package chainpulse

import (
	"net/http"
	"strings"
	"time"
)

const (
	defaultQueryBaseURL = "http://localhost:5001"
	defaultServerAPIURL = "http://localhost:3001/api"
)

type Client struct {
	queryBaseURL string
	serverAPIURL string
	apiKey       string
	apiSecret    string
	bearerToken  string
	httpClient   *http.Client
}

type Option func(*Client)

func NewClient(options ...Option) *Client {
	c := &Client{
		queryBaseURL: defaultQueryBaseURL,
		serverAPIURL: defaultServerAPIURL,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func WithQueryBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.queryBaseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	}
}

func WithServerAPIURL(baseURL string) Option {
	return func(c *Client) {
		c.serverAPIURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	}
}

func WithAPIKey(key string, secret string) Option {
	return func(c *Client) {
		c.apiKey = strings.TrimSpace(key)
		c.apiSecret = strings.TrimSpace(secret)
	}
}

func WithBearerToken(token string) Option {
	return func(c *Client) {
		c.bearerToken = strings.TrimSpace(token)
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}
