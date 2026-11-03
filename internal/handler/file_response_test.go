// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

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
