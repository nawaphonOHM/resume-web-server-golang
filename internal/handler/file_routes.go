// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

// RouteRootPath is the URL path for the root SPA entrypoint endpoint.
const RouteRootPath rest.Pathz = "/"

// RouteFilePath is the parameterized URL path for single-segment static file endpoints.
const RouteFilePath rest.Pathz = "/:fileName"

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
		},
	}
}

// RegisterRoutes creates a new BluePrint configured with SPA asset routes mounted without the /api prefix.
func RegisterRoutes(handler *FileHandler) *rest.BluePrint {
	return rest.NewBluePrint().WithAPIs(NewAPIRegistration(handler))
}
