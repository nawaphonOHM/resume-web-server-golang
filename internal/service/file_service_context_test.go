package service_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

type testContextKey string

const (
	traceIDKey testContextKey = "trace-id"
	traceIDVal string         = "trace-abc-123"
)

func TestGetAssetContextPropagation(t *testing.T) {
	ctx := context.WithValue(context.Background(), traceIDKey, traceIDVal)
	var repoCtxVal, storeCtxVal any

	repo, store := createContextComponents(&repoCtxVal, &storeCtxVal)
	svc := mustNewService(t, repo, store)

	asset, err := svc.GetAsset(ctx, testAppJS)
	if err != nil {
		t.Fatalf("GetAsset failed: %v", err)
	}
	defer closeBody(t, asset.Body)

	assertContextValue(t, repoCtxVal, "repository")
	assertContextValue(t, storeCtxVal, "storage")
}

func createContextComponents(repoVal, storeVal *any) (*mockRepository, *mockStorage) {
	repo := &mockRepository{
		findFileFunc: func(c context.Context, fileName string) (*domain.SPADocument, error) {
			*repoVal = c.Value(traceIDKey)
			return &domain.SPADocument{
				FileName: fileName,
				Location: "gs://my-bucket/app.js",
			}, nil
		},
	}
	store := &mockStorage{
		getObjectFunc: func(c context.Context, in storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			_ = in
			*storeVal = c.Value(traceIDKey)
			return &storage.GetObjectOutput{
				Body:          io.NopCloser(bytes.NewBufferString("ok")),
				ContentType:   testJSType,
				ETag:          `"etag"`,
				ContentLength: 2,
				LastModified:  fixedTestTime,
			}, nil
		},
	}
	return repo, store
}

func assertContextValue(t *testing.T, got any, target string) {
	t.Helper()
	if got != traceIDVal {
		t.Errorf("expected %s to receive context value %q, got %v", target, traceIDVal, got)
	}
}
