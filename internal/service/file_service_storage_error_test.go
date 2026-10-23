package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func TestGetAssetStorageTimeout(t *testing.T) {
	ctx := context.Background()
	repo := createAppJSRepo()
	store := createErrStorage(context.DeadlineExceeded)
	svc := mustNewService(t, repo, store)

	_, err := svc.GetAsset(ctx, testAppJS)
	assertNonNotFoundStorageError(t, err, context.DeadlineExceeded, "timeout")
}

func TestGetAssetStorageAccessDenied(t *testing.T) {
	ctx := context.Background()
	repo := createAppJSRepo()
	store := createErrStorage(storage.ErrAccessDenied)
	svc := mustNewService(t, repo, store)

	_, err := svc.GetAsset(ctx, testAppJS)
	assertNonNotFoundStorageError(t, err, storage.ErrAccessDenied, "access denied")
}

func TestGetIndexAssetDBError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("mongodb down")
	repo := &mockRepository{
		findFileFunc: func(ctx context.Context, name string) (*domain.SPADocument, error) {
			_ = ctx
			_ = name
			return nil, dbErr
		},
	}
	store := &mockStorage{}
	svc := mustNewService(t, repo, store)

	_, err := svc.GetIndexAsset(ctx)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func createErrStorage(err error) *mockStorage {
	return &mockStorage{
		getObjectFunc: func(ctx context.Context, in storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			_ = ctx
			_ = in
			return nil, err
		},
	}
}

func assertNonNotFoundStorageError(t *testing.T, err, wantErr error, desc string) {
	t.Helper()
	checkErrNonNil(t, err, desc)
	checkErrMatching(t, err, wantErr, desc)
}

func checkErrNonNil(t *testing.T, err error, desc string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error for %s, got nil", desc)
	}
}

func checkErrMatching(t *testing.T, err, wantErr error, desc string) {
	t.Helper()
	if !errors.Is(err, wantErr) {
		t.Errorf("expected %v for %s, got %v", wantErr, desc, err)
	}
	if errors.Is(err, domain.ErrFileNotFound) {
		t.Errorf("did not expect ErrFileNotFound for %s, got %v", desc, err)
	}
}
