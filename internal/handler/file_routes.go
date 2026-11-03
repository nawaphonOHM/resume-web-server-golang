// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

// RouteRootPath is the URL path for the root SPA entrypoint endpoint.
const RouteRootPath rest.Pathz = "/"

// RouteSinglePath is the parameterized URL path for single-segment static file endpoints.
const RouteSinglePath rest.Pathz = "/:p1"

// RouteNestedPath is the parameterized URL path for nested static file endpoints.
// The framework rejects root catch-all routes starting with "/*" to protect health probes.
// Registering both "/:p1" and "/:p1/*fileName" allows matching both single-segment and
// multi-segment paths under Gin's radix tree without router conflicts.
const RouteNestedPath rest.Pathz = "/:p1/*fileName"

func buildExportableAPIs(h *FileHandler) []*rest.ExportableAPI {
	return []*rest.ExportableAPI{
		{
			Method:  rest.GET,
			Path:    RouteRootPath,
			Handler: h.ServeRoot,
		},
		{
			Method:  rest.GET,
			Path:    RouteSinglePath,
			Handler: h.ServeFile,
		},
		{
			Method:  rest.GET,
			Path:    RouteNestedPath,
			Handler: h.ServeFile,
		},
	}
}

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
		Apis:    buildExportableAPIs(handler),
	}
}

// RegisterRoutes creates a new BluePrint configured with SPA asset routes mounted without the /api prefix.
func RegisterRoutes(handler *FileHandler) *rest.BluePrint {
	return rest.NewBluePrint().WithAPIs(NewAPIRegistration(handler))
}
