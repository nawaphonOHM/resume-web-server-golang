// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func TestServerServesNestedMediaFile(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	base := startServer(t, rec.service())

	resp, body := fetch(t, base+"/media/material-icons-LEZCGFVT.woff2")
	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, testType, resp.Header.Get(hdrContent))
	assertEqual(t, testETag, resp.Header.Get(headerETag))
	assertEqual(t, testBody, body)
	assertEqual(t, "material-icons-LEZCGFVT.woff2", <-rec.names)
	rec.spyBody.waitClosed(t)
}

func TestServerServesDeeplyNestedStaticFile(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	base := startServer(t, rec.service())

	resp, body := fetch(t, base+"/assets/chunks/app.js")
	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, testType, resp.Header.Get(hdrContent))
	assertEqual(t, testETag, resp.Header.Get(headerETag))
	assertEqual(t, testBody, body)
	assertEqual(t, "app.js", <-rec.names)
	rec.spyBody.waitClosed(t)
}

func TestServerNestedPathMissingFileReturns404(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	svc := &mockFileService{
		getAssetFunc: func(ctx context.Context, name string) (*domain.StaticAsset, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			rec.names <- name
			return nil, domain.ErrFileNotFound
		},
	}
	base := startServer(t, svc)

	resp, body := fetch(t, base+"/media/missing-icon.woff2")
	assertEqual(t, http.StatusNotFound, resp.StatusCode)
	assertEqual(t, problemCT, resp.Header.Get(hdrContent))
	assertEqual(t, true, strings.Contains(body, codeFileNotFound))
	assertEqual(t, "missing-icon.woff2", <-rec.names)
}

func TestServerKeepsReadyEndpoint(t *testing.T) {
	base := startServer(t, (&recordingService{}).service())

	resp, _ := fetch(t, base+"/ready")
	assertEqual(t, http.StatusOK, resp.StatusCode)
}

func TestServerHandlesProbeSubpathAsSPA(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	base := startServer(t, rec.service())

	resp, _ := fetch(t, base+"/health/x")
	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, "x", <-rec.names)
}

func TestServerTrailingSlashFallsBack(t *testing.T) {
	rec := &recordingService{names: make(chan string, nameBuffer)}
	base := startServer(t, rec.service())

	resp, _ := fetch(t, base+"/media/")
	assertEqual(t, http.StatusOK, resp.StatusCode)
	assertEqual(t, "", <-rec.names)

	respSingle, _ := fetch(t, base+"/app.js/")
	assertEqual(t, http.StatusOK, respSingle.StatusCode)
	assertEqual(t, "", <-rec.names)
}
