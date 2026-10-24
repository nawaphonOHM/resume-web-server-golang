package service_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func TestGetAssetSPAFallback(t *testing.T) {
	ctx := context.Background()
	indexContent := "<html>spa entrypoint</html>"
	var findCalls []string
	var getCalls []storage.GetObjectInput

	repo, store := createFallbackSpies(indexContent, &findCalls, &getCalls)
	svc := mustNewService(t, repo, store)

	fetchAndVerifyFallbackAsset(t, ctx, svc, indexContent)
	assertFallbackFindCalls(t, findCalls)
	assertFallbackGetCalls(t, getCalls)
}

func fetchAndVerifyFallbackAsset(t *testing.T, ctx context.Context, svc domain.FileService, content string) {
	t.Helper()
	asset, err := svc.GetAsset(ctx, "/dashboard/analytics")
	if err != nil {
		t.Fatalf("GetAsset fallback failed: %v", err)
	}
	defer closeBody(t, asset.Body)
	verifyFallbackAsset(t, asset, content)
}

func createFallbackSpies(
	content string,
	findCalls *[]string,
	getCalls *[]storage.GetObjectInput,
) (*mockRepository, *mockStorage) {
	repo := &mockRepository{
		findFileFunc: func(ctx context.Context, fileName string) (*domain.SPADocument, error) {
			_ = ctx
			*findCalls = append(*findCalls, fileName)
			if fileName == testIndexHTML {
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
				Body:          io.NopCloser(bytes.NewBufferString(content)),
				ContentType:   testHTMLType,
				ETag:          `"etag-index"`,
				ContentLength: int64(len(content)),
				LastModified:  fixedTestTime,
			}, nil
		},
	}
	return repo, store
}

func verifyFallbackAsset(t *testing.T, asset *domain.StaticAsset, content string) {
	t.Helper()
	verifyStaticAsset(t, asset, assetExpectation{
		lastModified:  fixedTestTime,
		content:       content,
		contentType:   testHTMLType,
		etag:          `"etag-index"`,
		contentLength: int64(len(content)),
	})
}

func assertFallbackFindCalls(t *testing.T, findCalls []string) {
	t.Helper()
	if len(findCalls) != 2 {
		t.Fatalf("got %d find calls, want 2", len(findCalls))
	}
	verifyFallbackCallIndices(t, findCalls)
}

func verifyFallbackCallIndices(t *testing.T, findCalls []string) {
	t.Helper()
	if findCalls[0] != "dashboard/analytics" {
		t.Errorf("first find call got %q, want %q", findCalls[0], "dashboard/analytics")
	}
	if findCalls[1] != testIndexHTML {
		t.Errorf("second find call got %q, want %q", findCalls[1], testIndexHTML)
	}
}

func assertFallbackGetCalls(t *testing.T, getCalls []storage.GetObjectInput) {
	t.Helper()
	if len(getCalls) != 1 {
		t.Fatalf("got %d storage calls, want 1", len(getCalls))
	}
	assertStorageTarget(t, getCalls[0], testBucketPrimary, testIndexHTML)
}
