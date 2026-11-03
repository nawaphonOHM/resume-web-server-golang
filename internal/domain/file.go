// Package domain defines the core domain models and interfaces for the resume web server.
package domain

import (
	"context"
	"errors"
	"io"
	"path"
	"strings"
	"time"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

// ErrFileNotFound Sentinel errors for domain operations.
var (
	ErrFileNotFound = errors.New("file not found")
)

// CleanFileName normalizes a file path by trimming whitespace and leading slashes.
// For example, "/index.html" becomes "index.html" and "///app.js" becomes "app.js".
func CleanFileName(fileName string) string {
	trimmed := strings.TrimSpace(fileName)
	return strings.TrimLeft(trimmed, "/")
}

// ExtractFileName extracts and normalizes the base file name from a path parameter.
// For example, "/media/material-icons-LEZCGFVT.woff2" becomes "material-icons-LEZCGFVT.woff2",
// and "/static/css/chunks/theme.css" becomes "theme.css".
// Trailing slash directory paths (e.g. "/media/"), whitespace-only segments, and
// terminal traversal tokens (".", "..", "/") safely evaluate to an empty string for SPA fallback.
func ExtractFileName(rawPath string) string {
	trimmed := strings.TrimSpace(rawPath)
	if isInvalidPath(trimmed) {
		return ""
	}
	base := strings.TrimSpace(path.Base(path.Clean(trimmed)))
	return sanitizeBase(base)
}

func isInvalidPath(p string) bool {
	return p == "" || strings.HasSuffix(p, "/")
}

func sanitizeBase(base string) string {
	switch base {
	case "", ".", "..", "/":
		return ""
	default:
		return base
	}
}

// SPADocument represents the metadata document stored in MongoDB.
// Documents in the resume-single-page-application collection follow the schema:
//
//	{
//	    "fileName": "index.html",
//	    "location": "gs://bucket-name/index.html",
//	    "hashValueFile": "<hash>"
//	}
//
// The fileName field is the primary lookup key and must be stored without leading slashes.
type SPADocument struct {
	ID            any    `bson:"_id,omitempty"`
	HashValueFile string `bson:"hashValueFile"`
	Location      string `bson:"location"`
	FileName      string `bson:"fileName"`
}

// StaticAsset holds the streamed object content and metadata from storage.
type StaticAsset struct {
	LastModified  time.Time
	Body          io.ReadCloser
	ContentType   string
	ETag          string
	ContentLength int64
}

// FileRepository defines operations for querying SPA file metadata.
type FileRepository interface {
	// FindFile resolves the SPA metadata document for the given file name.
	FindFile(context.Context, string) (*SPADocument, error)
}

// ObjectStorage defines operations for fetching object streams and metadata from cloud storage.
type ObjectStorage interface {
	// GetObject downloads an object stream and its metadata.
	GetObject(context.Context, storage.GetObjectInput) (*storage.GetObjectOutput, error)
}

// FileService defines operations for resolving and retrieving static SPA assets.
type FileService interface {
	// GetAsset retrieves a static asset by file name, falling back to index.html if not found.
	GetAsset(context.Context, string) (*StaticAsset, error)
	// GetIndexAsset retrieves the SPA entrypoint index.html asset.
	GetIndexAsset(context.Context) (*StaticAsset, error)
}
