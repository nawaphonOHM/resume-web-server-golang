// Package config_test tests configuration loading and validation.
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
)

const (
	dotEnvConfigsFile = "configs/.env"
	testDirMode       = 0o750
	configsEnvContent = "OHM9996_SERVER_API_GROUP=standard\nOHM9996_STORAGE_BUCKET=cfg-bucket\n" +
		"OHM9996_INDEX_FILE=cfg.html\n"
)

// chdirTemp switches the working directory to a fresh temp dir for the test.
func chdirTemp(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { restoreWorkdir(t, wd) })
}

func restoreWorkdir(t *testing.T, wd string) {
	t.Helper()
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("restore wd: %v", err)
	}
}

func writeFileAt(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), testDirMode); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), testFileMode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestLoadDefaultReadsConfigsDotEnvWhenDotEnvAbsent(t *testing.T) {
	prepareCleanEnv(t, envKeyServerAPIGroup, envKeyIndex, envKeyStorage)
	chdirTemp(t)
	writeFileAt(t, dotEnvConfigsFile, configsEnvContent)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, groupStandardMode, cfg.ServerAPIGroup)
	assertEqual(t, groupStandardMode, os.Getenv(envKeyServerAPIGroup))
	assertEqual(t, "cfg-bucket", cfg.StorageBucket)
	assertEqual(t, "cfg.html", cfg.IndexFile)
}

func TestLoadDefaultPrefersDotEnvOverConfigsDotEnv(t *testing.T) {
	prepareCleanEnv(t, envKeyServerAPIGroup, envKeyStorage)
	chdirTemp(t)
	writeFileAt(t, ".env", "OHM9996_STORAGE_BUCKET=root-bucket\n")
	writeFileAt(t, dotEnvConfigsFile, "OHM9996_STORAGE_BUCKET=cfg-bucket\nOHM9996_SERVER_API_GROUP=standard\n")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, "root-bucket", cfg.StorageBucket)
	assertEqual(t, config.DefaultServerAPIGroup, cfg.ServerAPIGroup)
}

func TestLoadDefaultAppliesRawGroupWhenNoDotEnvFiles(t *testing.T) {
	prepareCleanEnv(t, envKeyServerAPIGroup)
	chdirTemp(t)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, config.DefaultServerAPIGroup, cfg.ServerAPIGroup)
	assertEqual(t, config.DefaultServerAPIGroup, os.Getenv(envKeyServerAPIGroup))
}

func TestLoadDefaultExistingGroupEnvWinsOverConfigsDotEnv(t *testing.T) {
	prepareCleanEnv(t, envKeyServerAPIGroup)
	t.Setenv(envKeyServerAPIGroup, config.DefaultServerAPIGroup)
	chdirTemp(t)
	writeFileAt(t, dotEnvConfigsFile, "OHM9996_SERVER_API_GROUP=standard\n")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	assertEqual(t, config.DefaultServerAPIGroup, cfg.ServerAPIGroup)
}
