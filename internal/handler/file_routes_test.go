// Package handler provides HTTP handlers and responses for serving static SPA assets.
package handler

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const expectedRoutes = 2

func TestNewAPIRegistrationRawGroup(t *testing.T) {
	reg := NewAPIRegistration(NewFileHandler(&mockFileService{}))

	assertEqual(t, rest.APIGroupRaw, reg.Group)
	assertEqual(t, expectedRoutes, len(reg.Apis))
	assertEqual(t, RouteRootPath, reg.Apis[0].Path)
	assertEqual(t, RouteFilePath, reg.Apis[1].Path)
	for _, api := range reg.Apis {
		assertEqual(t, rest.GET, api.Method)
	}
}

func TestNewAPIRegistrationNilHandler(t *testing.T) {
	reg := NewAPIRegistration(nil)

	assertEqual(t, 0, len(reg.Apis))
	assertEqual(t, rest.APIGroupRaw, reg.Group)
}

func TestRegisterRoutesReturnsBluePrint(t *testing.T) {
	assertEqual(t, true, RegisterRoutes(NewFileHandler(&mockFileService{})) != nil)
}
