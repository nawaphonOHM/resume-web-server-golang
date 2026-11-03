// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

const (
	logKeyRawParam      = "rawParam"
	logKeyFileName      = "fileName"
	logKeyPath          = "path"
	logKeyContentType   = "contentType"
	logKeyContentLength = "contentLength"
	rootPathTarget      = "/"
	pathSlash           = "/"
)

func combineParams(p1, fn string) string {
	p1Clean := strings.TrimPrefix(p1, pathSlash)
	fnClean := strings.TrimPrefix(fn, pathSlash)
	if fnClean == "" {
		return pathSlash + p1Clean + pathSlash
	}
	return pathSlash + p1Clean + pathSlash + fnClean
}

func selectParam(p1, fn string) string {
	if p1 == "" {
		return fn
	}
	if fn == "" {
		return p1
	}
	return combineParams(p1, fn)
}

func extractRawParam(c rest.Context) string {
	if c == nil {
		return ""
	}
	return selectParam(c.Param("p1"), c.Param("fileName"))
}

func resolveFileParams(c rest.Context) (raw, name string, ctx context.Context) {
	raw = extractRawParam(c)
	name = resolveFileNameParam(c)
	ctx = extractRequestContext(c)
	return raw, name, ctx
}

func logFileRequest(ctx context.Context, rawParam, fileName string) {
	logging.DebugContext(ctx, "resolving static file request", logKeyRawParam, rawParam, logKeyFileName, fileName)
	logging.InfoContext(ctx, "serving static file request", logKeyFileName, fileName)
}

func logFilePrepared(ctx context.Context, fileName string, asset *domain.StaticAsset) {
	logging.DebugContext(ctx, "static asset response prepared",
		logKeyFileName, fileName,
		logKeyContentType, asset.ContentType,
		logKeyContentLength, asset.ContentLength,
	)
	logging.InfoContext(ctx, "static asset resolved, streaming response",
		logKeyFileName, fileName,
		logKeyContentType, asset.ContentType,
	)
}

func logRootRequest(ctx context.Context) {
	logging.InfoContext(ctx, "serving root asset request", logKeyPath, rootPathTarget)
}

func logRootPrepared(ctx context.Context, asset *domain.StaticAsset) {
	logging.DebugContext(ctx, "root asset response prepared",
		logKeyPath, rootPathTarget,
		logKeyContentType, asset.ContentType,
		logKeyContentLength, asset.ContentLength,
	)
	logging.InfoContext(ctx, "root asset resolved, streaming response",
		logKeyPath, rootPathTarget,
		logKeyContentType, asset.ContentType,
	)
}

func logNilAsset(ctx context.Context, fileName string) {
	logging.ErrorContext(ctx, "static asset resolution returned nil asset", logKeyFileName, fileName)
}

func logNilRootAsset(ctx context.Context) {
	logging.ErrorContext(ctx, "root asset resolution returned nil asset", logKeyPath, rootPathTarget)
}
