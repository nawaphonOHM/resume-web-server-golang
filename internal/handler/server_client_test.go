// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"io"
	"net/http"
	"testing"
)

// testClient disables keep-alive so no idle connection can delay graceful shutdown.
var testClient = &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}

func doGet(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	return testClient.Do(req)
}

// fetch performs a GET and returns the response with its fully read body.
func fetch(t *testing.T, url string) (*http.Response, string) {
	t.Helper()
	resp, err := doGet(context.Background(), url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer closeBodyOf(t, resp)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp, string(data)
}
