package domain

import (
	"testing"
)

func TestCleanFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"index.html", "index.html"},
		{"/index.html", "index.html"},
		{"///app.js", "app.js"},
		{"  style.css  ", "style.css"},
		{"", ""},
		{"///", ""},
	}

	for _, tt := range tests {
		got := CleanFileName(tt.input)
		if got != tt.expected {
			t.Errorf("CleanFileName(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
