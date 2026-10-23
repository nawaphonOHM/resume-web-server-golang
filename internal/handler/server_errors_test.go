// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func failingIndexService(err error) *mockFileService {
	return &mockFileService{getIndexAssetFunc: func(context.Context) (*domain.StaticAsset, error) {
		return nil, err
	}}
}

func TestServerMissingIndexReturnsProblemDetails(t *testing.T) {
	base := startServer(t, failingIndexService(domain.ErrFileNotFound))

	resp, body := fetch(t, base+"/")

	assertEqual(t, http.StatusNotFound, resp.StatusCode)
	assertEqual(t, problemCT, resp.Header.Get(hdrContent))
	assertEqual(t, true, strings.Contains(body, codeFileNotFound))
}

func TestServerInternalErrorDoesNotLeakDetails(t *testing.T) {
	base := startServer(t, failingIndexService(errBackend))

	resp, body := fetch(t, base+"/")

	assertEqual(t, http.StatusInternalServerError, resp.StatusCode)
	assertEqual(t, problemCT, resp.Header.Get(hdrContent))
	assertEqual(t, true, strings.Contains(body, detailInternal))
	assertEqual(t, false, strings.Contains(body, "secret"))
}

func blockingIndexService(started chan<- struct{}, result chan<- error) *mockFileService {
	return &mockFileService{getIndexAssetFunc: func(ctx context.Context) (*domain.StaticAsset, error) {
		close(started)
		<-ctx.Done()
		result <- ctx.Err()
		return nil, ctx.Err()
	}}
}

func awaitCanceled(t *testing.T, result <-chan error) {
	t.Helper()
	select {
	case got := <-result:
		assertEqual(t, true, errors.Is(got, context.Canceled))
	case <-time.After(waitTimeout):
		t.Fatal("service context was not canceled after client disconnect")
	}
}

func cancelOnStart(started <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-started
		cancel()
	}()
	return ctx, cancel
}

func TestServerPropagatesClientCancellation(t *testing.T) {
	started := make(chan struct{})
	result := make(chan error, 1)
	base := startServer(t, blockingIndexService(started, result))
	ctx, cancel := cancelOnStart(started)
	defer cancel()

	if resp, err := doGet(ctx, base+"/"); err == nil {
		closeBodyOf(t, resp)
	}

	awaitCanceled(t, result)
}
