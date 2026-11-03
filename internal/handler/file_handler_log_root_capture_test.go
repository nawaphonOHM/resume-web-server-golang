package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/testutil/logcapture"
)

func verifyRootSuccessLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	reqInfo := logcapture.RequireEntry(t, entries, logLvlInfo, "serving root asset request")
	logcapture.AssertField(t, reqInfo, logKeyPath, rootPathTarget)
	logcapture.AssertField(t, reqInfo, logKeyFileName, nil)

	prepInfo := logcapture.RequireEntry(t, entries, logLvlInfo, "root asset resolved, streaming response")
	logcapture.AssertField(t, prepInfo, logKeyPath, rootPathTarget)
	logcapture.AssertField(t, prepInfo, logKeyFileName, nil)
}

func TestServeRootLogging(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	h := NewFileHandler(&mockFileService{
		getIndexAssetFunc: func(context.Context) (*domain.StaticAsset, error) {
			return newFullAsset(), nil
		},
	})
	resp := h.ServeRoot(&mockRestContext{})
	assertEqual(t, http.StatusOK, resp.StatusCode())

	entries := stop()
	verifyRootSuccessLogs(t, entries)
}

func TestServeFileLoggingNilAsset(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	h := NewFileHandler(&mockFileService{
		getAssetFunc: func(context.Context, string) (*domain.StaticAsset, error) {
			return nil, nil
		},
	})
	resp := h.ServeFile(&mockRestContext{params: map[string]string{testParamFileName: testAppJS}})
	assertEqual(t, http.StatusInternalServerError, resp.StatusCode())

	entries := stop()
	nilErr := logcapture.RequireEntry(t, entries, logLvlError, "static asset resolution returned nil asset")
	logcapture.AssertField(t, nilErr, logKeyFileName, testAppJS)
}

func TestUnconfiguredHandlerNoLogs(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	h := NewFileHandler(nil)
	respFile := h.ServeFile(&mockRestContext{params: map[string]string{testParamFileName: testAppJS}})
	respRoot := h.ServeRoot(&mockRestContext{})

	assertEqual(t, http.StatusInternalServerError, respFile.StatusCode())
	assertEqual(t, http.StatusInternalServerError, respRoot.StatusCode())

	entries := stop()
	assertEqual(t, 0, len(entries))
}
