package service_test

import (
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
)

func TestParseLocationKeys(t *testing.T) {
	tests := []struct {
		name          string
		location      string
		defaultBucket string
		want          parseTarget
	}{
		{
			name:          "relative path with default bucket",
			location:      "assets/index.html",
			defaultBucket: testBucketPrimary,
			want:          parseTarget{bucket: testBucketPrimary, key: "assets/index.html"},
		},
		{
			name:          "leading slash path with default bucket",
			location:      "/assets/app.js",
			defaultBucket: testBucketPrimary,
			want:          parseTarget{bucket: testBucketPrimary, key: testAssetsAppJS},
		},
		{
			name:          "relative path without default bucket",
			location:      testAssetsAppJS,
			defaultBucket: "",
			want:          parseTarget{bucket: "", key: testAssetsAppJS},
		},
		{
			name:          "leading slash path without default bucket",
			location:      "/assets/app.js",
			defaultBucket: "",
			want:          parseTarget{bucket: "", key: testAssetsAppJS},
		},
		{
			name:          "nested key path without default bucket",
			location:      testNestedAppJS,
			defaultBucket: "",
			want:          parseTarget{bucket: "", key: testNestedAppJS},
		},
		{
			name:          "nested key path with default bucket",
			location:      testNestedAppJS,
			defaultBucket: testBucketPrimary,
			want:          parseTarget{bucket: testBucketPrimary, key: testNestedAppJS},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, k, err := service.ParseLocation(tt.location, tt.defaultBucket)
			assertParsedLocation(t, parseTarget{bucket: b, key: k}, tt.want, err)
		})
	}
}
