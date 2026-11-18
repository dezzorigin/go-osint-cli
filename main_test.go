package main

import (
	"strings"
	"testing"
)

func TestFormatResult(t *testing.T) {
	p := Platform{Name: "GitHub", URL: "https://github.com/%s"}
	url := "https://github.com/dezzorigin"

	tests := []struct {
		name   string
		status int
		want   string
	}{
		{"found", 200, "[GitHub] \033[32mFound\033[0m -> https://github.com/dezzorigin"},
		{"not found", 404, "[GitHub] \033[31mNot Found\033[0m (Available)"},
		{"redirect", 302, "[GitHub] \033[33mRedirect (Status: 302)\033[0m — cannot verify (login wall?)"},
		{"forbidden", 403, "[GitHub] Unexpected Response (Status: 403)"},
		{"server error", 500, "[GitHub] Unexpected Response (Status: 500)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatResult(p, url, tc.status); got != tc.want {
				t.Errorf("formatResult(%d) = %q, want %q", tc.status, got, tc.want)
			}
		})
	}
}

func TestPlatforms(t *testing.T) {
	if len(platforms) != 5 {
		t.Fatalf("expected 5 platforms, got %d", len(platforms))
	}
	for _, p := range platforms {
		if p.Name == "" || !strings.Contains(p.URL, "%s") {
			t.Errorf("platform %q has a bad URL pattern: %q", p.Name, p.URL)
		}
	}
}
