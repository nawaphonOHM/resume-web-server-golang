// Package config provides application configuration loading and validation.
package config

import (
	"fmt"
	"os"
	"strings"
)

func parseEnvLine(line string) error {
	trimmed := strings.TrimSpace(line)
	if isIgnoredLine(trimmed) {
		return nil
	}
	return processEnvPair(trimmed)
}

func isIgnoredLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#")
}

func processEnvPair(line string) error {
	key, val, found := strings.Cut(line, "=")
	if !found {
		return nil
	}
	trimmedKey := strings.TrimSpace(key)
	trimmedVal := trimQuotes(strings.TrimSpace(val))
	return setEnvIfNotSet(trimmedKey, trimmedVal)
}

func isKeySet(key string) bool {
	if key == "" {
		return true
	}
	_, exists := os.LookupEnv(key)
	return exists
}

func setEnvIfNotSet(key, val string) error {
	if isKeySet(key) {
		return nil
	}
	if err := os.Setenv(key, val); err != nil {
		return fmt.Errorf("set env %q: %w", key, err)
	}
	return nil
}

func trimQuotes(val string) string {
	if len(val) >= 2 {
		return stripMatchingQuotePair(val)
	}
	return val
}

func isEnclosedIn(val, quote string) bool {
	return strings.HasPrefix(val, quote) && strings.HasSuffix(val, quote)
}

func stripMatchingQuotePair(val string) string {
	if isEnclosedIn(val, "\"") || isEnclosedIn(val, "'") {
		return val[1 : len(val)-1]
	}
	return val
}
