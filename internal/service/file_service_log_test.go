package service_test

import (
	"context"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/testutil/logcapture"
)

const (
	logLvlDebug     = "DEBUG"
	logLvlInfo      = "INFO"
	logKeyFileName  = "fileName"
	logKeyCleanName = "cleanName"
	logKeyIndexFile = "indexFile"
	logKeyBucket    = "bucket"
	logKeyKey       = "key"
	logKeyLocation  = "location"
)

func verifyServiceTargetLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	evalLog := logcapture.RequireEntry(t, entries, logLvlDebug, "evaluating asset target")
	logcapture.AssertField(t, evalLog, logKeyFileName, testIndexHTML)
	logcapture.AssertField(t, evalLog, logKeyCleanName, testIndexHTML)

	metaLog := logcapture.RequireEntry(t, entries, logLvlDebug, "asset metadata found")
	logcapture.AssertField(t, metaLog, logKeyLocation, "gs://my-bucket/index.html")
}

func verifyServiceStorageLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	storeLog := logcapture.RequireEntry(t, entries, logLvlDebug, "storage object fetched successfully")
	logcapture.AssertField(t, storeLog, logKeyBucket, testBucketPrimary)
	logcapture.AssertField(t, storeLog, logKeyKey, testIndexHTML)

	logcapture.AssertNoEntry(t, entries, logLvlInfo, "storage object retrieved")
}

func verifyServiceFallbackLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	fbLog := logcapture.RequireEntry(t, entries, logLvlInfo, "asset not found in metadata, triggering SPA fallback")
	logcapture.AssertField(t, fbLog, logKeyFileName, "missing/route.js")
	logcapture.AssertField(t, fbLog, logKeyIndexFile, testIndexHTML)
	logcapture.AssertNoEntry(t, entries, logLvlDebug, "initiating index fallback retrieval")
}

func TestServiceLoggingSuccess(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	svc := mustNewService(t, createIndexRepo(), createIndexStorage("<html></html>"))

	asset, err := svc.GetAsset(context.Background(), testIndexHTML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer closeBody(t, asset.Body)

	entries := stop()
	verifyServiceTargetLogs(t, entries)
	verifyServiceStorageLogs(t, entries)
}

func TestServiceLoggingSPAFallback(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	svc := mustNewService(t, createIndexRepo(), createIndexStorage("<html></html>"))

	asset, err := svc.GetAsset(context.Background(), "/missing/route.js")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer closeBody(t, asset.Body)

	entries := stop()
	verifyServiceFallbackLogs(t, entries)
}

func TestServiceLoggingStorageError(t *testing.T) {
	stop := logcapture.Capture(t, "debug")
	store := &mockStorage{
		getObjectFunc: func(context.Context, storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			return nil, storage.ErrObjectNotFound
		},
	}
	svc, err := service.NewFileStorageService(createIndexRepo(), store, testBucketPrimary, "")
	if err != nil {
		t.Fatalf("unexpected error creating service: %v", err)
	}
	if _, err := svc.GetAsset(context.Background(), ""); err == nil {
		t.Fatal("expected storage error")
	}

	entries := stop()
	logcapture.RequireEntry(t, entries, logLvlDebug, "fetching object from storage")
	logcapture.AssertNoEntry(t, entries, logLvlDebug, "storage object fetched successfully")
}
