// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const (
	waitTimeout = 5 * time.Second
	testBody    = "console.log(1)"
	testType    = "application/javascript"
	testETag    = `"abc123"`
	testLength  = int64(len(testBody))
	headerType  = "Content-Type"
	headerLen   = "Content-Length"
	headerETag  = "ETag"
	headerMod   = "Last-Modified"
	testModTime = "Mon, 02 Jan 2006 15:04:05 GMT"
	testModUnix = 1136214245
)

func testModified() time.Time {
	return time.Unix(testModUnix, 0).UTC()
}

// spyReadCloser records close calls and wraps an io.Reader.
type spyReadCloser struct {
	reader io.Reader
	done   chan struct{}
	once   sync.Once
}

// newSpyReadCloser creates a new spyReadCloser wrapping reader.
func newSpyReadCloser(r io.Reader) *spyReadCloser {
	return &spyReadCloser{reader: r, done: make(chan struct{})}
}

// Read delegates read operations to the underlying reader.
func (s *spyReadCloser) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

// Close marks the closer as closed.
func (s *spyReadCloser) Close() error {
	s.once.Do(func() { close(s.done) })
	return nil
}

// waitClosed fails the test unless Close was called in time.
func (s *spyReadCloser) waitClosed(t *testing.T) {
	t.Helper()
	select {
	case <-s.done:
	case <-time.After(waitTimeout):
		t.Fatal("expected asset body to be closed")
	}
}

// mockFileService implements domain.FileService for unit testing.
type mockFileService struct {
	getAssetFunc      func(context.Context, string) (*domain.StaticAsset, error)
	getIndexAssetFunc func(context.Context) (*domain.StaticAsset, error)
}

// GetAsset delegates asset resolution to getAssetFunc.
func (m *mockFileService) GetAsset(ctx context.Context, fileName string) (*domain.StaticAsset, error) {
	if m.getAssetFunc != nil {
		return m.getAssetFunc(ctx, fileName)
	}
	return nil, domain.ErrFileNotFound
}

// GetIndexAsset delegates index resolution to getIndexAssetFunc.
func (m *mockFileService) GetIndexAsset(ctx context.Context) (*domain.StaticAsset, error) {
	if m.getIndexAssetFunc != nil {
		return m.getIndexAssetFunc(ctx)
	}
	return nil, domain.ErrFileNotFound
}

type contextEmbedder struct {
	rest.Context
}

// mockRestContext implements rest.Context for testing handler parameters and context.
type mockRestContext struct {
	contextEmbedder
	params map[string]string
	ctx    context.Context
}

// Param retrieves the named route parameter from the test map.
func (m *mockRestContext) Param(key string) string {
	return m.params[key]
}

// Context returns the attached context.Context.
func (m *mockRestContext) Context() context.Context {
	return m.ctx
}

func assertEqual[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}
