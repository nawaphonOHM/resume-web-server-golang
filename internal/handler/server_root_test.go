// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func indexAsset(body *spyReadCloser) *domain.StaticAsset {
	return &domain.StaticAsset{
		Body:        body,
		ContentType: indexType, ETag: `"idx"`, ContentLength: int64(len(indexBody)),
		LastModified: testModified(),
	}
}

func TestServerServesRootAtSlash(t *testing.T) {
	spy := newSpyReadCloser(strings.NewReader(indexBody))
	base := startServer(t, indexService(spy))

	resp, body := fetch(t, base+"/")

	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, indexType, resp.Header.Get(hdrContent))
	assertEqual(t, `"idx"`, resp.Header.Get(headerETag))
	assertEqual(t, "13", resp.Header.Get(headerLen))
	assertEqual(t, testModTime, resp.Header.Get(headerMod))
	assertEqual(t, indexBody, body)
	spy.waitClosed(t)
}

func indexService(spy *spyReadCloser) *mockFileService {
	return &mockFileService{getIndexAssetFunc: func(context.Context) (*domain.StaticAsset, error) {
		return indexAsset(spy), nil
	}}
}
