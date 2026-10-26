// Package config_test tests configuration loading and validation.
package config_test

import (
	"os"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
)

const (
	envKeyServerAPIGroup = "OHM9996_SERVER_API_GROUP"
	groupStandardMode    = "standard"
)

func TestLoadSetsDefaultServerAPIGroupWhenUnset(t *testing.T) {
	t.Setenv(envKeyServerAPIGroup, "")
	envPath := createEmptyEnv(t)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, config.DefaultServerAPIGroup, cfg.ServerAPIGroup)
	assertEqual(t, config.DefaultServerAPIGroup, os.Getenv(envKeyServerAPIGroup))
}

func TestLoadPreservesExplicitServerAPIGroup(t *testing.T) {
	t.Setenv(envKeyServerAPIGroup, groupStandardMode)
	envPath := createEmptyEnv(t)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, groupStandardMode, cfg.ServerAPIGroup)
	assertEqual(t, groupStandardMode, os.Getenv(envKeyServerAPIGroup))
}
