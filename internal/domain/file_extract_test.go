package domain

import (
	"testing"
)

func TestExtractFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"   \t\n  ", ""},
		{"/", ""},
		{"///", ""},
		{"/media/", ""},
		{"  /media/  ", ""},
		{"/media/ /", ""},
		{"/media/ font.woff2", "font.woff2"},
		{".", ""},
		{"./", ""},
		{"/.", ""},
		{"/./", ""},
		{"..", ""},
		{"/..", ""},
		{"../../", ""},
		{"/media/..", ""},
		{"/nested/path/../..", ""},
		{"index.html", "index.html"},
		{"/index.html", "index.html"},
		{"///app.js", "app.js"},
		{"  style.css  ", "style.css"},
		{"/media/material-icons-LEZCGFVT.woff2", "material-icons-LEZCGFVT.woff2"},
		{"media/material-icons-LEZCGFVT.woff2", "material-icons-LEZCGFVT.woff2"},
		{"/static/css/chunks/theme.css", "theme.css"},
		{"static/css/chunks/theme.css", "theme.css"},
		{"   /media/material-icons-LEZCGFVT.woff2   ", "material-icons-LEZCGFVT.woff2"},
		{"///media///fonts//icon.woff2", "icon.woff2"},
		{"/static/css/../theme.css", "theme.css"},
		{"../secret", "secret"},
	}

	for _, tt := range tests {
		got := ExtractFileName(tt.input)
		if got != tt.expected {
			t.Errorf("ExtractFileName(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
