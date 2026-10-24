// Package config_test tests configuration loading and validation.
package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
)

const (
	testFileMode     = 0600
	envKeyStorage    = "OHM9996_STORAGE_BUCKET"
	envKeyIndex      = "OHM9996_INDEX_FILE"
	envKeyCollection = "OHM9996_MONGODB_COLLECTION"
	errUnexpectedFmt = "unexpected error: %v"
)

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	assertEqual(t, config.DefaultMongoCollection, cfg.MongoCollection)
	assertEqual(t, config.DefaultIndexFile, cfg.IndexFile)
	assertEqual(t, "", cfg.StorageBucket)
	assertEqual(t, config.DefaultServerAPIGroup, cfg.ServerAPIGroup)
}

func TestValidate(t *testing.T) {
	cfg := config.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got %v", err)
	}
}

func TestValidateEmptyFields(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.MongoCollection = "  "
	if err := cfg.Validate(); !errors.Is(err, config.ErrEmptyMongoCollection) {
		t.Errorf("expected ErrEmptyMongoCollection, got %v", err)
	}

	cfg = config.DefaultConfig()
	cfg.IndexFile = ""
	if err := cfg.Validate(); !errors.Is(err, config.ErrEmptyIndexFile) {
		t.Errorf("expected ErrEmptyIndexFile, got %v", err)
	}
}

func createEmptyEnv(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.empty")
	if err := os.WriteFile(envPath, []byte(""), testFileMode); err != nil {
		t.Fatalf("failed to write empty env file: %v", err)
	}
	return envPath
}

func TestLoadWithEnvOverrides(t *testing.T) {
	setTestEnv(t)
	envPath := createEmptyEnv(t)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, "custom-collection", cfg.MongoCollection)
	assertEqual(t, "main.html", cfg.IndexFile)
	assertEqual(t, "my-bucket", cfg.StorageBucket)
}

func setTestEnv(t *testing.T) {
	t.Helper()
	t.Setenv(envKeyCollection, "custom-collection")
	t.Setenv(envKeyIndex, "main.html")
	t.Setenv(envKeyStorage, "my-bucket")
}
