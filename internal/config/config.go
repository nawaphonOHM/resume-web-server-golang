// Package config provides application configuration loading and validation.
package config

import (
	"errors"
	"os"
	"strings"
)

// Default constants for application configuration.
const (
	DefaultMongoCollection = "resume-single-page-application"
	DefaultIndexFile       = "index.html"
	DefaultServerAPIGroup  = "raw"
)

// Sentinel errors for configuration.
var (
	ErrEmptyMongoCollection = errors.New("mongodb collection cannot be empty")
	ErrEmptyIndexFile       = errors.New("index file cannot be empty")
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	MongoCollection string
	IndexFile       string
	StorageBucket   string
	ServerAPIGroup  string
}

// DefaultConfig returns the default application configuration.
func DefaultConfig() *Config {
	return &Config{
		MongoCollection: DefaultMongoCollection,
		IndexFile:       DefaultIndexFile,
		StorageBucket:   "",
		ServerAPIGroup:  DefaultServerAPIGroup,
	}
}

func initEnv(filenames ...string) error {
	if err := loadDotEnvFiles(filenames...); err != nil {
		return err
	}
	return ensureDefaultServerAPIGroup()
}

// Load loads application configuration from environment variables and optional .env files.
func Load(filenames ...string) (*Config, error) {
	if err := initEnv(filenames...); err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	applyEnvOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration fields.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.MongoCollection) == "" {
		return ErrEmptyMongoCollection
	}
	if strings.TrimSpace(c.IndexFile) == "" {
		return ErrEmptyIndexFile
	}
	return nil
}

func applyEnvOverrides(cfg *Config) {
	applyEnvString(&cfg.MongoCollection, "OHM9996_MONGODB_COLLECTION")
	applyEnvString(&cfg.IndexFile, "OHM9996_INDEX_FILE")
	applyEnvString(&cfg.StorageBucket, "OHM9996_STORAGE_BUCKET")
	applyEnvString(&cfg.ServerAPIGroup, "OHM9996_SERVER_API_GROUP")
}

func applyEnvString(target *string, key string) {
	if val := os.Getenv(key); val != "" {
		*target = val
	}
}
