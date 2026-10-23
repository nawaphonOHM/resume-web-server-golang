// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"context"

	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

// extractRequestContext returns the request-scoped context or context.Background() when unavailable.
func extractRequestContext(c rest.Context) context.Context {
	if c == nil {
		return context.Background()
	}
	if ctx := c.Context(); ctx != nil {
		return ctx
	}
	return context.Background()
}
