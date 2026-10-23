package service_test

import (
	"io"
	"testing"
	"time"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
)

var (
	fixedTestTime = time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
)

type assetExpectation struct {
	lastModified  time.Time
	content       string
	contentType   string
	etag          string
	contentLength int64
}

func mustNewService(
	t *testing.T,
	repo domain.FileRepository,
	store domain.ObjectStorage,
) *service.FileStorageService {
	t.Helper()
	svc, err := service.NewFileStorageService(repo, store, testBucketPrimary, testIndexHTML)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	return svc
}

func closeBody(t *testing.T, body io.Closer) {
	t.Helper()
	if err := body.Close(); err != nil {
		t.Errorf("close body failed: %v", err)
	}
}

func verifyStaticAsset(t *testing.T, asset *domain.StaticAsset, exp assetExpectation) {
	t.Helper()
	verifyAssetContent(t, asset.Body, exp.content)
	verifyAssetHeaders(t, asset, exp)
}

func verifyAssetContent(t *testing.T, body io.Reader, wantContent string) {
	t.Helper()
	b, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if string(b) != wantContent {
		t.Fatalf("content mismatch: got %q, want %q", string(b), wantContent)
	}
}

func verifyAssetHeaders(t *testing.T, asset *domain.StaticAsset, exp assetExpectation) {
	t.Helper()
	verifyStringHeaders(t, asset, exp)
	verifyMetadataValues(t, asset, exp)
}

func verifyStringHeaders(t *testing.T, asset *domain.StaticAsset, exp assetExpectation) {
	t.Helper()
	if asset.ContentType != exp.contentType {
		t.Errorf("content-type mismatch: got %q, want %q", asset.ContentType, exp.contentType)
	}
	if asset.ETag != exp.etag {
		t.Errorf("etag mismatch: got %q, want %q", asset.ETag, exp.etag)
	}
}

func verifyMetadataValues(t *testing.T, asset *domain.StaticAsset, exp assetExpectation) {
	t.Helper()
	if asset.ContentLength != exp.contentLength {
		t.Errorf("content-length mismatch: got %d, want %d", asset.ContentLength, exp.contentLength)
	}
	if !asset.LastModified.Equal(exp.lastModified) {
		t.Errorf("last-modified mismatch: got %v, want %v", asset.LastModified, exp.lastModified)
	}
}

func assertStorageTarget(t *testing.T, in storage.GetObjectInput, wantBucket, wantKey string) {
	t.Helper()
	if in.Bucket != wantBucket {
		t.Errorf("got bucket %q, want %q", in.Bucket, wantBucket)
	}
	if in.Key != wantKey {
		t.Errorf("got key %q, want %q", in.Key, wantKey)
	}
}
