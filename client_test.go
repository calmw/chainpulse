package chainpulse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignRequest(t *testing.T) {
	got := SignRequest("secret", "get", "/v1/events?limit=1", "1710000000", "nonce")
	want := "sha256=036aaa9846eb8f5686cbb044f9e14f3737b7318c0e0d652cd0453a34c87e8b8d"
	if got != want {
		t.Fatalf("signature mismatch: got %s want %s", got, want)
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	body := []byte(`{"eventType":"coin.balance_changed"}`)
	signature := SignWebhookPayload("whsec_test", timestamp, body)
	if !VerifyWebhookSignature("whsec_test", timestamp, body, signature, 5*time.Minute) {
		t.Fatal("expected valid webhook signature")
	}
	if VerifyWebhookSignature("wrong", timestamp, body, signature, 5*time.Minute) {
		t.Fatal("expected invalid webhook signature")
	}
}

func TestDefaultHTTPClientHasTimeout(t *testing.T) {
	client := NewClient()
	if client.httpClient == nil {
		t.Fatal("expected default http client")
	}
	if client.httpClient.Timeout == 0 {
		t.Fatal("expected default timeout")
	}
}

func TestRequestSetsUserAgent(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{"status":"ok","service":"query"}`))
	}))
	defer server.Close()

	client := NewClient(WithQueryBaseURL(server.URL))
	_, err := client.Health(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if gotUserAgent == "" {
		t.Fatal("expected user agent")
	}
}

func TestWebhookEventsUnmarshalArray(t *testing.T) {
	var webhook Webhook
	err := json.Unmarshal([]byte(`{"id":"1","url":"https://example.com","events":["coin.balance_changed"],"isActive":true}`), &webhook)
	if err != nil {
		t.Fatal(err)
	}
	if len(webhook.Events) != 1 || webhook.Events[0] != "coin.balance_changed" {
		t.Fatalf("unexpected events: %#v", webhook.Events)
	}
}
