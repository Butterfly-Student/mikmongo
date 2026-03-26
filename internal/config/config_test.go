package config

import "testing"

func TestParseOrigins_Wildcard(t *testing.T) {
	origins := parseOrigins("*")
	if len(origins) != 1 || origins[0] != "*" {
		t.Errorf("expected [*], got %v", origins)
	}
}

func TestParseOrigins_Multiple(t *testing.T) {
	origins := parseOrigins("https://app.example.com, https://admin.example.com")
	if len(origins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(origins))
	}
	if origins[0] != "https://app.example.com" {
		t.Errorf("expected first origin to be trimmed, got %q", origins[0])
	}
	if origins[1] != "https://admin.example.com" {
		t.Errorf("expected second origin to be trimmed, got %q", origins[1])
	}
}

func TestParseOrigins_Empty(t *testing.T) {
	origins := parseOrigins("")
	if origins != nil {
		t.Errorf("expected nil for empty string, got %v", origins)
	}
}

func TestParseOrigins_SpacesOnly(t *testing.T) {
	origins := parseOrigins(",  ,  ,")
	if len(origins) != 0 {
		t.Errorf("expected empty slice for whitespace-only, got %v", origins)
	}
}
