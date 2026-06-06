package chainpulse

import "time"

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ChainTransaction struct {
	ChainID          uint64    `json:"chainId"`
	BlockNumber      uint64    `json:"blockNumber"`
	BlockHash        string    `json:"blockHash"`
	TxHash           string    `json:"txHash"`
	From             string    `json:"from"`
	To               string    `json:"to"`
	Value            string    `json:"value"`
	Nonce            uint64    `json:"nonce"`
	Gas              uint64    `json:"gas"`
	GasPrice         string    `json:"gasPrice"`
	Input            string    `json:"input,omitempty"`
	Status           uint8     `json:"status"`
	TransactionIndex uint64    `json:"transactionIndex"`
	BlockTime        time.Time `json:"blockTime"`
	Finalized        uint8     `json:"finalized"`
	Removed          uint8     `json:"removed"`
	Version          uint64    `json:"version"`
}

type InternalTransaction struct {
	ChainID      uint64    `json:"chainId"`
	BlockNumber  uint64    `json:"blockNumber"`
	BlockHash    string    `json:"blockHash"`
	TxHash       string    `json:"txHash"`
	TraceAddress string    `json:"traceAddress"`
	CallType     string    `json:"callType"`
	From         string    `json:"from"`
	To           string    `json:"to"`
	Value        string    `json:"value"`
	Removed      uint8     `json:"removed"`
	BlockTime    time.Time `json:"blockTime"`
	Version      uint64    `json:"version"`
}

type ChainLog struct {
	ChainID     uint64    `json:"chainId"`
	BlockNumber uint64    `json:"blockNumber"`
	BlockHash   string    `json:"blockHash"`
	TxHash      string    `json:"txHash"`
	LogIndex    uint64    `json:"logIndex"`
	Address     string    `json:"address"`
	Topic0      string    `json:"topic0"`
	Topic1      string    `json:"topic1"`
	Topic2      string    `json:"topic2"`
	Topic3      string    `json:"topic3"`
	Data        string    `json:"data"`
	EventName   string    `json:"eventName"`
	Removed     uint8     `json:"removed"`
	BlockTime   time.Time `json:"blockTime"`
	Version     uint64    `json:"version"`
}

type AddressBalance struct {
	ChainID      uint64    `json:"chainId"`
	Address      string    `json:"address"`
	AssetType    string    `json:"assetType"`
	TokenAddress string    `json:"tokenAddress,omitempty"`
	Balance      string    `json:"balance"`
	UpdatedBlock uint64    `json:"updatedBlock"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type NFTTokenBalance struct {
	ChainID       uint64    `json:"chainId"`
	OwnerAddress  string    `json:"ownerAddress"`
	TokenAddress  string    `json:"tokenAddress"`
	TokenID       string    `json:"tokenId"`
	TokenStandard string    `json:"tokenStandard"`
	Balance       string    `json:"balance"`
	UpdatedBlock  uint64    `json:"updatedBlock"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type WebhookEvent struct {
	EventID            string         `json:"eventId"`
	EventType          string         `json:"eventType"`
	ChainID            uint64         `json:"chainId"`
	BlockNumber        uint64         `json:"blockNumber"`
	BlockHash          string         `json:"blockHash"`
	TxHash             string         `json:"txHash"`
	Address            string         `json:"address"`
	ConfirmationStatus string         `json:"confirmationStatus,omitempty"`
	Confirmations      uint64         `json:"confirmations,omitempty"`
	Finalized          uint8          `json:"finalized,omitempty"`
	Removed            uint8          `json:"removed,omitempty"`
	Payload            map[string]any `json:"payload"`
	BlockTime          *time.Time     `json:"blockTime,omitempty"`
	CreatedAt          time.Time      `json:"createdAt"`
}

type Webhook struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	ChainID   *int     `json:"chainId"`
	Secret    string   `json:"secret,omitempty"`
	IsActive  bool     `json:"isActive"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}

type CreateWebhookResponse struct {
	WebhookID string `json:"webhookId"`
	Secret    string `json:"secret"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type listResponse[T any] struct {
	Items []T `json:"items"`
}
