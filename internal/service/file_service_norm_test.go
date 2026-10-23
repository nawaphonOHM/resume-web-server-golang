package service_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

const (
	testRootContent = "root"
)

func TestGetAssetRootNormalization(t *testing.T) {
	ctx := context.Background()
	for _, name := range []string{"", "/", "/index.html", "index.html"} {
		t.Run("path_"+name, func(t *testing.T) {
			assertSingleIndexResolution(t, ctx, name)
		})
	}
}

func assertSingleIndexResolution(t *testing.T, ctx context.Context, name string) {
	t.Helper()
	var findCalls []string
	var getCalls []storage.GetObjectInput

	repo, store := createSpyIndexComponents(&findCalls, &getCalls)
	svc := mustNewService(t, repo, store)

	fetchAndVerifyIndexAsset(t, ctx, svc, name)
	assertIndexFindCalls(t, findCalls)
	assertIndexGetCalls(t, getCalls)
}

func fetchAndVerifyIndexAsset(t *testing.T, ctx context.Context, svc domain.FileService, name string) {
	t.Helper()
	asset, err := svc.GetAsset(ctx, name)
	if err != nil {
		t.Fatalf("failed for path %q: %v", name, err)
	}
	defer closeBody(t, asset.Body)
	verifyIndexAssetResult(t, asset)
}

func verifyIndexAssetResult(t *testing.T, asset *domain.StaticAsset) {
	t.Helper()
	verifyStaticAsset(t, asset, assetExpectation{
		lastModified:  fixedTestTime,
		content:       testRootContent,
		contentType:   testHTMLType,
		etag:          `"etag-index"`,
		contentLength: int64(len(testRootContent)),
	})
}

func createSpyIndexComponents(
	findCalls *[]string,
	getCalls *[]storage.GetObjectInput,
) (*mockRepository, *mockStorage) {
	repo := &mockRepository{
		findFileFunc: func(ctx context.Context, name string) (*domain.SPADocument, error) {
			_ = ctx
			*findCalls = append(*findCalls, name)
			if name == testIndexHTML {
				return &domain.SPADocument{
					FileName: testIndexHTML,
					Location: "gs://my-bucket/index.html",
				}, nil
			}
			return nil, domain.ErrFileNotFound
		},
	}
	store := &mockStorage{
		getObjectFunc: func(ctx context.Context, in storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			_ = ctx
			*getCalls = append(*getCalls, in)
			return &storage.GetObjectOutput{
				Body:          io.NopCloser(bytes.NewBufferString(testRootContent)),
				ContentType:   testHTMLType,
				ETag:          `"etag-index"`,
				ContentLength: int64(len(testRootContent)),
				LastModified:  fixedTestTime,
			}, nil
		},
	}
	return repo, store
}

func assertIndexFindCalls(t *testing.T, findCalls []string) {
	t.Helper()
	if len(findCalls) != 1 || findCalls[0] != testIndexHTML {
		t.Errorf("got find calls %v, want ['index.html']", findCalls)
	}
}

func assertIndexGetCalls(t *testing.T, getCalls []storage.GetObjectInput) {
	t.Helper()
	if len(getCalls) != 1 {
		t.Fatalf("got %d storage calls, want 1", len(getCalls))
	}
	assertStorageTarget(t, getCalls[0], testBucketPrimary, testIndexHTML)
}
