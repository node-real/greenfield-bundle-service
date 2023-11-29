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
	CreateBundleMethod   = "createBundle"
	FinalizeBundleMethod = "finalizeBundle"
)

type BundleSignMessage struct {
	Method     string
	BucketName string
	BundleName string
	Timestamp  int64
}

func (s *BundleSignMessage) SignBytes() ([]byte, error) {
	return json.Marshal(s)
}

// SigCheckCreateBundle checks the signature of create bundle request
func SigCheckCreateBundle(params bundle.CreateBundleParams) (common.Address, error) {
	signMessage := BundleSignMessage{
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

// HandleCreateBundle handles create bundle request
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

// SigCheckFinalizeBundle checks the signature of finalize bundle request
func SigCheckFinalizeBundle(params bundle.FinalizeBundleParams) (common.Address, error) {
	signMessage := BundleSignMessage{
		Method:     FinalizeBundleMethod,
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

// HandleFinalizeBundle handles finalize bundle request
func HandleFinalizeBundle() func(params bundle.FinalizeBundleParams) middleware.Responder {
	return func(params bundle.FinalizeBundleParams) middleware.Responder {
		// todo: make sure the owner can only finalize the bundle created by himself manually
		// todo: handle the errors
		if params.Body.BundleName != nil {
			return bundle.NewFinalizeBundleBadRequest()
		}
		if params.Body.BucketName != nil {
			return bundle.NewFinalizeBundleBadRequest()
		}
		if params.Body.Timestamp != nil {
			return bundle.NewFinalizeBundleBadRequest()
		}

		// query bundle
		queriedBundle, err := service.BundleSvc.QueryBundle(*params.Body.BucketName, *params.Body.BundleName)
		if err != nil {
			return bundle.NewFinalizeBundleBadRequest()
		}

		// check signature
		signerAddress, err := SigCheckFinalizeBundle(params)
		if err != nil {
			return bundle.NewFinalizeBundleBadRequest()
		}

		// check owner
		if signerAddress.String() != queriedBundle.Owner {
			return bundle.NewFinalizeBundleBadRequest()
		}

		// finalize bundle
		_, err = service.BundleSvc.FinalizeBundle(*params.Body.BucketName, *params.Body.BundleName)
		if err != nil {
			return bundle.NewFinalizeBundleBadRequest()
		}

		return bundle.NewFinalizeBundleOK()
	}
}
