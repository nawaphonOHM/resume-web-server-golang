// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const (
	codeFileNotFound  = "FILE_NOT_FOUND"
	codeInternalError = "INTERNAL_SERVER_ERROR"
	detailNotFound    = "the requested file was not found"
	detailInternal    = "the server failed to serve the requested file"
)

// FileHandler handles HTTP requests for serving single page application assets.
type FileHandler struct {
	fileService domain.FileService
}

// NewFileHandler initializes a new FileHandler with the given FileService.
func NewFileHandler(fileService domain.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

// isHandlerUnconfigured checks whether handler or its service dependency is nil.
func isHandlerUnconfigured(h *FileHandler) bool {
	if h == nil {
		return true
	}
	return h.fileService == nil
}

// resolveFileNameParam extracts the fileName parameter from the request context.
func resolveFileNameParam(c rest.Context) string {
	if c == nil {
		return ""
	}
	return c.Param("fileName")
}

// mapServiceError translates domain service errors into standardized RFC 9457 responses.
// Internal error details are logged and never exposed to clients.
func mapServiceError(ctx context.Context, err error) rest.Response {
	if errors.Is(err, domain.ErrFileNotFound) {
		return rest.NotFound(codeFileNotFound, detailNotFound)
	}
	logging.ErrorContext(ctx, "failed to serve static asset", "err", err)
	return rest.InternalServerError(codeInternalError, detailInternal)
}

// buildAssetResponse builds a streaming REST response forwarding asset headers and body.
func buildAssetResponse(asset *domain.StaticAsset) rest.Response {
	if asset == nil {
		return rest.InternalServerError(codeInternalError, detailInternal)
	}
	resp := rest.Stream(http.StatusOK, asset.Body)
	resp = applyOptionalHeader(resp, "Content-Type", asset.ContentType)
	resp = applyOptionalIntHeader(resp, "Content-Length", asset.ContentLength)
	resp = applyOptionalHeader(resp, "ETag", asset.ETag)
	resp = applyOptionalStringerHeader(resp, "Last-Modified", asset.LastModified)
	return resp
}

// ServeRoot handles GET / by resolving and serving the index.html asset.
func (h *FileHandler) ServeRoot(c rest.Context) rest.Response {
	if isHandlerUnconfigured(h) {
		return rest.InternalServerError("SERVICE_UNAVAILABLE", "file service not configured")
	}
	ctx := extractRequestContext(c)
	asset, err := h.fileService.GetIndexAsset(ctx)
	if err != nil {
		return mapServiceError(ctx, err)
	}
	return buildAssetResponse(asset)
}

// ServeFile handles GET /:fileName by resolving the asset or falling back to index.html.
func (h *FileHandler) ServeFile(c rest.Context) rest.Response {
	if isHandlerUnconfigured(h) {
		return rest.InternalServerError("SERVICE_UNAVAILABLE", "file service not configured")
	}
	fileName := resolveFileNameParam(c)
	ctx := extractRequestContext(c)
	asset, err := h.fileService.GetAsset(ctx, fileName)
	if err != nil {
		return mapServiceError(ctx, err)
	}
	return buildAssetResponse(asset)
}
