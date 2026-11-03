// Package logcapture provides test helpers for capturing and asserting structured log output.
package logcapture

import (
	"testing"
)

func matchEntry(entry map[string]any, level, msg string) bool {
	if entry["level"] != level {
		return false
	}
	return entry["msg"] == msg
}

func findEntry(entries []map[string]any, level, msg string) (map[string]any, bool) {
	for _, entry := range entries {
		if matchEntry(entry, level, msg) {
			return entry, true
		}
	}
	return nil, false
}

// RequireEntry asserts that an entry with the given level and message is present.
func RequireEntry(t *testing.T, entries []map[string]any, level, msg string) map[string]any {
	t.Helper()
	entry, ok := findEntry(entries, level, msg)
	if !ok {
		t.Fatalf("expected log entry [%s] %q not found", level, msg)
	}
	return entry
}

// AssertNoEntry asserts that an entry with the given level and message is not present.
func AssertNoEntry(t *testing.T, entries []map[string]any, level, msg string) {
	t.Helper()
	if _, ok := findEntry(entries, level, msg); ok {
		t.Errorf("unexpected log entry [%s] %q was found", level, msg)
	}
}

// AssertField asserts that the entry contains key with the expected value.
func AssertField(t *testing.T, entry map[string]any, key string, want any) {
	t.Helper()
	got := entry[key]
	if got != want {
		t.Errorf("field %q: expected %v, got %v", key, want, got)
	}
}
