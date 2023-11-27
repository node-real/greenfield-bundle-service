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
	CreateBundleMethod = "createBundle"
)

type CreateBundleSignMessage struct {
	Method     string
	BucketName string
	BundleName string
	Timestamp  int64
}

func (s *CreateBundleSignMessage) SignBytes() ([]byte, error) {
	return json.Marshal(s)
}

func SigCheckCreateBundle(params bundle.CreateBundleParams) (common.Address, error) {
	signMessage := CreateBundleSignMessage{
		Method:     CreateBundleMethod,
		BucketName: *params.Body.BucketName,
		BundleName: *params.Body.BundleName,
		Timestamp:  *params.Body.Timestamp,
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

	return address, err
}

func HandleCreateBundle() func(params bundle.CreateBundleParams) middleware.Responder {
	return func(params bundle.CreateBundleParams) middleware.Responder {
		// todo: handle the errors
		if params.Body.BundleName != nil {
			return bundle.NewCreateBundleBadRequest()
		}
		if params.Body.Timestamp != nil {
			return bundle.NewCreateBundleBadRequest()
		}
		if params.Body.BundleName != nil {
			return bundle.NewCreateBundleBadRequest()
		}
		signerAddress, err := SigCheckCreateBundle(params)
		if err != nil {
			return bundle.NewCreateBundleBadRequest()
		}

		newBundle := database.Bundle{
			Owner:  signerAddress.String(),
			Bucket: *params.Body.BucketName,
			Name:   *params.Body.BundleName,
		}

		_, err = service.BundleSvc.CreateBundle(newBundle)
		if err != nil {
			return bundle.NewCreateBundleBadRequest()
		}

		return bundle.NewCreateBundleOK()
	}
}

func HandleFinalizeBundle() func(params bundle.FinalizeBundleParams) middleware.Responder {
	return func(params bundle.FinalizeBundleParams) middleware.Responder {
		return bundle.NewFinalizeBundleOK()
	}
}
