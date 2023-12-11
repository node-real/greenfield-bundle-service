package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

// HandleCreateBundle handles create bundle request
func HandleCreateBundle() func(params bundle.CreateBundleParams) middleware.Responder {
	return func(params bundle.CreateBundleParams) middleware.Responder {
		if params.XBundleBucketName == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidBucketName)
		}
		if params.XBundleName == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidBundleName)
		}
		if params.Authorization == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// check signature
		signerAddress, err := types.VerifySignature(params.HTTPRequest)
		if err != nil {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// check expiry timestamp
		if err := types.ValidateExpiryTimestamp(params.HTTPRequest); err != nil {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidExpiryTimestamp)
		}

		bucketInfo, err := service.BundleSvc.QueryBucketFromGndf(params.XBundleBucketName)
		if err != nil {
			util.Logger.Errorf("query bucket error, err=%s", err.Error())
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInternalError)
		}

		// check if the signer is the owner of the bucket
		if bucketInfo.Owner != signerAddress.String() {
			util.Logger.Errorf("signer is not the owner of the bucket, signer=%s, bucket=%s", signerAddress.String(), params.XBundleBucketName)
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// todo: validate bundle params

		newBundle := database.Bundle{
			Owner:  signerAddress.String(),
			Bucket: params.XBundleBucketName,
			Name:   params.XBundleName,
		}

		// get bundler account for the user
		bundlerAccount, err := service.UserBundlerAccountSvc.GetOrCreateUserBundlerAccount(newBundle.Owner)
		if err != nil {
			util.Logger.Errorf("get bundler account for user error, user=%s, err=%s", newBundle.Owner, err.Error())
			return bundle.NewCreateBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}
		newBundle.BundlerAccount = bundlerAccount.BundlerAddress

		// create bundle
		_, err = service.BundleSvc.CreateBundle(newBundle)
		if err != nil {
			util.Logger.Errorf("create bundle error, bundle=%+v, err=%s", newBundle, err.Error())
			return bundle.NewCreateBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}

		return bundle.NewCreateBundleOK()
	}
}

// HandleFinalizeBundle handles finalize bundle request
func HandleFinalizeBundle() func(params bundle.FinalizeBundleParams) middleware.Responder {
	return func(params bundle.FinalizeBundleParams) middleware.Responder {
		if params.XBundleName == "" {
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.ErrorInvalidBundleName)
		}
		if params.XBundleBucketName == "" {
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.ErrorInvalidBucketName)
		}
		if params.Authorization == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// query bundle
		queriedBundle, err := service.BundleSvc.QueryBundle(params.XBundleBucketName, params.XBundleName)
		if err != nil {
			util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", params.XBundleBucketName, params.XBundleName, err.Error())
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}

		// check signature
		signerAddress, err := types.VerifySignature(params.HTTPRequest)
		if err != nil {
			util.Logger.Errorf("sig check error, err=%s", err.Error())
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// check expiry timestamp
		if err := types.ValidateExpiryTimestamp(params.HTTPRequest); err != nil {
			util.Logger.Errorf("validate expiry timestamp error, err=%s", err.Error())
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.ErrorInvalidExpiryTimestamp)
		}

		// check owner
		if signerAddress.String() != queriedBundle.Owner {
			util.Logger.Errorf("invalid bundle owner, signer=%s, bundleOwner=%s", signerAddress.String(), queriedBundle.Owner)
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.ErrorInvalidBundleOwner)
		}

		// finalize bundle
		_, err = service.BundleSvc.FinalizeBundle(params.XBundleBucketName, params.XBundleName)
		if err != nil {
			util.Logger.Errorf("finalize bundle error, bucket=%s, bundle=%s, err=%s", params.XBundleBucketName, params.XBundleName, err.Error())
			return bundle.NewFinalizeBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}

		return bundle.NewFinalizeBundleOK()
	}
}
