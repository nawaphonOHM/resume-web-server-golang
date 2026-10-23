package main

import (
	"context"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestDisconnectClientNil(t *testing.T) {
	t.Parallel()
	disconnectClient(context.Background(), nil)
}

func TestBuildBluePrintSuccess(t *testing.T) {
	t.Parallel()
	cfg := config.DefaultConfig()
	bp, err := buildBluePrint(cfg, &mongodb.Client{}, &storage.Client{})
	if err != nil {
		t.Fatalf("expected blueprint creation to succeed, got %v", err)
	}
	if bp == nil {
		t.Fatal("expected non-nil blueprint")
	}
}

func TestBuildBluePrintNilStorage(t *testing.T) {
	t.Parallel()
	cfg := config.DefaultConfig()
	bp, err := buildBluePrint(cfg, &mongodb.Client{}, nil)
	if err == nil {
		t.Fatal("expected error when storage client is nil, got nil")
	}
	if bp != nil {
		t.Fatalf("expected nil blueprint on error, got %v", bp)
	}
}

func TestRunConfigError(t *testing.T) {
	t.Setenv("OHM9996_INDEX_FILE", "   ")
	err := run(context.Background())
	if err == nil {
		t.Fatal("expected error on invalid config, got nil")
	}
}

func TestRunMongoConnectError(t *testing.T) {
	restore := mongodb.SetExitFunc(func(code int) {
		if code != 0 {
			return
		}
	})
	defer func() { mongodb.SetExitFunc(restore) }()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Setenv("OHM9996_INDEX_FILE", "index.html")
	t.Setenv("OHM9996_MONGODB_HOST", "invalid-mongo-host.local")
	err := run(ctx)
	if err == nil {
		t.Fatal("expected error on canceled context mongo connection, got nil")
	}
}

func TestStartServerStorageError(t *testing.T) {
	restore := storage.SetExitFunc(func(code int) {
		if code != 0 {
			return
		}
	})
	defer func() { storage.SetExitFunc(restore) }()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := config.DefaultConfig()
	t.Setenv("OHM9996_STORAGE_BUCKET", "")
	err := startServer(ctx, cfg, nil)
	if err == nil {
		t.Fatal("expected error on storage connect failure, got nil")
	}
}
