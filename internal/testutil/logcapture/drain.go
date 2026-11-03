// Package logcapture provides test helpers for capturing and asserting structured log output.
package logcapture

import (
	"encoding/json"
	"os"
)

func decodeNextEntry(dec *json.Decoder) (map[string]any, bool) {
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return nil, false
	}
	return m, true
}

func collectJSONEntries(dec *json.Decoder) []map[string]any {
	var entries []map[string]any
	for dec.More() {
		if m, ok := decodeNextEntry(dec); ok {
			entries = append(entries, m)
		}
	}
	return entries
}

func readJSONStream(r *os.File) []map[string]any {
	entries := collectJSONEntries(json.NewDecoder(r))
	if err := r.Close(); err != nil {
		return entries
	}
	return entries
}

func startLogReader(r *os.File) chan []map[string]any {
	readDone := make(chan []map[string]any, 1)
	go func() {
		readDone <- readJSONStream(r)
	}()
	return readDone
}
