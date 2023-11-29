package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/rule"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	SetBundleRuleMethod = "setBundleRule"
)

type BundleRuleSignMessage struct {
	Method          string
	BucketName      string
	MaxFiles        int64
	MaxSize         int64
	MaxFinalizeTime int64
	Timestamp       int64
}

// SignBytes returns the bytes to be signed
func (s *BundleRuleSignMessage) SignBytes() ([]byte, error) {
	return json.Marshal(s)
}

// SigCheckSetBundleRule checks the signature of set bundle rule request
func SigCheckSetBundleRule(params rule.SetBundleRuleParams) (common.Address, error) {
	signMessage := BundleRuleSignMessage{
		Method:          SetBundleRuleMethod,
		BucketName:      *params.Body.BucketName,
		MaxFiles:        *params.Body.MaxBundleFiles,
		MaxSize:         *params.Body.MaxBundleSize,
		MaxFinalizeTime: *params.Body.MaxFinalizeTime,
		Timestamp:       *params.Body.Timestamp,
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

func HandleSetBundleRule() func(params rule.SetBundleRuleParams) middleware.Responder {
	return func(params rule.SetBundleRuleParams) middleware.Responder {
		// check params
		if params.Body.Timestamp != nil {
			return rule.NewSetBundleRuleBadRequest()
		}
		if params.Body.BucketName != nil {
			return rule.NewSetBundleRuleBadRequest()
		}
		if params.Body.MaxBundleFiles != nil {
			return rule.NewSetBundleRuleBadRequest()
		}
		if params.Body.MaxBundleSize != nil {
			return rule.NewSetBundleRuleBadRequest()
		}
		if params.Body.MaxFinalizeTime != nil {
			return rule.NewSetBundleRuleBadRequest()
		}

		// check signature
		signerAddress, err := SigCheckSetBundleRule(params)
		if err != nil {
			return rule.NewSetBundleRuleBadRequest()
		}

		// create or update bundle rule
		_, err = service.BundleRuleSvc.CreateOrUpdateBundleRule(signerAddress, *params.Body.BucketName, *params.Body.MaxBundleFiles, *params.Body.MaxBundleSize, *params.Body.MaxFinalizeTime)
		if err != nil {
			return rule.NewSetBundleRuleBadRequest()
		}

		return rule.NewSetBundleRuleOK()
	}
}
