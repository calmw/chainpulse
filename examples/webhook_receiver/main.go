package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	chainpulse "github.com/calmw/chainpulse"
)

func main() {
	secret := os.Getenv("CHAINPULSE_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("set CHAINPULSE_WEBHOOK_SECRET")
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body failed", http.StatusBadRequest)
			return
		}

		ok := chainpulse.VerifyWebhookSignature(
			secret,
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

		fmt.Printf("event=%s chain=%d tx=%s\n", event.EventType, event.ChainID, event.TxHash)
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
