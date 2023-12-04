package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/rule"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

func HandleSetBundleRule() func(params rule.SetBundleRuleParams) middleware.Responder {
	return func(params rule.SetBundleRuleParams) middleware.Responder {
		// check params
		if params.XBundleBucketName == "" {
			return rule.NewSetBundleRuleBadRequest().WithPayload(ErrorInvalidBucketName)
		}
		if params.Authorization == "" {
			return rule.NewSetBundleRuleBadRequest().WithPayload(ErrorInvalidSignature)
		}
		// check signature
		signerAddress, err := types.VerifySignature(params.HTTPRequest)
		if err != nil {
			util.Logger.Errorf("sig check error, err=%s", err.Error())
			return rule.NewSetBundleRuleBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// check expiry timestamp
		if err := types.ValidateExpiryTimestamp(params.HTTPRequest); err != nil {
			util.Logger.Errorf("validate expiry timestamp error, err=%s", err.Error())
			return rule.NewSetBundleRuleBadRequest().WithPayload(ErrorInvalidExpiryTimestamp)
		}

		// create or update bundle rule
		_, err = service.BundleRuleSvc.CreateOrUpdateBundleRule(signerAddress,
			params.XBundleBucketName,
			params.XBundleMaxBundleFiles,
			params.XBundleMaxBundleSize,
			params.XBundleMaxFinalizeTime,
		)
		if err != nil {
			util.Logger.Errorf("create or update bundle rule error, err=%s", err.Error())
			return rule.NewSetBundleRuleInternalServerError().WithPayload(ErrorInternalError)
		}

		return rule.NewSetBundleRuleOK()
	}
}
