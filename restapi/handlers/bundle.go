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
			return bundle.NewCreateBundleBadRequest().WithPayload(ErrorInvalidBucketName)
		}
		if params.XBundleName == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(ErrorInvalidBundleName)
		}
		if params.Authorization == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// check signature
		signerAddress, err := types.VerifySignature(params.HTTPRequest)
		if err != nil {
			return bundle.NewCreateBundleBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// check expiry timestamp
		if err := types.ValidateExpiryTimestamp(params.HTTPRequest); err != nil {
			return bundle.NewCreateBundleBadRequest().WithPayload(ErrorInvalidExpiryTimestamp)
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
			return bundle.NewCreateBundleBadRequest()
		}
		newBundle.BundlerAccount = bundlerAccount.BundlerAddress

		// create bundle
		_, err = service.BundleSvc.CreateBundle(newBundle)
		if err != nil {
			util.Logger.Errorf("create bundle error, bundle=%+v, err=%s", newBundle, err.Error())
			return bundle.NewCreateBundleBadRequest() // todo: return proper error
		}

		return bundle.NewCreateBundleOK()
	}
}

// HandleFinalizeBundle handles finalize bundle request
func HandleFinalizeBundle() func(params bundle.FinalizeBundleParams) middleware.Responder {
	return func(params bundle.FinalizeBundleParams) middleware.Responder {
		// todo: make sure the owner can only finalize the bundle created by himself manually
		if params.XBundleName == "" {
			return bundle.NewFinalizeBundleBadRequest().WithPayload(ErrorInvalidBundleName)
		}
		if params.XBundleBucketName == "" {
			return bundle.NewFinalizeBundleBadRequest().WithPayload(ErrorInvalidBucketName)
		}
		if params.Authorization == "" {
			return bundle.NewCreateBundleBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// query bundle
		queriedBundle, err := service.BundleSvc.QueryBundle(params.XBundleBucketName, params.XBundleName)
		if err != nil {
			util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", params.XBundleBucketName, params.XBundleName, err.Error())
			return bundle.NewFinalizeBundleBadRequest()
		}

		// check signature
		signerAddress, err := types.VerifySignature(params.HTTPRequest)
		if err != nil {
			util.Logger.Errorf("sig check error, err=%s", err.Error())
			return bundle.NewFinalizeBundleBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// check expiry timestamp
		if err := types.ValidateExpiryTimestamp(params.HTTPRequest); err != nil {
			util.Logger.Errorf("validate expiry timestamp error, err=%s", err.Error())
			return bundle.NewFinalizeBundleBadRequest().WithPayload(ErrorInvalidExpiryTimestamp)
		}

		// check owner
		if signerAddress.String() != queriedBundle.Owner {
			util.Logger.Errorf("invalid bundle owner, signer=%s, bundleOwner=%s", signerAddress.String(), queriedBundle.Owner)
			return bundle.NewFinalizeBundleBadRequest().WithPayload(ErrorInvalidBundleOwner)
		}

		// finalize bundle
		_, err = service.BundleSvc.FinalizeBundle(params.XBundleBucketName, params.XBundleName)
		if err != nil {
			util.Logger.Errorf("finalize bundle error, bucket=%s, bundle=%s, err=%s", params.XBundleBucketName, params.XBundleName, err.Error())
			return bundle.NewFinalizeBundleBadRequest() // todo: return proper error
		}

		return bundle.NewFinalizeBundleOK()
	}
}
