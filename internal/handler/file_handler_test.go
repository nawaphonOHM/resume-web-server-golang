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

	resp := h.ServeFile(&mockRestContext{params: map[string]string{"p1": "app.js"}})

	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "app.js", svc.name)
}

func TestServeFileNestedPathPassesExtractedFileName(t *testing.T) {
	svc := &nameCapturingService{}
	h := NewFileHandler(&mockFileService{getAssetFunc: svc.getAsset})

	resp := h.ServeFile(&mockRestContext{params: map[string]string{
		"p1":       "media",
		"fileName": "/material-icons-LEZCGFVT.woff2",
	}})

	assertEqual(t, http.StatusOK, resp.StatusCode())
	assertEqual(t, "material-icons-LEZCGFVT.woff2", svc.name)
}

func TestExtractRequestContextFallsBack(t *testing.T) {
	assertEqual(t, context.Background(), extractRequestContext(nil))
	assertEqual(t, context.Background(), extractRequestContext(&mockRestContext{}))
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
