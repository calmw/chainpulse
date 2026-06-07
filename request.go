package chainpulse

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const userAgent = "chainpulse-go/0.1"

func (c *Client) doQuery(ctx context.Context, method string, path string, query url.Values, body any, out any) error {
	if c.apiKey == "" || c.apiSecret == "" {
		return fmt.Errorf("chainpulse: missing API key or API secret")
	}
	return c.do(ctx, c.queryBaseURL, method, path, query, body, out, true)
}

func (c *Client) doServer(ctx context.Context, method string, path string, query url.Values, body any, out any) error {
	if c.bearerToken == "" {
		return fmt.Errorf("chainpulse: missing bearer token")
	}
	return c.do(ctx, c.serverAPIURL, method, path, query, body, out, false)
}

func (c *Client) do(ctx context.Context, baseURL string, method string, path string, query url.Values, body any, out any, sign bool) error {
	if strings.TrimSpace(baseURL) == "" {
		return fmt.Errorf("chainpulse: missing base URL")
	}

	pathWithQuery := path
	if len(query) > 0 {
		pathWithQuery += "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(baseURL, "/")+pathWithQuery, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	if sign {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		nonce, err := randomHex(12)
		if err != nil {
			return err
		}
		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("X-API-Timestamp", timestamp)
		req.Header.Set("X-API-Nonce", nonce)
		req.Header.Set("X-API-Signature", SignRequest(c.apiSecret, method, pathWithQuery, timestamp, nonce))
	} else {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp.StatusCode, data)
	}
	if out == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		if strings.Contains(err.Error(), `cannot parse "`) && strings.Contains(err.Error(), `as "2006-01-02T15:04:05Z07:00"`) {
			return fmt.Errorf("decode response: %w (API timestamps use YYYY-MM-DD HH:mm:ss; upgrade chainpulse to v1.0.2+)", err)
		}
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func SignRequest(secret string, method string, pathWithQuery string, timestamp string, nonce string) string {
	signingString := strings.Join([]string{
		strings.ToUpper(method),
		pathWithQuery,
		timestamp,
		nonce,
	}, "\n")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingString))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func randomHex(size int) (string, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func pathEscape(value string) string {
	return url.PathEscape(strings.TrimSpace(value))
}
