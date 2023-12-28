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
		// validate headers
		signerAddress, merr := types.ValidateHeaders(params.HTTPRequest)
		if merr != nil {
			util.Logger.Errorf("sig check error, code=%d, msg=%s", merr.Code, merr.Message)
			return rule.NewSetBundleRuleBadRequest().WithPayload(merr)
		}

		bucket, err := service.BundleSvc.QueryBucketFromGndf(params.XBundleBucketName)
		if err != nil {
			util.Logger.Errorf("query bucket error, err=%s", err.Error())
			return rule.NewSetBundleRuleInternalServerError().WithPayload(types.ErrorInternalError)
		}

		// check if the signer is the owner of the bucket
		if bucket.Owner != signerAddress.String() {
			util.Logger.Errorf("signer is not the owner of the bucket, signer=%s, bucket=%s", signerAddress.String(), params.XBundleBucketName)
			return rule.NewSetBundleRuleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// check rule params
		if params.XBundleMaxBundleFiles > types.MaxBundleFiles || params.XBundleMaxBundleSize > types.MaxBundleSize || params.XBundleMaxFinalizeTime > types.MaxFinalizeTime {
			util.Logger.Errorf("invalid rule params, maxBundleFiles=%d, maxBundleSize=%d, maxFinalizeTime=%d", params.XBundleMaxBundleFiles, params.XBundleMaxBundleSize, params.XBundleMaxFinalizeTime)
			return rule.NewSetBundleRuleBadRequest().WithPayload(types.ErrorInvalidBundleRuleParams)
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
			return rule.NewSetBundleRuleInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}

		return rule.NewSetBundleRuleOK()
	}
}
