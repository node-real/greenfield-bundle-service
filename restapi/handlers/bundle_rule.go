package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/rule"
)

func HandleAddBundleRule() func(params rule.AddBundleRuleParams) middleware.Responder {
	return func(params rule.AddBundleRuleParams) middleware.Responder {
		return rule.NewAddBundleRuleOK()
	}
}
