package service_test

import (
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
)

const (
	testBucketPrimary  = "my-bucket"
	testBucketFallback = "fallback-bucket"
	testAssetsAppJS    = "assets/app.js"
	testNestedAppJS    = "my-bucket/assets/app.js"
)

type parseTarget struct {
	bucket string
	key    string
}

func TestParseLocationURIs(t *testing.T) {
	tests := []struct {
		name          string
		location      string
		defaultBucket string
		want          parseTarget
	}{
		{
			name:          "gs URI standard",
			location:      "gs://my-bucket/index.html",
			defaultBucket: "",
			want:          parseTarget{bucket: testBucketPrimary, key: "index.html"},
		},
		{
			name:          "gs URI nested key",
			location:      "gs://my-bucket/assets/js/app.js",
			defaultBucket: testBucketFallback,
			want:          parseTarget{bucket: testBucketPrimary, key: "assets/js/app.js"},
		},
		{
			name:          "s3 URI standard",
			location:      "s3://asset-bucket/style.css",
			defaultBucket: "",
			want:          parseTarget{bucket: "asset-bucket", key: "style.css"},
		},
		{
			name:          "gs URI empty bucket with fallback",
			location:      "gs:///index.html",
			defaultBucket: testBucketFallback,
			want:          parseTarget{bucket: testBucketFallback, key: "index.html"},
		},
		{
			name:          "gs URI empty bucket without fallback",
			location:      "gs:///index.html",
			defaultBucket: "",
			want:          parseTarget{bucket: "", key: "index.html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, k, err := service.ParseLocation(tt.location, tt.defaultBucket)
			assertParsedLocation(t, parseTarget{bucket: b, key: k}, tt.want, err)
		})
	}
}

func assertParsedLocation(t *testing.T, got, want parseTarget, err error) {
	t.Helper()
	checkParseError(t, err)
	checkTarget(t, got, want)
}

func checkParseError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func checkTarget(t *testing.T, got, want parseTarget) {
	t.Helper()
	if got.bucket != want.bucket || got.key != want.key {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
