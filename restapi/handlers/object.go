package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
)

func HandleUploadObject() func(params bundle.UploadObjectParams) middleware.Responder {
	return func(params bundle.UploadObjectParams) middleware.Responder {
		return bundle.NewUploadObjectOK()
	}
}

func HandleBundleObject() func(params bundle.BundleObjectParams) middleware.Responder {
	return func(params bundle.BundleObjectParams) middleware.Responder {
		return bundle.NewBundleObjectOK()
	}
}
