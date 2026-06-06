# ChainPulse Go SDK

Go SDK for ChainPulse query and webhook services.

```bash
go get github.com/calmw/chainpulse
```

## Layout

```text
.
├── client.go          # Client options and construction
├── request.go         # HTTP transport, HMAC signing, error handling
├── query.go           # Query API methods
├── query_params.go    # Query API common parameters
├── webhook.go         # Webhook management and signature helpers
├── models.go          # Public response/request models
├── errors.go          # APIError
└── examples/
    ├── query/
    └── webhook_receiver/
```

## Query API

Query service uses HMAC headers:

- `X-API-Key`
- `X-API-Timestamp`
- `X-API-Nonce`
- `X-API-Signature`

The SDK signs requests automatically.

```go
package main

import (
	"context"
	"fmt"

	chainpulse "github.com/calmw/chainpulse"
)

func main() {
	client := chainpulse.NewClient(
		chainpulse.WithQueryBaseURL("https://query.chainpulse.cc"),
		chainpulse.WithAPIKey("ck_live_xxx", "sk_live_xxx"),
	)

	tx, err := client.GetTransaction(context.Background(), "0xabc")
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.TxHash)
}
```

Common query methods (optional `QueryParams`: `chainId`, `consistency`, `includeRemoved`, `limit`):

- `GetTransaction(ctx, txHash, params...)`
- `ListTransactionInternalTransactions(ctx, txHash, params...)`
- `ListAddressTransactions(ctx, address, params...)`
- `ListAddressBalances(ctx, address, params...)`
- `ListNFTBalances(ctx, address, standard, params...)`
- `ListEvents(ctx, options, params...)`
- `ListContractEvents(ctx, contractAddress, options, params...)`

Example with consistency:

```go
tx, err := client.GetTransaction(ctx, "0xabc", chainpulse.QueryParams{
	ChainID:     56,
	Consistency: "stable",
})
```

## Webhook Management

Webhook subscription management uses the server API bearer token from user login.

```go
client := chainpulse.NewClient(
	chainpulse.WithServerAPIURL("https://chainpulse.cc/api"),
	chainpulse.WithBearerToken("user_jwt_token"),
)

created, err := client.CreateWebhook(context.Background(), chainpulse.WebhookConfig{
	URL:    "https://example.com/webhook",
	Events: []string{"coin.balance_changed", "token.balance_changed", "contract.event"},
})
if err != nil {
	panic(err)
}
fmt.Println(created.WebhookID, created.Secret)
```

## Examples

Query a transaction:

```bash
cd examples/query
CHAINPULSE_API_KEY=ck_live_xxx \
CHAINPULSE_API_SECRET=sk_live_xxx \
CHAINPULSE_TX_HASH=0x... \
go run .
```

Run a webhook receiver:

```bash
cd examples/webhook_receiver
CHAINPULSE_WEBHOOK_SECRET=whsec_xxx go run .
```

Webhook methods:

- `CreateWebhook(ctx, config)`
- `ListWebhooks(ctx)`
- `UpdateWebhook(ctx, webhookID, config)`
- `DeleteWebhook(ctx, webhookID)`

## Webhook Receiver

Use `VerifyWebhookSignature` to verify payloads received from the webhook service.

```go
func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	ok := chainpulse.VerifyWebhookSignature(
		"whsec_xxx",
		r.Header.Get("X-Bot-Timestamp"),
		body,
		r.Header.Get("X-Bot-Signature"),
		5*time.Minute,
	)
	if !ok {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	event, err := chainpulse.ParseWebhookEvent(body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	fmt.Println(event.EventType, event.TxHash)
}
```
