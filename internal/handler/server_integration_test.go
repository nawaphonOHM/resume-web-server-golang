// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

const (
	indexBody  = "<html></html>"
	indexType  = "text/html; charset=utf-8"
	problemCT  = "application/problem+json"
	hdrContent = "Content-Type"
	nameBuffer = 4
)

// recordingService serves a fixed index and records requested file names.
type recordingService struct {
	names   chan string
	spyBody *spyReadCloser
}

func (r *recordingService) getAsset(ctx context.Context, name string) (*domain.StaticAsset, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.names <- name
	r.spyBody = newSpyReadCloser(strings.NewReader(testBody))
	return &domain.StaticAsset{
		Body:        r.spyBody,
		ContentType: testType, ETag: testETag, ContentLength: testLength,
	}, nil
}

func (r *recordingService) service() *mockFileService {
	return &mockFileService{
		getIndexAssetFunc: func(context.Context) (*domain.StaticAsset, error) {
			return indexAsset(newSpyReadCloser(strings.NewReader(indexBody))), nil
		},
		getAssetFunc: r.getAsset,
	}
}

func TestServerServesSingleSegmentFile(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	base := startServer(t, rec.service())

	resp, body := fetch(t, base+"/app.js")
	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, testType, resp.Header.Get(hdrContent))
	assertEqual(t, testETag, resp.Header.Get(headerETag))
	assertEqual(t, testBody, body)
	assertEqual(t, "app.js", <-rec.names)
	rec.spyBody.waitClosed(t)
}

func TestServerNestedPathDoesNotHitAssetHandler(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	base := startServer(t, rec.service())

	resp, _ := fetch(t, base+"/assets/chunks/app.js")
	assertEqual(t, http.StatusNotFound, resp.StatusCode)
	assertEqual(t, 0, len(rec.names))
}

func TestServerServesFileWithoutOptionalMetadata(t *testing.T) {
	rec := &recordingService{names: make(chan string, 1)}
	svc := &mockFileService{
		getAssetFunc: func(context.Context, string) (*domain.StaticAsset, error) {
			rec.spyBody = newSpyReadCloser(strings.NewReader("minimal"))
			return &domain.StaticAsset{Body: rec.spyBody, ContentType: "text/plain"}, nil
		},
	}
	base := startServer(t, svc)

	resp, body := fetch(t, base+"/minimal.txt")
	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, "text/plain", resp.Header.Get(hdrContent))
	assertEqual(t, "", resp.Header.Get(headerETag))
	assertEqual(t, "", resp.Header.Get(headerMod))
	assertEqual(t, "minimal", body)
	rec.spyBody.waitClosed(t)
}

func TestServerTreatsAPIPathAsPlainSPARoute(t *testing.T) {
	rec := &recordingService{names: make(chan string, 1)}
	base := startServer(t, rec.service())

	resp, _ := fetch(t, base+"/api")

	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, "api", <-rec.names)
}

func TestServerKeepsHealthEndpoints(t *testing.T) {
	base := startServer(t, (&recordingService{}).service())

	resp, _ := fetch(t, base+"/health")

	assertEqual(t, http.StatusOK, resp.StatusCode)
}
