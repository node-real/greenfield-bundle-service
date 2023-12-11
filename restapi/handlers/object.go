package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/models"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

func ValidateUploadObjectParams(params bundle.UploadObjectParams) *models.Error {
	if params.XBundleFileName == "" {
		return types.ErrorInvalidFileName
	}
	if params.XBundleContentType == "" {
		return types.ErrorInvalidContentType
	}
	if params.XBundleBucketName == "" {
		return types.ErrorInvalidBucketName
	}
	if params.Authorization == "" {
		return types.ErrorInvalidSignature
	}

	return nil
}

// ValidateFileContent validates the file content against the hash in the header
func ValidateFileContent(file io.ReadCloser, headerHash string) (io.ReadCloser, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		util.Logger.Errorf("read file error, err=%s", err.Error())
		return nil, err
	}

	hash := sha256.New()
	hash.Write(fileBytes)

	hashInBytes := hash.Sum(nil)[:]
	calculatedHash := hex.EncodeToString(hashInBytes)
	if calculatedHash != headerHash {
		util.Logger.Errorf("file hash does not match header hash, calculatedHash=%s, headerHash=%s", calculatedHash, headerHash)
		return nil, fmt.Errorf("file hash does not match header hash, calculatedHash=%s, headerHash=%s", calculatedHash, headerHash)
	}

	return io.NopCloser(bytes.NewReader(fileBytes)), nil
}

// HandleUploadObject handles the upload object request
func HandleUploadObject() func(params bundle.UploadObjectParams) middleware.Responder {
	return func(params bundle.UploadObjectParams) middleware.Responder {
		// check params
		if err := ValidateUploadObjectParams(params); err != nil {
			return bundle.NewUploadObjectBadRequest().WithPayload(err)
		}

		// check file content
		file, err := ValidateFileContent(params.File, params.XBundleFileSha256)
		if err != nil {
			util.Logger.Errorf("validate file content error, err=%s", err.Error())
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorInvalidFileContent)
		}

		// check signature
		signerAddress, err := types.VerifySignature(params.HTTPRequest)
		if err != nil {
			util.Logger.Errorf("sig check error, err=%s", err.Error())
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// check expiry timestamp
		if err := types.ValidateExpiryTimestamp(params.HTTPRequest); err != nil {
			util.Logger.Errorf("validate expiry timestamp error, err=%s", err.Error())
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorInvalidExpiryTimestamp)
		}

		// check if the signer is the owner of the bucket
		bucketInfo, err := service.BundleSvc.QueryBucketFromGndf(params.XBundleBucketName)
		if err != nil {
			util.Logger.Errorf("query bucket error, err=%s", err.Error())
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorInternalError)
		}

		if bucketInfo.Owner != signerAddress.String() {
			util.Logger.Errorf("signer is not the owner of the bucket, signer=%s, bucket=%s", signerAddress.String(), params.XBundleBucketName)
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorInvalidSignature)
		}

		// get bundling bundle
		bundlingBundle, err := service.BundleSvc.GetBundlingBundle(params.XBundleBucketName)
		if err != nil {
			util.Logger.Errorf("get bundling bundle error, bucket=%s, err=%s", params.XBundleBucketName, err.Error())
			return bundle.NewUploadObjectInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}

		// bundle not found
		if bundlingBundle.Id == 0 {
			// create new bundle
			newBundle := database.Bundle{
				Owner:  signerAddress.String(),
				Bucket: params.XBundleBucketName,
			}

			// get bundler account for the user
			bundlerAccount, err := service.UserBundlerAccountSvc.GetOrCreateUserBundlerAccount(newBundle.Owner)
			if err != nil {
				util.Logger.Errorf("get bundler account for user error, user=%s, err=%s", newBundle.Owner, err.Error())
				return bundle.NewUploadObjectInternalServerError()
			}
			newBundle.BundlerAccount = bundlerAccount.BundlerAddress

			// create bundle
			bundlingBundle, err = service.BundleSvc.CreateBundle(newBundle)
			if err != nil {
				util.Logger.Errorf("create bundle error, bundle=%+v, err=%s", newBundle, err.Error())
				return bundle.NewUploadObjectInternalServerError().WithPayload(types.InternalErrorWithError(err))
			}
		}

		// save object file to local storage
		_, fileSize, err := service.ObjectSvc.StoreObjectFile(params.XBundleBucketName, bundlingBundle.Name, params.XBundleFileName, file)
		if err != nil {
			util.Logger.Errorf("store object file error, err=%s", err.Error())
			return bundle.NewUploadObjectInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}

		// create object
		newObject := database.Object{
			Bucket:      params.XBundleBucketName,
			BundleName:  bundlingBundle.Name,
			ObjectName:  params.XBundleFileName,
			Owner:       signerAddress.String(),
			ContentType: params.XBundleContentType,
			Size:        fileSize,
		}
		if params.XBundleAttributes != nil {
			newObject.Attributes = *params.XBundleAttributes
		}

		_, err = service.ObjectSvc.CreateObjectForBundling(newObject)
		if err != nil {
			util.Logger.Errorf("create object error, object=%+v, err=%s", newObject, err.Error())
			return bundle.NewUploadObjectInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}

		return bundle.NewUploadObjectOK().WithPayload(&models.UploadObjectResponse{
			BundleName: bundlingBundle.Name,
		})
	}
}

// HandleViewBundleObject handles the view bundle object request
func HandleViewBundleObject() func(params bundle.ViewBundleObjectParams) middleware.Responder {
	return func(params bundle.ViewBundleObjectParams) middleware.Responder {
		object, err := service.ObjectSvc.GetObject(params.BucketName, params.BundleName, params.ObjectName)
		if err != nil {
			util.Logger.Errorf("get object error, bucket=%s, bundle=%s, object=%s, err=%s", params.BucketName, params.BundleName, params.ObjectName, err.Error())
			return bundle.NewViewBundleObjectInternalServerError()
		}

		if object.Id == 0 {
			return bundle.NewViewBundleObjectNotFound()
		}

		objectFile, _, err := service.ObjectSvc.GetObjectFile(params.BucketName, params.BundleName, params.ObjectName)
		if err != nil {
			util.Logger.Errorf("get object file error, bucket=%s, bundle=%s, object=%s, err=%s", params.BucketName, params.BundleName, params.ObjectName, err.Error())
			return bundle.NewViewBundleObjectInternalServerError()
		}

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       objectFile,
		}

		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.Header().Set("Content-Disposition", "inline")
			w.Header().Set("Content-Type", object.ContentType)
			io.Copy(w, response.Body)
		})
	}
}

// HandleDownloadBundleObject handles the download bundle object request
func HandleDownloadBundleObject() func(params bundle.DownloadBundleObjectParams) middleware.Responder {
	return func(params bundle.DownloadBundleObjectParams) middleware.Responder {
		object, err := service.ObjectSvc.GetObject(params.BucketName, params.BundleName, params.ObjectName)
		if err != nil {
			util.Logger.Errorf("get object error, bucket=%s, bundle=%s, object=%s, err=%s", params.BucketName, params.BundleName, params.ObjectName, err.Error())
			return bundle.NewViewBundleObjectInternalServerError()
		}

		if object.Id == 0 {
			return bundle.NewViewBundleObjectNotFound()
		}

		objectFile, _, err := service.ObjectSvc.GetObjectFile(params.BucketName, params.BundleName, params.ObjectName)
		if err != nil {
			util.Logger.Errorf("get object file error, bucket=%s, bundle=%s, object=%s, err=%s", params.BucketName, params.BundleName, params.ObjectName, err.Error())
			return bundle.NewViewBundleObjectInternalServerError()
		}

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       objectFile,
		}

		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", object.ObjectName))
			w.Header().Set("Content-Type", object.ContentType)
			io.Copy(w, response.Body)
		})
	}
}
