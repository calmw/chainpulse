// Package chainpulse provides a Go SDK for ChainPulse Query API,
// webhook management, and webhook signature verification.
//
// Query API and Webhook payloads use timestamps in YYYY-MM-DD HH:mm:ss (UTC).
// Models use FlexTime (also accepts RFC3339). Published v1.0.0/v1.0.1 used time.Time and failed to decode.
package chainpulse
