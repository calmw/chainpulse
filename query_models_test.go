package chainpulse

import (
	"encoding/json"
	"testing"
)

// Query API JSON uses YYYY-MM-DD HH:mm:ss (UTC) via timeformat.FormatValue.
const queryDisplayTime = "2026-06-07 18:32:03"

func TestQueryModelsUnmarshalDisplayTime(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "ChainTransaction",
			body: `{"chainId":968,"blockNumber":12035315,"blockHash":"0xabc","txHash":"0xdef","from":"0x1","to":"0x2","value":"0","nonce":1,"gas":21000,"gasPrice":"1","status":1,"transactionIndex":0,"blockTime":"` + queryDisplayTime + `","finalized":0,"removed":0,"version":1}`,
		},
		{
			name: "InternalTransaction",
			body: `{"chainId":968,"blockNumber":12035315,"blockHash":"0xabc","txHash":"0xdef","traceAddress":"0.0.0","callType":"CALL","from":"0x1","to":"0x2","value":"1000000000000000000","removed":0,"blockTime":"` + queryDisplayTime + `","version":1}`,
		},
		{
			name: "ChainLog",
			body: `{"chainId":968,"blockNumber":12035315,"blockHash":"0xabc","txHash":"0xdef","logIndex":0,"address":"0x1","topic0":"0x2","topic1":"","topic2":"","topic3":"","data":"0x","eventName":"Transfer","removed":0,"blockTime":"` + queryDisplayTime + `","version":1}`,
		},
		{
			name: "AddressBalance",
			body: `{"chainId":968,"address":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723","assetType":"coin","balance":"1000000000000000000","updatedBlock":12035315,"updatedAt":"` + queryDisplayTime + `"}`,
		},
		{
			name: "NFTTokenBalance",
			body: `{"chainId":968,"ownerAddress":"0x1","tokenAddress":"0x2","tokenId":"1","tokenStandard":"erc721","balance":"1","updatedBlock":12035315,"updatedAt":"` + queryDisplayTime + `"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "ChainTransaction":
				var v ChainTransaction
				if err := json.Unmarshal([]byte(tt.body), &v); err != nil {
					t.Fatal(err)
				}
				if v.BlockTime.IsZero() {
					t.Fatal("expected blockTime")
				}
			case "InternalTransaction":
				var v InternalTransaction
				if err := json.Unmarshal([]byte(tt.body), &v); err != nil {
					t.Fatal(err)
				}
				if v.BlockTime.IsZero() {
					t.Fatal("expected blockTime")
				}
			case "ChainLog":
				var v ChainLog
				if err := json.Unmarshal([]byte(tt.body), &v); err != nil {
					t.Fatal(err)
				}
				if v.BlockTime.IsZero() {
					t.Fatal("expected blockTime")
				}
			case "AddressBalance":
				var v AddressBalance
				if err := json.Unmarshal([]byte(tt.body), &v); err != nil {
					t.Fatal(err)
				}
				if v.UpdatedAt.IsZero() {
					t.Fatal("expected updatedAt")
				}
			case "NFTTokenBalance":
				var v NFTTokenBalance
				if err := json.Unmarshal([]byte(tt.body), &v); err != nil {
					t.Fatal(err)
				}
				if v.UpdatedAt.IsZero() {
					t.Fatal("expected updatedAt")
				}
			}
		})
	}
}

func TestListAddressBalancesResponseUnmarshal(t *testing.T) {
	body := []byte(`{"items":[{"chainId":968,"address":"0x0f4b9fc118dc2428745a10970f680ff06b0d5723","assetType":"coin","balance":"1000000000000000000","updatedBlock":12035315,"updatedAt":"2026-06-07 18:32:03"}]}`)
	var out listResponse[AddressBalance]
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("unmarshal list response: %v", err)
	}
	if len(out.Items) != 1 || out.Items[0].UpdatedAt.IsZero() {
		t.Fatalf("unexpected items: %#v", out.Items)
	}
}
