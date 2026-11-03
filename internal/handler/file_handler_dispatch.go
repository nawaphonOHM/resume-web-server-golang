// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func logNilAssetError(ctx context.Context, fileName string) {
	logNilAsset(ctx, fileName)
}

func logNilRootAssetError(ctx context.Context) {
	logNilRootAsset(ctx)
}

func dispatchAssetResponse(ctx context.Context, name string, asset *domain.StaticAsset) rest.Response {
	if asset == nil {
		logNilAssetError(ctx, name)
		return rest.InternalServerError(codeInternalError, detailInternal)
	}
	logFilePrepared(ctx, name, asset)
	return buildAssetResponse(asset)
}

func dispatchRootResponse(ctx context.Context, asset *domain.StaticAsset) rest.Response {
	if asset == nil {
		logNilRootAssetError(ctx)
		return rest.InternalServerError(codeInternalError, detailInternal)
	}
	logRootPrepared(ctx, asset)
	return buildAssetResponse(asset)
}
