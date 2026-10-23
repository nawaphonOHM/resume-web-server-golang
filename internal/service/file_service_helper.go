// Package service provides domain business logic for static file resolution.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

func (s *FileStorageService) resolveAsset(
	ctx context.Context,
	cleanName string,
) (*domain.StaticAsset, error) {
	doc, err := s.repo.FindFile(ctx, cleanName)
	if err != nil {
		return s.handleLookupError(ctx, err)
	}
	return s.fetchObject(ctx, doc.Location)
}

func (s *FileStorageService) handleLookupError(
	ctx context.Context,
	err error,
) (*domain.StaticAsset, error) {
	if errors.Is(err, domain.ErrFileNotFound) {
		return s.GetIndexAsset(ctx)
	}
	return nil, err
}

func (s *FileStorageService) fetchObject(
	ctx context.Context,
	location string,
) (*domain.StaticAsset, error) {
	bucket, key, err := ParseLocation(location, s.defaultBucket)
	if err != nil {
		return nil, err
	}
	out, err := s.storage.GetObject(ctx, storage.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return nil, mapStorageError(err, key)
	}
	return toStaticAsset(out), nil
}

func isIndexTarget(name, indexFile string) bool {
	return name == "" || name == indexFile
}

func resolveIndexFile(indexFile string) string {
	if indexFile == "" {
		return config.DefaultIndexFile
	}
	return domain.CleanFileName(indexFile)
}

func mapStorageError(err error, key string) error {
	if errors.Is(err, storage.ErrObjectNotFound) {
		return fmt.Errorf("storage object %q: %w: %w", key, domain.ErrFileNotFound, err)
	}
	return fmt.Errorf("fetch storage object %q failed: %w", key, err)
}

func toStaticAsset(out *storage.GetObjectOutput) *domain.StaticAsset {
	return &domain.StaticAsset{
		LastModified:  out.LastModified,
		Body:          out.Body,
		ContentType:   out.ContentType,
		ETag:          out.ETag,
		ContentLength: out.ContentLength,
	}
}
