package handler

import (
	"net/http"
	"testing"
)

const (
	testParamFileName = "fileName"
	testAppJS         = "app.js"
	testWoff2         = "material-icons-LEZCGFVT.woff2"
	testThemeCSS      = "theme.css"
	testMediaDir      = "/media/"
	testMedia         = "media"
)

func TestResolveFileNameParam(t *testing.T) {
	if got := resolveFileNameParam(nil); got != "" {
		t.Errorf("expected empty string for nil context, got %q", got)
	}

	tests := []struct {
		params   map[string]string
		expected string
	}{
		{make(map[string]string), ""},
		{map[string]string{testParamP1: "   "}, ""},
		{map[string]string{testParamP1: testAppJS}, testAppJS},
		{map[string]string{testParamFileName: testAppJS}, testAppJS},
		{map[string]string{testParamP1: testMedia, testParamFileName: testSlash + testWoff2}, testWoff2},
		{map[string]string{testParamP1: "assets", testParamFileName: "/chunks/" + testAppJS}, testAppJS},
		{map[string]string{testParamP1: testMedia, testParamFileName: testSlash}, ""},
		{map[string]string{testParamP1: testAppJS, testParamFileName: testSlash}, ""},
		{map[string]string{testParamFileName: testMediaDir + testWoff2}, testWoff2},
		{map[string]string{testParamFileName: "/static/css/chunks/" + testThemeCSS}, testThemeCSS},
		{map[string]string{testParamFileName: testMediaDir}, ""},
		{map[string]string{testParamP1: testSlash}, ""},
		{map[string]string{testParamP1: "."}, ""},
		{map[string]string{testParamP1: ".."}, ""},
	}

	for _, tt := range tests {
		ctx := &mockRestContext{params: tt.params}
		assertEqual(t, tt.expected, resolveFileNameParam(ctx))
	}
}

func TestServeFileExtractsNestedFileName(t *testing.T) {
	tests := []struct {
		params   map[string]string
		expected string
	}{
		{map[string]string{testParamP1: testMedia, testParamFileName: testSlash + testWoff2}, testWoff2},
		{map[string]string{testParamP1: "assets", testParamFileName: "/chunks/" + testAppJS}, testAppJS},
		{map[string]string{testParamP1: testAppJS}, testAppJS},
		{map[string]string{testParamP1: testMedia, testParamFileName: testSlash}, ""},
		{map[string]string{testParamP1: testAppJS, testParamFileName: testSlash}, ""},
		{map[string]string{testParamFileName: testMediaDir + testWoff2}, testWoff2},
		{map[string]string{testParamFileName: "/index.html"}, "index.html"},
		{map[string]string{testParamFileName: testMediaDir}, ""},
		{map[string]string{testParamP1: ".."}, ""},
	}

	for _, tt := range tests {
		svc := &nameCapturingService{}
		h := NewFileHandler(&mockFileService{getAssetFunc: svc.getAsset})
		resp := h.ServeFile(&mockRestContext{params: tt.params})

		assertEqual(t, http.StatusOK, resp.StatusCode())
		assertEqual(t, tt.expected, svc.name)
	}
}
