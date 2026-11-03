package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
)

const (
	testNestedWoff2 = "/media/material-icons-LEZCGFVT.woff2"
	testParamP1     = "p1"
	testSlash       = "/"
)

func TestExtractRawParam(t *testing.T) {
	assertEqual(t, "", extractRawParam(nil))
	assertEqual(t, "", extractRawParam(&mockRestContext{}))

	tests := []struct {
		params map[string]string
		want   string
	}{
		{map[string]string{testParamP1: "media", testParamFileName: testSlash + testWoff2}, testNestedWoff2},
		{map[string]string{testParamP1: "media", testParamFileName: testSlash}, "/media/"},
		{map[string]string{testParamP1: testAppJS, testParamFileName: testSlash}, "/app.js/"},
		{map[string]string{testParamP1: testAppJS}, testAppJS},
		{map[string]string{testParamFileName: testAppJS}, testAppJS},
		{map[string]string{testParamFileName: testNestedWoff2}, testNestedWoff2},
	}
	for _, tt := range tests {
		ctx := &mockRestContext{params: tt.params}
		assertEqual(t, tt.want, extractRawParam(ctx))
	}
}

func TestResolveFileParams(t *testing.T) {
	raw, name, ctx := resolveFileParams(nil)
	assertEqual(t, "", raw)
	assertEqual(t, "", name)
	assertEqual(t, context.Background(), ctx)

	reqCtx := &mockRestContext{params: map[string]string{
		testParamP1:       "media",
		testParamFileName: testSlash + testWoff2,
	}}
	raw, name, _ = resolveFileParams(reqCtx)
	assertEqual(t, testNestedWoff2, raw)
	assertEqual(t, testWoff2, name)
}

func TestDispatchResponses(t *testing.T) {
	ctx := context.Background()
	asset := &domain.StaticAsset{
		Body:        newSpyReadCloser(strings.NewReader(testBody)),
		ContentType: testType,
	}
	respFile := dispatchAssetResponse(ctx, testAppJS, asset)
	assertEqual(t, http.StatusOK, respFile.StatusCode())

	respRoot := dispatchRootResponse(ctx, asset)
	assertEqual(t, http.StatusOK, respRoot.StatusCode())
}

func TestDispatchResponsesNilAsset(t *testing.T) {
	ctx := context.Background()
	respFile := dispatchAssetResponse(ctx, testAppJS, nil)
	assertEqual(t, http.StatusInternalServerError, respFile.StatusCode())

	respRoot := dispatchRootResponse(ctx, nil)
	assertEqual(t, http.StatusInternalServerError, respRoot.StatusCode())
}
