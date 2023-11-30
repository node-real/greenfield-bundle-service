package handlers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/models"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/service"
)

func HandleGetUserBundlerAccount() func(params bundle.BundlerAccountParams) middleware.Responder {
	return func(params bundle.BundlerAccountParams) middleware.Responder {
		userAddress := common.HexToAddress(params.UserAddress)
		bundlerAccount, err := service.UserBundlerAccountSvc.GetOrCreateUserBundlerAccount(userAddress.String())
		if err != nil {
			return bundle.NewBundlerAccountInternalServerError()
		}

		return bundle.NewBundlerAccountOK().WithPayload(&models.BundlerAccount{
			Address: bundlerAccount.BundlerAddress,
		})
	}
}
