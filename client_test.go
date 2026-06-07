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

func TestParseWebhookEventDisplayTime(t *testing.T) {
	body := []byte(`{
		"eventId":"e1",
		"eventType":"coin.balance_changed",
		"chainId":968,
		"blockNumber":12035315,
		"blockHash":"0xabc",
		"txHash":"0xb9308cad0162f61d29b0fe0f892a83fcd62777f4a99a8c5abaf54eda94c2c131",
		"address":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723",
		"confirmationStatus":"stable",
		"confirmations":12,
		"payload":{"role":"internal_to"},
		"blockTime":"2026-06-07 18:30:00",
		"createdAt":"2026-06-08 02:12:25"
	}`)
	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatalf("ParseWebhookEvent failed: %v", err)
	}
	if event.EventType != "coin.balance_changed" || event.ChainID != 968 {
		t.Fatalf("unexpected event: %#v", event)
	}
}

func TestParseWebhookEventInternalTransferPayload(t *testing.T) {
	body := []byte(`{
		"eventId":"e2",
		"eventType":"coin.balance_changed",
		"chainId":968,
		"blockNumber":12035315,
		"blockHash":"0x579623b1962f73d0cbfdccc83a4d3d8b8a04d23d964ba50a06dcd30bb13daa1e",
		"txHash":"0xb9308cad0162f61d29b0fe0f892a83fcd62777f4a99a8c5abaf54eda94c2c131",
		"address":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723",
		"confirmationStatus":"stable",
		"confirmations":12,
		"finalized":0,
		"removed":0,
		"payload":{
			"matchedAddress":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723",
			"assetType":"coin",
			"role":"internal_to",
			"balance":{
				"chainId":968,
				"address":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723",
				"assetType":"coin",
				"balance":"1000000000000000000",
				"updatedBlock":12035315,
				"updatedAt":"2026-06-07 18:30:00"
			},
			"internalTransfer":{
				"txHash":"0xb9308cad0162f61d29b0fe0f892a83fcd62777f4a99a8c5abaf54eda94c2c131",
				"traceAddress":"0.0.0",
				"type":"CALL",
				"from":"0x3e817cf1b58d00a1787c3a49a2f405433cb89ace",
				"to":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723",
				"value":"1000000000000000000"
			},
			"blockTime":"2026-06-07 18:30:00",
			"subscriptionType":"account_coin",
			"webhookId":1
		},
		"blockTime":"2026-06-07 18:30:00",
		"createdAt":"2026-06-08 02:12:25"
	}`)
	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatalf("ParseWebhookEvent failed: %v", err)
	}
	if event.Payload["role"] != "internal_to" {
		t.Fatalf("unexpected role: %#v", event.Payload["role"])
	}
	if event.CreatedAt.IsZero() || event.BlockTime == nil || event.BlockTime.IsZero() {
		t.Fatalf("expected parsed timestamps: %#v", event)
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
