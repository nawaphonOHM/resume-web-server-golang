// Package service provides domain business logic for static file resolution.
package service

import (
	"strings"
)

func parseURILocation(loc, defaultBucket string) (bucket, key string, err error) {
	before, after, _ := strings.Cut(loc, schemeDelimiter)
	scheme := strings.ToLower(before)
	if !isSupportedScheme(scheme) {
		return "", "", ErrInvalidLocation
	}
	rawTarget := after
	return parseURITarget(rawTarget, defaultBucket)
}

func isSupportedScheme(scheme string) bool {
	return scheme == schemeGS || scheme == schemeS3
}

func parseURITarget(rawTarget, defaultBucket string) (bucket, key string, err error) {
	if rawTarget == "" {
		return "", "", ErrInvalidLocation
	}
	if strings.HasPrefix(rawTarget, slashSeparator) {
		return parseURIRootKey(rawTarget, defaultBucket)
	}
	return parseURIBucketAndKey(rawTarget, defaultBucket)
}

func parseURIRootKey(rawTarget, defaultBucket string) (bucket, key string, err error) {
	keyName := strings.TrimLeft(rawTarget, slashSeparator)
	if keyName == "" {
		return "", "", ErrInvalidLocation
	}
	return defaultBucket, keyName, nil
}

func parseURIBucketAndKey(rawTarget, defaultBucket string) (bucket, key string, err error) {
	before, after, ok := strings.Cut(rawTarget, slashSeparator)
	if !ok {
		return "", "", ErrInvalidLocation
	}
	b := before
	k := strings.TrimLeft(after, slashSeparator)
	return validateURIBucketAndKey(b, k, defaultBucket)
}

func validateURIBucketAndKey(b, k, defaultBucket string) (bucket, key string, err error) {
	if k == "" {
		return "", "", ErrInvalidLocation
	}
	if b == "" {
		return defaultBucket, k, nil
	}
	return b, k, nil
}
