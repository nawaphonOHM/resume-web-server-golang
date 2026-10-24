// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const decimalBase = 10

// applyOptionalHeader adds a header if value is non-empty.
func applyOptionalHeader(resp rest.Response, key, value string) rest.Response {
	if value != "" {
		return resp.WithHeader(key, value)
	}
	return resp
}

// applyOptionalIntHeader adds a header if val > 0; int64 is formatted directly to avoid 32-bit overflow.
func applyOptionalIntHeader(resp rest.Response, key string, val int64) rest.Response {
	if val > 0 {
		return resp.WithHeader(key, strconv.FormatInt(val, decimalBase))
	}
	return resp
}

// applyOptionalStringerHeader adds a header if time is non-zero.
func applyOptionalStringerHeader(resp rest.Response, key string, t time.Time) rest.Response {
	if !t.IsZero() {
		return resp.WithHeaderStringer(key, httpTime{t: t.UTC()})
	}
	return resp
}

// httpTime wraps a time.Time to format it as an HTTP-compliant date string.
type httpTime struct {
	t time.Time
}

func (h httpTime) String() string {
	return h.t.Format(http.TimeFormat)
}
