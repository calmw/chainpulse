// Package chainpulse provides a Go SDK for ChainPulse Query API,
// webhook management, and webhook signature verification.
//
// Webhook POST bodies use createdAt/blockTime in YYYY-MM-DD HH:mm:ss (UTC).
// ParseWebhookEvent relies on FlexTime; use v1.0.1 or later (v1.0.0 used time.Time and failed to parse).
package chainpulse
