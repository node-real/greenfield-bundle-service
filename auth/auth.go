package auth

import (
	"context"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield/x/permission/types"
	"github.com/ethereum/go-ethereum/common"
)

type AuthManager struct {
	gnfdClient client.IClient
}

func NewAuthManager(gnfdClient client.IClient) *AuthManager {
	return &AuthManager{
		gnfdClient: gnfdClient,
	}
}

// IsBucketPermissionGranted check if the bucket permission is granted
func (a *AuthManager) IsBucketPermissionGranted(bundlerAddress common.Address, bucket string) (bool, error) {
	effect, err := a.gnfdClient.IsBucketPermissionAllowed(context.Background(), bundlerAddress.Hex(), bucket, types.ACTION_CREATE_OBJECT)
	if err != nil {
		return false, err
	}

	return effect == types.EFFECT_ALLOW, nil
}
