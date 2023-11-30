package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	UploadObjectMethod = "uploadObject"
)

type ObjectSignMessage struct {
	Method      string
	BucketName  string
	BundleName  string
	FileName    string
	ContentType string
	Timestamp   int64
}

func (s *ObjectSignMessage) SignBytes() ([]byte, error) {
	return json.Marshal(s)
}

// SigCheckUploadObject checks the signature of upload object request
func SigCheckUploadObject(params bundle.UploadObjectParams) (common.Address, error) {
	signMessage := ObjectSignMessage{
		Method:      UploadObjectMethod,
		BucketName:  params.BucketName,
		BundleName:  *params.BundleName,
		FileName:    params.FileName,
		ContentType: params.ContentType,
		Timestamp:   params.Timestamp,
	}

	signBytes, err := signMessage.SignBytes()
	if err != nil {
		return common.Address{}, err
	}
	messageHash := crypto.Keccak256Hash(signBytes)

	sigBytes, err := hex.DecodeString(params.XSignature)
	if err != nil {
		return common.Address{}, err
	}
	isValid, err := util.VerifySignature(messageHash.Bytes(), sigBytes)
	if err != nil {
		return common.Address{}, err
	}
	if !isValid {
		return common.Address{}, fmt.Errorf("invalid signature")
	}

	address, err := util.RecoverAddress(messageHash, sigBytes)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

func HandleUploadObject() func(params bundle.UploadObjectParams) middleware.Responder {
	return func(params bundle.UploadObjectParams) middleware.Responder {
		// check params
		if params.FileName == "" {
			return bundle.NewUploadObjectBadRequest()
		}
		if params.Timestamp == 0 {
			return bundle.NewUploadObjectBadRequest()
		}
		if params.ContentType == "" {
			return bundle.NewUploadObjectBadRequest()
		}
		if params.BucketName == "" {
			return bundle.NewUploadObjectBadRequest()
		}
		if params.XSignature == "" {
			return bundle.NewUploadObjectBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// check signature
		signerAddress, err := SigCheckUploadObject(params)
		if err != nil {
			util.Logger.Errorf("sig check upload object error, err=%s", err.Error())
			return bundle.NewUploadObjectBadRequest().WithPayload(ErrorInvalidSignature)
		}

		// get bundling bundle
		bundlingBundle, err := service.BundleSvc.GetBundlingBundle(params.BucketName)
		if err != nil {
			util.Logger.Errorf("get bundling bundle error, bucket=%s, err=%s", params.BucketName, err.Error())
			return bundle.NewUploadObjectInternalServerError()
		}

		// bundle not found
		if bundlingBundle.Id == 0 {
			// create new bundle
			newBundle := database.Bundle{
				Owner:  signerAddress.String(),
				Bucket: params.BucketName,
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
				return bundle.NewUploadObjectInternalServerError()
			}
		}

		// create object
		newObject := database.Object{ // TODO: add more fields
			Bucket:      params.BucketName,
			BundleName:  bundlingBundle.Name,
			ObjectName:  params.FileName,
			Owner:       signerAddress.String(),
			ContentType: params.ContentType,
		}

		_, err = service.ObjectSvc.CreateObjectForBundling(newObject)
		if err != nil {
			util.Logger.Errorf("create object error, object=%+v, err=%s", newObject, err.Error())
			return bundle.NewUploadObjectInternalServerError()
		}

		return bundle.NewUploadObjectOK()
	}
}

func HandleBundleObject() func(params bundle.BundleObjectParams) middleware.Responder {
	return func(params bundle.BundleObjectParams) middleware.Responder {
		return bundle.NewBundleObjectOK()
	}
}
