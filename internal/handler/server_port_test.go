// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"
	"net"
	"net/http"
	"testing"
)

func listenLoopback(t *testing.T) net.Listener {
	t.Helper()
	lc := net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", loopbackHost+":0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	return l
}

func freePort(t *testing.T) string {
	t.Helper()
	l := listenLoopback(t)
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}
	if closeErr := l.Close(); closeErr != nil {
		t.Fatalf("close listener: %v", closeErr)
	}
	return port
}

func closeBodyOf(t *testing.T, resp *http.Response) {
	t.Helper()
	if err := resp.Body.Close(); err != nil {
		t.Logf("close response body: %v", err)
	}
}
