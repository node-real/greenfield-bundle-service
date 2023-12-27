package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/models"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

// ValidateFileContent validates the file content against the hash in the header
func ValidateFileContent(params bundle.UploadObjectParams) (io.ReadCloser, error) {
	// check tags
	if params.XBundleTags != nil && *params.XBundleTags != "" {
		// json unmarshal
		var tags map[string]string
		err := json.Unmarshal([]byte(*params.XBundleTags), &tags)
		if err != nil {
			util.Logger.Errorf("invalid tags, err=%s", err.Error())
			return nil, fmt.Errorf("invalid tags, err=%s", err.Error())
		}
	}

	contentLength := params.HTTPRequest.Header.Get("Content-Length")

	fileSize, err := strconv.Atoi(contentLength)
	if err != nil {
		util.Logger.Errorf("invalid Content-Length header, err=%s", err.Error())
		return nil, err
	}
	if fileSize > types.DefaultMaxFileSize {
		util.Logger.Errorf("file size exceeds limit, size=%d", fileSize)
		return nil, fmt.Errorf("file size exceeds limit, size=%d", fileSize)
	}

	fileBytes, err := io.ReadAll(params.File)
	if err != nil {
		util.Logger.Errorf("read file error, err=%s", err.Error())
		return nil, err
	}

	hash := sha256.New()
	hash.Write(fileBytes)

	hashInBytes := hash.Sum(nil)[:]
	calculatedHash := hex.EncodeToString(hashInBytes)
	if calculatedHash != params.XBundleFileSha256 {
		util.Logger.Errorf("file hash does not match header hash, calculatedHash=%s, headerHash=%s", calculatedHash, params.XBundleFileSha256)
		return nil, fmt.Errorf("file hash does not match header hash, calculatedHash=%s, headerHash=%s", calculatedHash, params.XBundleFileSha256)
	}

	return io.NopCloser(bytes.NewReader(fileBytes)), nil
}

func GetBundlingBundle(bucketName string, signerAddress string) (database.Bundle, *models.Error) {
	// get bundling bundle
	bundlingBundle, err := service.BundleSvc.GetBundlingBundle(bucketName)
	if err != nil {
		util.Logger.Errorf("get bundling bundle error, bucket=%s, err=%s", bucketName, err.Error())
		return database.Bundle{}, types.InternalErrorWithError(err)
	}

	// bundle not found
	if bundlingBundle.Id == 0 {
		// create new bundle
		newBundle := database.Bundle{
			Owner:  signerAddress,
			Bucket: bucketName,
		}

		// get bundler account for the user
		bundlerAccount, err := service.UserBundlerAccountSvc.GetOrCreateUserBundlerAccount(newBundle.Owner)
		if err != nil {
			util.Logger.Errorf("get bundler account for user error, user=%s, err=%s", newBundle.Owner, err.Error())
			return database.Bundle{}, types.InternalErrorWithError(err)
		}
		newBundle.BundlerAccount = bundlerAccount.BundlerAddress

		// create bundle
		bundlingBundle, err = service.BundleSvc.CreateBundle(newBundle)
		if err != nil {
			util.Logger.Errorf("create bundle error, bundle=%+v, err=%s", newBundle, err.Error())
			return database.Bundle{}, types.InternalErrorWithError(err)
		}
	}
	return bundlingBundle, nil
}

// HandleUploadObject handles the upload object request
func HandleUploadObject() func(params bundle.UploadObjectParams) middleware.Responder {
	return func(params bundle.UploadObjectParams) middleware.Responder {
		// check file content
		file, err := ValidateFileContent(params)
		if err != nil {
			util.Logger.Errorf("validate file content error, err=%s", err.Error())
			return bundle.NewUploadObjectBadRequest().WithPayload(types.InvalidFileContentErrorWithError(err))
		}

		// validate headers
		signerAddress, merr := types.ValidateHeaders(params.HTTPRequest)
		if merr != nil {
			util.Logger.Errorf("sig check error, code=%d, msg=%s", merr.Code, merr.Message)
			return bundle.NewUploadObjectBadRequest().WithPayload(merr)
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
		bundlingBundle, merr := GetBundlingBundle(params.XBundleBucketName, signerAddress.String())
		if err != nil {
			util.Logger.Errorf("get bundling bundle error, bucket=%s, err=%s", params.XBundleBucketName, err.Error())
			return bundle.NewUploadObjectInternalServerError().WithPayload(merr)
		}

		// check bundle size and files against the limit
		if bundlingBundle.Size > bundlingBundle.MaxSize || bundlingBundle.Files > bundlingBundle.MaxFiles {
			util.Logger.Errorf("bundle size exceeds limit, size=%d, maxSize=%d, files=%d, maxFiles=%d", bundlingBundle.Size, bundlingBundle.MaxSize, bundlingBundle.Files, bundlingBundle.MaxFiles)
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorBundleSizeExceedsLimit)
		}

		// check if the object already exists
		queriedObject, err := service.ObjectSvc.GetObject(params.XBundleBucketName, bundlingBundle.Name, params.XBundleFileName)
		if err != nil {
			util.Logger.Errorf("get object error, bucket=%s, bundle=%s, object=%s, err=%s", params.XBundleBucketName, bundlingBundle.Name, params.XBundleFileName, err.Error())
			return bundle.NewUploadObjectInternalServerError().WithPayload(types.InternalErrorWithError(err))
		}
		if queriedObject.Id != 0 {
			util.Logger.Errorf("object already exists, bucket=%s, bundle=%s, object=%s", params.XBundleBucketName, bundlingBundle.Name, params.XBundleFileName)
			return bundle.NewUploadObjectBadRequest().WithPayload(types.ErrorObjectAlreadyExists)
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

		// check tags
		if params.XBundleTags != nil && *params.XBundleTags != "" {
			newObject.Tags = *params.XBundleTags

			var tags map[string]string
			err = json.Unmarshal([]byte(newObject.Tags), &tags)
			if err != nil {
				util.Logger.Warnf("unmarshal tags failed, tags=%s, err=%v", newObject.Tags, err.Error())
				return bundle.NewUploadObjectInternalServerError().WithPayload(types.ErrorInvalidTags)
			}
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

		objectFile, err := service.ObjectSvc.GetObjectFile(params.BucketName, params.BundleName, params.ObjectName)
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
			_, err := io.Copy(w, response.Body)
			if err != nil {
				util.Logger.Errorf("copy object file error, err=%s", err.Error())
			}
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

		objectFile, err := service.ObjectSvc.GetObjectFile(params.BucketName, params.BundleName, params.ObjectName)
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
			_, err := io.Copy(w, response.Body)
			if err != nil {
				util.Logger.Errorf("copy object file error, err=%s", err.Error())
			}
		})
	}
}
