// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

// RouteRootPath is the URL path for the root SPA entrypoint endpoint.
const RouteRootPath rest.Pathz = "/"

// RouteFilePath is the parameterized URL path for single-segment static file endpoints.
const RouteFilePath rest.Pathz = "/:fileName"

// RouteNestedFilePath resolves nested asset and SPA paths such as /assets/app.js.
// Three technical Gin route registrations ("/", "/:fileName", and "/:fileName/*rest")
// are used to serve the two logical SPA endpoints (root index and static asset resolution).
// A single root catch-all "/*fileName" cannot be used because Gin rejects wildcard
// conflicts with the framework's static /health and /ready probe routes.
// Note: static assets named exactly "health" or "ready" are shadowed by probe endpoints.
const RouteNestedFilePath rest.Pathz = "/:fileName/*rest"

// NewAPIRegistration builds the RRestAPIRegistration for SPA asset routes.
func NewAPIRegistration(handler *FileHandler) *rest.RRestAPIRegistration {
	if handler == nil {
		return &rest.RRestAPIRegistration{
			Version: 0,
			Prefix:  "",
			Group:   rest.APIGroupRaw,
		}
	}
	return &rest.RRestAPIRegistration{
		Version: 0,
		Prefix:  "",
		Group:   rest.APIGroupRaw,
		Apis: []*rest.ExportableAPI{
			{
				Method:  rest.GET,
				Path:    RouteRootPath,
				Handler: handler.ServeRoot,
			},
			{
				Method:  rest.GET,
				Path:    RouteFilePath,
				Handler: handler.ServeFile,
			},
			{
				Method:  rest.GET,
				Path:    RouteNestedFilePath,
				Handler: handler.ServeFile,
			},
		},
	}
}

// RegisterRoutes creates a new BluePrint configured with SPA asset routes mounted without the /api prefix.
func RegisterRoutes(handler *FileHandler) *rest.BluePrint {
	return rest.NewBluePrint().WithAPIs(NewAPIRegistration(handler))
}
