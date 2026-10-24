package service_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func TestGetAssetExactMatch(t *testing.T) {
	ctx := context.Background()
	content := "console.log('hello');"
	var recordedInput storage.GetObjectInput

	repo := createAppJSRepo()
	store := createAppJSStore(content, &recordedInput)
	svc := mustNewService(t, repo, store)

	fetchAndVerifyExactAsset(t, ctx, svc, content)
	assertStorageTarget(t, recordedInput, testBucketPrimary, testAppJS)
}

func fetchAndVerifyExactAsset(t *testing.T, ctx context.Context, svc domain.FileService, content string) {
	t.Helper()
	asset, err := svc.GetAsset(ctx, "/app.js")
	if err != nil {
		t.Fatalf("GetAsset failed: %v", err)
	}
	defer closeBody(t, asset.Body)
	verifyExactAssetResult(t, asset, content)
}

func createAppJSStore(content string, recordedInput *storage.GetObjectInput) *mockStorage {
	return &mockStorage{
		getObjectFunc: func(ctx context.Context, in storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			_ = ctx
			*recordedInput = in
			return &storage.GetObjectOutput{
				Body:          io.NopCloser(bytes.NewBufferString(content)),
				ContentType:   testJSType,
				ETag:          `"etag-app"`,
				ContentLength: int64(len(content)),
				LastModified:  fixedTestTime,
			}, nil
		},
	}
}

func verifyExactAssetResult(t *testing.T, asset *domain.StaticAsset, content string) {
	t.Helper()
	verifyStaticAsset(t, asset, assetExpectation{
		lastModified:  fixedTestTime,
		content:       content,
		contentType:   testJSType,
		etag:          `"etag-app"`,
		contentLength: int64(len(content)),
	})
}

func createAppJSRepo() *mockRepository {
	return &mockRepository{
		findFileFunc: func(ctx context.Context, fileName string) (*domain.SPADocument, error) {
			_ = ctx
			if fileName == testAppJS {
				return &domain.SPADocument{
					FileName: testAppJS,
					Location: "gs://my-bucket/app.js",
				}, nil
			}
			return nil, domain.ErrFileNotFound
		},
	}
}
