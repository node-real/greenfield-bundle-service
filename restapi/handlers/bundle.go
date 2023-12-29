package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/models"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

// HandleDeleteBundle handles delete bundle request
func HandleDeleteBundle() func(params bundle.DeleteBundleParams) middleware.Responder {
	return func(params bundle.DeleteBundleParams) middleware.Responder {
		// validate headers
		signerAddress, merr := types.ValidateHeaders(params.HTTPRequest)
		if merr != nil {
			util.Logger.Errorf("sig check error, code=%d, msg=%s", merr.Code, merr.Message)
			return bundle.NewCreateBundleBadRequest().WithPayload(merr)
		}

		// check bundle name prefix
		if service.IsAutoGeneratedBundleName(params.XBundleName) {
			util.Logger.Errorf("bundle name should not start with %s", service.BundleNamePrefix)
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.ErrorInvalidBundleName)
		}

		// check existence and status of the bundle
		queriedBundle, err := service.BundleSvc.QueryBundle(params.XBundleBucketName, params.XBundleName)
		if err != nil {
			util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", params.XBundleBucketName, params.XBundleName, err.Error())
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}
		if queriedBundle == nil {
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.ErrorBundleNotExist)
		}

		// check bundle status, can not delete finalized bundle
		if queriedBundle.Status == database.BundleStatusFinalized {
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.ErrorInvalidBundleStatus)
		}

		// check if the signer is the owner of the bundle
		bucketInfo, err := service.BundleSvc.QueryBucketFromGndf(params.XBundleBucketName)
		if err != nil {
			util.Logger.Errorf("query bucket error, err=%s", err.Error())
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.ErrorInternalError)
		}
		if bucketInfo.Owner != signerAddress.String() {
			util.Logger.Errorf("signer is not the owner of the bucket, signer=%s, bucket=%s", signerAddress.String(), params.XBundleBucketName)
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// delete bundle
		err = service.BundleSvc.DeleteBundle(params.XBundleBucketName, params.XBundleName)
		if err != nil {
			util.Logger.Errorf("delete bundle error, bucket=%s, bundle=%s, err=%s", params.XBundleBucketName, params.XBundleName, err.Error())
			return bundle.NewDeleteBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}

		return bundle.NewDeleteBundleOK()
	}
}

// HandleCreateBundle handles create bundle request
func HandleCreateBundle() func(params bundle.CreateBundleParams) middleware.Responder {
	return func(params bundle.CreateBundleParams) middleware.Responder {
		// validate headers
		signerAddress, merr := types.ValidateHeaders(params.HTTPRequest)
		if merr != nil {
			util.Logger.Errorf("sig check error, code=%d, msg=%s", merr.Code, merr.Message)
			return bundle.NewCreateBundleBadRequest().WithPayload(merr)
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

		// check bundle name prefix
		if service.IsAutoGeneratedBundleName(params.XBundleName) {
			util.Logger.Errorf("bundle name should not start with %s", service.BundleNamePrefix)
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorInvalidBundleName)
		}

		// validate bundle name
		if err := types.ValidateBundleName(params.XBundleName); err != nil {
			util.Logger.Errorf("invalid bundle name, err=%s", err.Message)
			return bundle.NewCreateBundleBadRequest().WithPayload(err)
		}

		// check the existence of the bundle in Greenfield
		_, err = service.BundleSvc.HeadObjectFromGnfd(params.XBundleBucketName, params.XBundleName)
		if err == nil {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.ErrorObjectExist)
		}
		if !service.IsObjectNotFoundError(err) {
			return bundle.NewCreateBundleBadRequest().WithPayload(types.InternalErrorWithError(err))
		}

		// create new bundle
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

// HandleQueryBundle handles the query bundle request
func HandleQueryBundle() func(params bundle.QueryBundleParams) middleware.Responder {
	return func(params bundle.QueryBundleParams) middleware.Responder {
		bundleInfo, err := service.BundleSvc.QueryBundle(params.BucketName, params.BundleName)
		if err != nil {
			util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", params.BucketName, params.BundleName, err.Error())
			return bundle.NewQueryBundleInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}

		if bundleInfo.Id == 0 {
			return bundle.NewQueryBundleNotFound()
		}

		return bundle.NewQueryBundleOK().WithPayload(&models.QueryBundleResponse{
			BucketName:       bundleInfo.Bucket,
			BundleName:       bundleInfo.Name,
			Status:           int64(bundleInfo.Status),
			Files:            bundleInfo.Files,
			Size:             bundleInfo.Size,
			ErrorMessage:     bundleInfo.ErrMessage,
			CreatedTimestamp: bundleInfo.CreatedAt.Unix(),
		})
	}
}

// HandleQueryBundlingBundle handles the query bundling bundle request
func HandleQueryBundlingBundle() func(params bundle.QueryBundlingBundleParams) middleware.Responder {
	return func(params bundle.QueryBundlingBundleParams) middleware.Responder {
		bundleInfo, err := service.BundleSvc.GetBundlingBundle(params.BucketName)
		if err != nil {
			util.Logger.Errorf("query bundle error, bucket=%s, err=%s", params.BucketName, err.Error())
			return bundle.NewQueryBundleInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}

		if bundleInfo.Id == 0 {
			return bundle.NewQueryBundlingBundleNotFound()
		}

		return bundle.NewQueryBundleOK().WithPayload(&models.QueryBundleResponse{
			BucketName:       bundleInfo.Bucket,
			BundleName:       bundleInfo.Name,
			Status:           int64(bundleInfo.Status),
			Files:            bundleInfo.Files,
			Size:             bundleInfo.Size,
			ErrorMessage:     bundleInfo.ErrMessage,
			CreatedTimestamp: bundleInfo.CreatedAt.Unix(),
		})
	}
}
