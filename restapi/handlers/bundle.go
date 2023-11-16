package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
)

func HandleManageBundle() func(params bundle.ManageBundleParams) middleware.Responder {
	return func(params bundle.ManageBundleParams) middleware.Responder {
		return bundle.NewManageBundleOK()
	}
}
