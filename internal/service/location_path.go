// Package service provides domain business logic for static file resolution.
package service

import (
	"strings"
)

func parseURILocation(loc, defaultBucket string) (bucket, key string, err error) {
	idx := strings.Index(loc, schemeDelimiter)
	scheme := strings.ToLower(loc[:idx])
	if !isSupportedScheme(scheme) {
		return "", "", ErrInvalidLocation
	}
	rawTarget := loc[idx+schemeOffsetBytes:]
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
	slashIdx := strings.Index(rawTarget, slashSeparator)
	if slashIdx == -1 {
		return "", "", ErrInvalidLocation
	}
	b := rawTarget[:slashIdx]
	k := strings.TrimLeft(rawTarget[slashIdx+1:], slashSeparator)
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
