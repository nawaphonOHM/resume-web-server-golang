package service_test

import (
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
)

func TestParseLocationErrors(t *testing.T) {
	tests := []struct {
		name          string
		location      string
		defaultBucket string
	}{
		{
			name:          "empty location",
			location:      "",
			defaultBucket: testBucketPrimary,
		},
		{
			name:          "whitespace location",
			location:      "   ",
			defaultBucket: testBucketPrimary,
		},
		{
			name:          "only slash",
			location:      "/",
			defaultBucket: testBucketPrimary,
		},
		{
			name:          "multiple slashes",
			location:      "///",
			defaultBucket: testBucketPrimary,
		},
		{
			name:          "gs scheme empty without fallback",
			location:      "gs://",
			defaultBucket: "",
		},
		{
			name:          "gs URI empty key",
			location:      "gs://my-bucket/",
			defaultBucket: "",
		},
		{
			name:          "gs URI no slash no fallback",
			location:      "gs://my-bucket",
			defaultBucket: "",
		},
		{
			name:          "gs URI no slash with fallback",
			location:      "gs://my-bucket",
			defaultBucket: testBucketFallback,
		},
		{
			name:          "s3 URI no slash",
			location:      "s3://asset-bucket",
			defaultBucket: testBucketFallback,
		},
		{
			name:          "https URL unsupported scheme",
			location:      "https://storage.googleapis.com/bucket/key",
			defaultBucket: testBucketPrimary,
		},
		{
			name:          "http URL unsupported scheme",
			location:      "http://example.com/asset.js",
			defaultBucket: testBucketPrimary,
		},
		{
			name:          "ftp URL unsupported scheme",
			location:      "ftp://example.com/file",
			defaultBucket: testBucketPrimary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.ParseLocation(tt.location, tt.defaultBucket)
			if !errors.Is(err, service.ErrInvalidLocation) {
				t.Errorf("got error %v, want %v", err, service.ErrInvalidLocation)
			}
		})
	}
}
