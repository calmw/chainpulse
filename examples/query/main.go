package main

import (
	"context"
	"fmt"
	"log"
	"os"

	chainpulse "github.com/calmw/chainpulse"
)

func main() {
	queryURL := env("CHAINPULSE_QUERY_URL", "https://query.chainpulse.cc")
	apiKey := os.Getenv("CHAINPULSE_API_KEY")
	apiSecret := os.Getenv("CHAINPULSE_API_SECRET")
	txHash := os.Getenv("CHAINPULSE_TX_HASH")
	if apiKey == "" || apiSecret == "" || txHash == "" {
		log.Fatal("set CHAINPULSE_API_KEY, CHAINPULSE_API_SECRET, and CHAINPULSE_TX_HASH")
	}

	client := chainpulse.NewClient(
		chainpulse.WithQueryBaseURL(queryURL),
		chainpulse.WithAPIKey(apiKey, apiSecret),
	)

	tx, err := client.GetTransaction(context.Background(), txHash, chainpulse.QueryParams{
		ChainID:     56,
		Consistency: "stable",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx=%s block=%d from=%s to=%s\n", tx.TxHash, tx.BlockNumber, tx.From, tx.To)
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
