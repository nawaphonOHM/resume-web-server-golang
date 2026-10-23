// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const (
	loopbackHost   = "127.0.0.1"
	pollInterval   = 20 * time.Millisecond
	envKeyAPIGroup = "OHM9996_SERVER_API_GROUP"
)

func waitReady(t *testing.T, base string, done <-chan error) {
	t.Helper()
	deadline := time.Now().Add(waitTimeout)
	for time.Now().Before(deadline) {
		if resp, err := doGet(context.Background(), base+"/health"); err == nil {
			closeBodyOf(t, resp)
			return
		}
		time.Sleep(pollInterval)
	}
	t.Fatalf("server did not become ready, last start result: %v", done)
}

func awaitShutdown(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		assertEqual(t, true, err == nil)
	case <-time.After(waitTimeout):
		t.Fatal("server did not shut down gracefully")
	}
}

func stopServer(t *testing.T, done <-chan error) {
	t.Helper()
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatalf("send SIGTERM: %v", err)
	}
	awaitShutdown(t, done)
}

func setServerEnv(t *testing.T, port string) {
	t.Helper()
	t.Setenv(envKeyAPIGroup, "")
	if err := os.Unsetenv(envKeyAPIGroup); err != nil {
		t.Fatalf("unset %s: %v", envKeyAPIGroup, err)
	}
	if _, err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}
	t.Setenv("OHM9996_GIN_MODE", "test")
	t.Setenv("OHM9996_SERVER_HOST", loopbackHost)
	t.Setenv("OHM9996_SERVER_PORT", port)
	t.Setenv("OHM9996_SERVER_CORS_PATH", "")
}

// startServer boots the real rest.StartREST server and returns its base URL.
func startServer(t *testing.T, svc domain.FileService) string {
	t.Helper()
	port := freePort(t)
	setServerEnv(t, port)
	done := make(chan error, 1)
	go func() { done <- rest.StartREST(RegisterRoutes(NewFileHandler(svc))) }()
	base := "http://" + loopbackHost + ":" + port
	waitReady(t, base, done)
	t.Cleanup(func() { stopServer(t, done) })
	return base
}
