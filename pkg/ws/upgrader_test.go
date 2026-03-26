package ws

import (
	"net/http"
	"testing"
)

func TestCheckOrigin_NoOriginHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/ws", nil)
	if !checkOrigin(r) {
		t.Error("expected true when no Origin header (same-origin or non-browser)")
	}
}

func TestCheckOrigin_AllowAll(t *testing.T) {
	SetAllowedOrigins([]string{"*"})
	defer SetAllowedOrigins(nil)

	r, _ := http.NewRequest("GET", "/ws", nil)
	r.Header.Set("Origin", "https://evil.com")
	if !checkOrigin(r) {
		t.Error("expected true when * is in allowed origins")
	}
}

func TestCheckOrigin_AllowedOrigin(t *testing.T) {
	SetAllowedOrigins([]string{"https://app.example.com", "https://admin.example.com"})
	defer SetAllowedOrigins(nil)

	r, _ := http.NewRequest("GET", "/ws", nil)
	r.Header.Set("Origin", "https://app.example.com")
	if !checkOrigin(r) {
		t.Error("expected true for allowed origin")
	}
}

func TestCheckOrigin_CaseInsensitive(t *testing.T) {
	SetAllowedOrigins([]string{"https://App.Example.Com"})
	defer SetAllowedOrigins(nil)

	r, _ := http.NewRequest("GET", "/ws", nil)
	r.Header.Set("Origin", "https://app.example.com")
	if !checkOrigin(r) {
		t.Error("expected true for case-insensitive match")
	}
}

func TestCheckOrigin_RejectedOrigin(t *testing.T) {
	SetAllowedOrigins([]string{"https://app.example.com"})
	defer SetAllowedOrigins(nil)

	r, _ := http.NewRequest("GET", "/ws", nil)
	r.Header.Set("Origin", "https://evil.com")
	if checkOrigin(r) {
		t.Error("expected false for disallowed origin")
	}
}

func TestCheckOrigin_EmptyAllowedOrigins(t *testing.T) {
	SetAllowedOrigins(nil)

	r, _ := http.NewRequest("GET", "/ws", nil)
	r.Header.Set("Origin", "https://any.com")
	if checkOrigin(r) {
		t.Error("expected false when no origins configured and Origin header present")
	}
}
