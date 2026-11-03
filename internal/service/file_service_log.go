// Package service provides domain business logic for static file resolution.
package service

import (
	"context"

	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	logKeyFileName      = "fileName"
	logKeyCleanName     = "cleanName"
	logKeyIndexFile     = "indexFile"
	logKeyBucket        = "bucket"
	logKeyKey           = "key"
	logKeyLocation      = "location"
	logKeyContentType   = "contentType"
	logKeyContentLength = "contentLength"
)

func logEvaluateTarget(ctx context.Context, fileName, cleanName string) {
	logging.DebugContext(ctx, "evaluating asset target", logKeyFileName, fileName, logKeyCleanName, cleanName)
}

func logTargetIsIndex(ctx context.Context, fileName, indexFile string) {
	logging.DebugContext(ctx, "target evaluated as index entrypoint",
		logKeyFileName, fileName,
		logKeyIndexFile, indexFile,
	)
}

func logQueryMetadata(ctx context.Context, fileName string) {
	logging.DebugContext(ctx, "querying asset metadata", logKeyFileName, fileName)
}

func logMetadataFound(ctx context.Context, fileName, location string) {
	logging.DebugContext(ctx, "asset metadata found", logKeyFileName, fileName, logKeyLocation, location)
}

func logSPAFallback(ctx context.Context, fileName, indexFile string) {
	logging.InfoContext(ctx, "asset not found in metadata, triggering SPA fallback",
		logKeyFileName, fileName,
		logKeyIndexFile, indexFile,
	)
}

func logFetchObject(ctx context.Context, bucket, key string) {
	logging.DebugContext(ctx, "fetching object from storage", logKeyBucket, bucket, logKeyKey, key)
}

func logObjectRetrieved(ctx context.Context, bucket, key string, out *storage.GetObjectOutput) {
	if out == nil {
		return
	}
	logging.DebugContext(ctx, "storage object fetched successfully",
		logKeyBucket, bucket,
		logKeyKey, key,
		logKeyContentType, out.ContentType,
		logKeyContentLength, out.ContentLength,
	)
}
