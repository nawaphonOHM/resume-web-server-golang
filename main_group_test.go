package main

import (
	"errors"
	"os"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	envKeyServerGroup = "OHM9996_SERVER_API_GROUP"
	groupStandard     = "standard"
)

func checkGroupEnvAndConfig(t *testing.T, cfg *config.Config) {
	t.Helper()
	if cfg.ServerAPIGroup != config.DefaultServerAPIGroup {
		t.Fatalf("expected ServerAPIGroup %s, got %s", config.DefaultServerAPIGroup, cfg.ServerAPIGroup)
	}
	if os.Getenv(envKeyServerGroup) != config.DefaultServerAPIGroup {
		t.Fatalf("expected env %s, got %s", config.DefaultServerAPIGroup, os.Getenv(envKeyServerGroup))
	}
}

func TestDefaultAPIGroupMatchesRoutes(t *testing.T) {
	t.Setenv(envKeyServerGroup, "")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}
	checkGroupEnvAndConfig(t, cfg)
}

func loadStandardConfig(t *testing.T) *config.Config {
	t.Helper()
	t.Setenv(envKeyServerGroup, groupStandard)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}
	return cfg
}

func TestExplicitStandardAPIGroupMismatchesRoutes(t *testing.T) {
	cfg := loadStandardConfig(t)
	bp, err := buildBluePrint(cfg, &mongodb.Client{}, &storage.Client{})
	if err != nil {
		t.Fatalf("expected blueprint creation to succeed, got %v", err)
	}
	err = rest.StartREST(bp)
	if !errors.Is(err, rest.ErrAPIGroupMismatch) {
		t.Fatalf("expected ErrAPIGroupMismatch, got %v", err)
	}
}
