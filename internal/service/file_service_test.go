package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
)

const (
	testIndexHTML = "index.html"
	testAppJS     = "app.js"
	testHTMLType  = "text/html"
	testJSType    = "application/javascript"
)

func TestNewFileStorageServiceNilRepo(t *testing.T) {
	store := &mockStorage{}
	_, err := service.NewFileStorageService(nil, store, testBucketPrimary, "")
	if !errors.Is(err, service.ErrNilRepository) {
		t.Fatalf("expected ErrNilRepository, got %v", err)
	}
}

func TestNewFileStorageServiceNilStorage(t *testing.T) {
	repo := &mockRepository{}
	_, err := service.NewFileStorageService(repo, nil, testBucketPrimary, "")
	if !errors.Is(err, service.ErrNilStorage) {
		t.Fatalf("expected ErrNilStorage, got %v", err)
	}
}

func TestNewFileStorageServiceSuccess(t *testing.T) {
	repo := &mockRepository{}
	store := &mockStorage{}
	svc, err := service.NewFileStorageService(repo, store, testBucketPrimary, "")
	if err != nil || svc == nil {
		t.Fatalf("unexpected error creating service: %v", err)
	}
}

func TestGetIndexAssetSuccess(t *testing.T) {
	ctx := context.Background()
	content := "<html>hello</html>"
	svc := mustNewService(t, createIndexRepo(), createIndexStorage(content))

	asset, err := svc.GetIndexAsset(ctx)
	if err != nil {
		t.Fatalf("GetIndexAsset failed: %v", err)
	}
	defer closeBody(t, asset.Body)

	verifyStaticAsset(t, asset, assetExpectation{
		lastModified:  fixedTestTime,
		content:       content,
		contentType:   testHTMLType,
		etag:          `"etag-index"`,
		contentLength: int64(len(content)),
	})
}

func createIndexRepo() *mockRepository {
	return &mockRepository{
		findFileFunc: func(ctx context.Context, fileName string) (*domain.SPADocument, error) {
			_ = ctx
			if fileName == testIndexHTML {
				return &domain.SPADocument{
					FileName: testIndexHTML,
					Location: "gs://my-bucket/index.html",
				}, nil
			}
			return nil, domain.ErrFileNotFound
		},
	}
}

func createIndexStorage(content string) *mockStorage {
	return &mockStorage{
		getObjectFunc: func(ctx context.Context, in storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			_ = ctx
			if in.Bucket == testBucketPrimary && in.Key == testIndexHTML {
				return &storage.GetObjectOutput{
					Body:          io.NopCloser(bytes.NewBufferString(content)),
					ContentType:   testHTMLType,
					ETag:          `"etag-index"`,
					ContentLength: int64(len(content)),
					LastModified:  fixedTestTime,
				}, nil
			}
			return nil, storage.ErrObjectNotFound
		},
	}
}
