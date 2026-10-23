// Package config_test tests environment file loading logic.
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
)

func prepareCleanEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset env %q: %v", key, err)
		}
	}
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env.test")
	if err := os.WriteFile(envPath, []byte(content), testFileMode); err != nil {
		t.Fatalf("failed to write test env file: %v", err)
	}
	return envPath
}

func TestLoadWithDotEnvFile(t *testing.T) {
	prepareCleanEnv(t, envKeyIndex, envKeyStorage)
	content := "# comment line\n\nOHM9996_INDEX_FILE=\"test.html\"\nOHM9996_STORAGE_BUCKET='test-bucket'\n"
	envPath := writeTempEnv(t, content)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, "test.html", cfg.IndexFile)
	assertEqual(t, "test-bucket", cfg.StorageBucket)
}

func TestLoadEnvExistingEnvPrecedence(t *testing.T) {
	t.Setenv(envKeyStorage, "existing-bucket")
	content := "OHM9996_STORAGE_BUCKET=dotenv-bucket\n"
	envPath := writeTempEnv(t, content)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, "existing-bucket", cfg.StorageBucket)
}

func TestLoadEnvMismatchedQuotes(t *testing.T) {
	prepareCleanEnv(t, envKeyIndex, envKeyStorage)
	content := "OHM9996_INDEX_FILE=\"mismatched'\nOHM9996_STORAGE_BUCKET='other\"\n"
	envPath := writeTempEnv(t, content)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, "\"mismatched'", cfg.IndexFile)
	assertEqual(t, "'other\"", cfg.StorageBucket)
}

func TestLoadEnvMalformedLines(t *testing.T) {
	prepareCleanEnv(t, envKeyStorage)
	content := "INVALID_LINE_WITHOUT_EQUALS\n# Comment\n\nOHM9996_STORAGE_BUCKET=valid-bucket\n"
	envPath := writeTempEnv(t, content)

	cfg, err := config.Load(envPath)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, "valid-bucket", cfg.StorageBucket)
}

func TestLoadExplicitMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/file/path")
	if err == nil {
		t.Error("expected error loading nonexistent explicit file, got nil")
	}
}

func TestLoadDefaultMissingFile(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected nil error when default .env is missing, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}
