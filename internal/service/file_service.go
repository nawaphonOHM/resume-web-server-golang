// Package service provides domain business logic for static file resolution.
package service

import (
	"context"
	"fmt"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

// FileStorageService retrieves static assets with SPA fallback and cloud storage streaming.
type FileStorageService struct {
	repo          domain.FileRepository
	storage       domain.ObjectStorage
	defaultBucket string
	indexFile     string
}

// NewFileStorageService creates a new FileStorageService instance.
func NewFileStorageService(
	repo domain.FileRepository,
	store domain.ObjectStorage,
	defaultBucket, indexFile string,
) (*FileStorageService, error) {
	if repo == nil {
		return nil, ErrNilRepository
	}
	if store == nil {
		return nil, ErrNilStorage
	}
	return &FileStorageService{
		repo:          repo,
		storage:       store,
		defaultBucket: defaultBucket,
		indexFile:     resolveIndexFile(indexFile),
	}, nil
}

// GetAsset resolves and retrieves a static asset by file name.
// If the requested file is not found in MongoDB metadata, it falls back to index.html.
func (s *FileStorageService) GetAsset(ctx context.Context, fileName string) (*domain.StaticAsset, error) {
	cleanName := domain.CleanFileName(fileName)
	logEvaluateTarget(ctx, fileName, cleanName)
	if isIndexTarget(cleanName, s.indexFile) {
		logTargetIsIndex(ctx, fileName, s.indexFile)
		return s.GetIndexAsset(ctx)
	}
	return s.resolveAsset(ctx, cleanName)
}

// GetIndexAsset resolves and retrieves the SPA entrypoint index.html asset.
func (s *FileStorageService) GetIndexAsset(ctx context.Context) (*domain.StaticAsset, error) {
	logQueryMetadata(ctx, s.indexFile)
	doc, err := s.repo.FindFile(ctx, s.indexFile)
	if err != nil {
		return nil, fmt.Errorf("lookup index file %q failed: %w", s.indexFile, err)
	}
	logMetadataFound(ctx, s.indexFile, doc.Location)
	return s.fetchObject(ctx, doc.Location)
}
