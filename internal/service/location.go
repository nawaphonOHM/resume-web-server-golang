// Package service provides domain business logic for static file resolution.
package service

import (
	"errors"
	"strings"
)

const (
	slashSeparator    = "/"
	schemeDelimiter   = "://"
	schemeOffsetBytes = len(schemeDelimiter)
	schemeGS          = "gs"
	schemeS3          = "s3"
)

// Sentinel errors for service operations.
var (
	ErrInvalidLocation = errors.New("invalid storage location")
	ErrNilRepository   = errors.New("repository cannot be nil")
	ErrNilStorage      = errors.New("storage client cannot be nil")
)

// ParseLocation parses a document location string into a target bucket and object key.
// It supports cloud storage URIs (gs://bucket/key, s3://bucket/key, gs:///key) and
// non-URI object keys (such as assets/app.js or /index.html). Non-URI paths are
// always treated as object keys in defaultBucket.
func ParseLocation(location, defaultBucket string) (bucket, key string, err error) {
	loc := strings.TrimSpace(location)
	if loc == "" {
		return "", "", ErrInvalidLocation
	}
	if hasScheme(loc) {
		return parseURILocation(loc, defaultBucket)
	}
	return parseKeyLocation(loc, defaultBucket)
}

func hasScheme(loc string) bool {
	idx := strings.Index(loc, schemeDelimiter)
	if idx <= 0 {
		return false
	}
	return !strings.Contains(loc[:idx], slashSeparator)
}

func parseKeyLocation(loc, defaultBucket string) (bucket, key string, err error) {
	cleanedKey := strings.TrimLeft(loc, slashSeparator)
	if cleanedKey == "" {
		return "", "", ErrInvalidLocation
	}
	return defaultBucket, cleanedKey, nil
}
