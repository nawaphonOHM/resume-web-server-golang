// Package config provides application configuration loading and validation.
package config

import (
	"fmt"
	"os"
	"strings"
)

const envKeyServerAPIGroup = "OHM9996_SERVER_API_GROUP"

func shouldSetDefaultAPIGroup() bool {
	val, ok := os.LookupEnv(envKeyServerAPIGroup)
	if !ok {
		return true
	}
	return strings.TrimSpace(val) == ""
}

func ensureDefaultServerAPIGroup() error {
	if !shouldSetDefaultAPIGroup() {
		return nil
	}
	if err := os.Setenv(envKeyServerAPIGroup, DefaultServerAPIGroup); err != nil {
		return fmt.Errorf("set default api group env: %w", err)
	}
	return nil
}
