package service_test

import (
	"context"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

type mockRepository struct {
	findFileFunc func(context.Context, string) (*domain.SPADocument, error)
}

func (m *mockRepository) FindFile(ctx context.Context, fileName string) (*domain.SPADocument, error) {
	if m.findFileFunc != nil {
		return m.findFileFunc(ctx, fileName)
	}
	return nil, domain.ErrFileNotFound
}

type mockStorage struct {
	getObjectFunc func(context.Context, storage.GetObjectInput) (*storage.GetObjectOutput, error)
}

func (m *mockStorage) GetObject(
	ctx context.Context,
	in storage.GetObjectInput,
) (*storage.GetObjectOutput, error) {
	if m.getObjectFunc != nil {
		return m.getObjectFunc(ctx, in)
	}
	return nil, storage.ErrObjectNotFound
}
