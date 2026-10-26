package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
)

func TestGetAssetMissingIndexFile(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		findFileFunc: func(ctx context.Context, name string) (*domain.SPADocument, error) {
			_ = ctx
			_ = name
			return nil, domain.ErrFileNotFound
		},
	}
	store := &mockStorage{}
	svc := mustNewService(t, repo, store)

	_, err := svc.GetAsset(ctx, "/unknown")
	if !errors.Is(err, domain.ErrFileNotFound) {
		t.Fatalf("expected ErrFileNotFound, got %v", err)
	}
}

func TestGetAssetStorageNotFound(t *testing.T) {
	ctx := context.Background()
	repo := createAppJSRepo()
	store := &mockStorage{
		getObjectFunc: func(ctx context.Context, in storage.GetObjectInput) (*storage.GetObjectOutput, error) {
			_ = ctx
			_ = in
			return nil, storage.ErrObjectNotFound
		},
	}
	svc := mustNewService(t, repo, store)

	_, err := svc.GetAsset(ctx, testAppJS)
	if !errors.Is(err, domain.ErrFileNotFound) {
		t.Fatalf("expected ErrFileNotFound, got %v", err)
	}
	if !errors.Is(err, storage.ErrObjectNotFound) {
		t.Fatalf("expected storage.ErrObjectNotFound cause, got %v", err)
	}
}

func TestGetAssetDatabaseError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db connection pool exhausted")
	repo := &mockRepository{
		findFileFunc: func(ctx context.Context, name string) (*domain.SPADocument, error) {
			_ = ctx
			_ = name
			return nil, dbErr
		},
	}
	store := &mockStorage{}
	svc := mustNewService(t, repo, store)

	_, err := svc.GetAsset(ctx, testAppJS)
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestGetAssetInvalidLocation(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{
		findFileFunc: func(ctx context.Context, name string) (*domain.SPADocument, error) {
			_ = ctx
			_ = name
			return &domain.SPADocument{
				FileName: testAppJS,
				Location: "",
			}, nil
		},
	}
	store := &mockStorage{}
	svc := mustNewService(t, repo, store)

	_, err := svc.GetAsset(ctx, testAppJS)
	if !errors.Is(err, service.ErrInvalidLocation) {
		t.Fatalf("expected ErrInvalidLocation, got %v", err)
	}
}
