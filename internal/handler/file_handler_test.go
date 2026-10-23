// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

type ctxKey string

var errBackend = errors.New("mongodb://user:secret@db:27017 unreachable")

func newAsset() *domain.StaticAsset {
	return &domain.StaticAsset{Body: newSpyReadCloser(strings.NewReader(testBody))}
}

type nameCapturingService struct {
	name string
}

func (n *nameCapturingService) getAsset(ctx context.Context, name string) (*domain.StaticAsset, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	n.name = name
	return newAsset(), nil
}

func TestServeRootUsesRequestContext(t *testing.T) {
	var got context.Context
	h := NewFileHandler(&mockFileService{getIndexAssetFunc: func(ctx context.Context) (*domain.StaticAsset, error) {
		got = ctx
		return newAsset(), nil
	}})
	want := context.WithValue(context.Background(), ctxKey("k"), "v")

	resp := h.ServeRoot(&mockRestContext{ctx: want})

	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, want, got)
}

func TestServeFileUsesParamAndContext(t *testing.T) {
	svc := &nameCapturingService{}
	h := NewFileHandler(&mockFileService{getAssetFunc: svc.getAsset})

	resp := h.ServeFile(&mockRestContext{params: map[string]string{"fileName": "assets", "rest": "/app.js"}})

	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "assets/app.js", svc.name)
}

func TestExtractRequestContextFallsBack(t *testing.T) {
	assertEqual(t, context.Background(), extractRequestContext(nil))
	assertEqual(t, context.Background(), extractRequestContext(&mockRestContext{}))
}

func TestMapServiceErrorStatuses(t *testing.T) {
	ctx := context.Background()
	assertEqual(t, http.StatusNotFound, mapServiceError(ctx, domain.ErrFileNotFound).StatusCode())
	assertEqual(t, http.StatusInternalServerError, mapServiceError(ctx, errBackend).StatusCode())
}

func TestUnconfiguredHandlerReturns500(t *testing.T) {
	var nilHandler *FileHandler
	for _, h := range []*FileHandler{nilHandler, NewFileHandler(nil)} {
		assertEqual(t, http.StatusInternalServerError, h.ServeRoot(&mockRestContext{}).StatusCode())
		assertEqual(t, http.StatusInternalServerError, h.ServeFile(&mockRestContext{}).StatusCode())
	}
}

func TestBuildAssetResponse(t *testing.T) {
	t.Run("nil asset returns internal error", func(t *testing.T) {
		resp := buildAssetResponse(nil)
		assertEqual(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("valid asset with metadata returns stream response", func(t *testing.T) {
		body := newSpyReadCloser(strings.NewReader(testBody))
		asset := &domain.StaticAsset{
			Body:          body,
			ContentType:   testType,
			ETag:          testETag,
			ContentLength: testLength,
			LastModified:  testModified(),
		}
		resp := buildAssetResponse(asset)
		assertEqual(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("valid asset with empty metadata returns stream response", func(t *testing.T) {
		body := newSpyReadCloser(strings.NewReader(""))
		asset := &domain.StaticAsset{Body: body}
		resp := buildAssetResponse(asset)
		assertEqual(t, http.StatusOK, resp.StatusCode())
	})
}

func TestServeRootServiceError(t *testing.T) {
	h := NewFileHandler(&mockFileService{
		getIndexAssetFunc: func(context.Context) (*domain.StaticAsset, error) {
			return nil, domain.ErrFileNotFound
		},
	})
	resp := h.ServeRoot(&mockRestContext{})
	assertEqual(t, http.StatusNotFound, resp.StatusCode())
}

func TestServeFileServiceError(t *testing.T) {
	h := NewFileHandler(&mockFileService{
		getAssetFunc: func(context.Context, string) (*domain.StaticAsset, error) {
			return nil, domain.ErrFileNotFound
		},
	})
	resp := h.ServeFile(&mockRestContext{params: map[string]string{"fileName": "missing.js"}})
	assertEqual(t, http.StatusNotFound, resp.StatusCode())
}
