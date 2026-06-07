package chainpulse

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexTimeUnmarshalDisplayLayout(t *testing.T) {
	var ft FlexTime
	if err := json.Unmarshal([]byte(`"2026-06-08 02:12:25"`), &ft); err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 6, 8, 2, 12, 25, 0, time.UTC)
	if !ft.Time.Equal(want) {
		t.Fatalf("got %v want %v", ft.Time, want)
	}
}

func TestFlexTimeUnmarshalRFC3339(t *testing.T) {
	var ft FlexTime
	if err := json.Unmarshal([]byte(`"2026-06-08T02:12:25Z"`), &ft); err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 6, 8, 2, 12, 25, 0, time.UTC)
	if !ft.Time.Equal(want) {
		t.Fatalf("got %v want %v", ft.Time, want)
	}
}

func TestFlexTimeUnmarshalNull(t *testing.T) {
	var ft FlexTime
	if err := json.Unmarshal([]byte(`null`), &ft); err != nil {
		t.Fatal(err)
	}
	if !ft.Time.IsZero() {
		t.Fatal("expected zero time")
	}
}
