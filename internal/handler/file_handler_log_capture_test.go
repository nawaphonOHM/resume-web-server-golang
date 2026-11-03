package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/testutil/logcapture"
)

const (
	logLvlDebug = "DEBUG"
	logLvlInfo  = "INFO"
	logLvlError = "ERROR"
)

func newFullAsset() *domain.StaticAsset {
	return &domain.StaticAsset{
		Body:          newSpyReadCloser(strings.NewReader(testBody)),
		ContentType:   testType,
		ContentLength: testLength,
	}
}

func verifyFileRequestLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	rawDebug := logcapture.RequireEntry(t, entries, logLvlDebug, "resolving static file request")
	logcapture.AssertField(t, rawDebug, logKeyRawParam, testNestedWoff2)
	logcapture.AssertField(t, rawDebug, logKeyFileName, testWoff2)

	reqInfo := logcapture.RequireEntry(t, entries, logLvlInfo, "serving static file request")
	logcapture.AssertField(t, reqInfo, logKeyFileName, testWoff2)
}

func verifyFileResponseLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	prepInfo := logcapture.RequireEntry(t, entries, logLvlInfo, "static asset resolved, streaming response")
	logcapture.AssertField(t, prepInfo, logKeyFileName, testWoff2)

	prepDebug := logcapture.RequireEntry(t, entries, logLvlDebug, "static asset response prepared")
	logcapture.AssertField(t, prepDebug, logKeyContentLength, float64(testLength))
}

func TestServeFileLoggingSuccess(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	h := NewFileHandler(&mockFileService{
		getAssetFunc: func(context.Context, string) (*domain.StaticAsset, error) {
			return newFullAsset(), nil
		},
	})
	resp := h.ServeFile(&mockRestContext{params: map[string]string{
		"p1":              "media",
		testParamFileName: "/" + testWoff2,
	}})
	assertEqual(t, http.StatusOK, resp.StatusCode())

	entries := stop()
	verifyFileRequestLogs(t, entries)
	verifyFileResponseLogs(t, entries)
}

func TestServeFileLoggingErrors(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	h := NewFileHandler(&mockFileService{
		getAssetFunc: func(context.Context, string) (*domain.StaticAsset, error) {
			return nil, errBackend
		},
	})
	resp := h.ServeFile(&mockRestContext{params: map[string]string{testParamFileName: testAppJS}})
	assertEqual(t, http.StatusInternalServerError, resp.StatusCode())

	entries := stop()
	logcapture.AssertNoEntry(t, entries, logLvlInfo, "static asset resolved, streaming response")
	logcapture.RequireEntry(t, entries, logLvlError, "failed to serve static asset")
}
