package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	"gorm.io/gorm"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/pkg/utils"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	types2 "github.com/bnb-chain/greenfield/sdk/types"
	permTypes "github.com/bnb-chain/greenfield/x/permission/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"

	"github.com/node-real/greenfield-bundle-service/bundler"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/storage"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	ownerPrivateKey = "" // Set your own private key for testing
	bucketName      = "gnfd-bundler-test"
	objectContent   = "Hello, Greenfield Bundle! Welcome to test and use it!"
)

var (
	ownerAcc   *types.Account
	bundlerAcc *types.Account
	gnfdClient client.IClient
)

func setupAccounts(config *util.ServerConfig) {
	var err error

	ownerAcc, err = types.NewAccountFromPrivateKey("owner-account", ownerPrivateKey)
	if err != nil {
		util.Logger.Fatalf("create owner account failed, err=%v", err.Error())
	}

	bundlerAcc, err = types.NewAccountFromPrivateKey("bundler-account", config.BundleConfig.BundlerPrivateKeys[0])
	if err != nil {
		util.Logger.Fatalf("create bundler account failed, err=%v", err.Error())
	}

	gnfdClient, err = client.New(config.GnfdConfig.ChainId, config.GnfdConfig.RpcUrl, client.Option{DefaultAccount: ownerAcc})
	if err != nil {
		util.Logger.Fatalf("create greenfield client failed, account=%s, err=%v", ownerAcc.GetAddress().String(), err.Error())
	}
}

func setupBucket() {
	ctx := context.Background()
	_, err := gnfdClient.HeadBucket(ctx, bucketName)
	if err == nil {
		return
	}

	_, err = gnfdClient.CreateBucket(ctx, bucketName, "0xbb0cac4668d970db8b51b6c1f80e6627f2bd3c02", types.CreateBucketOptions{Visibility: storageTypes.VISIBILITY_TYPE_PUBLIC_READ})
	if err != nil {
		util.Logger.Fatalf("create bucket %s failed: %v", bucketName, err)
	}
}

func grantFeesAndPermission() {
	ctx := context.Background()

	bucketActions := []permTypes.ActionType{permTypes.ACTION_CREATE_OBJECT}
	statements := utils.NewStatement(bucketActions, permTypes.EFFECT_ALLOW, nil, types.NewStatementOptions{})
	principal, err := utils.NewPrincipalWithAccount(bundlerAcc.GetAddress())
	if err != nil {
		util.Logger.Fatalf("fail to generate marshaled principal: %v", err)
	}
	txHash, err := gnfdClient.PutBucketPolicy(ctx, bucketName, principal, []*permTypes.Statement{&statements}, types.PutPolicyOption{})
	if err != nil {
		util.Logger.Fatalf("put policy failed: %v", err)
	}
	_, err = gnfdClient.WaitForTx(ctx, txHash)
	if err != nil {
		util.Logger.Fatalf("wait for grant permission tx failed: %v", err)
	}

	allowanceAmount := math.NewIntWithDecimal(1, 18)
	allowance, err := gnfdClient.QueryBasicAllowance(ctx, ownerAcc.GetAddress().String(), bundlerAcc.GetAddress().String())
	if err == nil {
		for _, coin := range allowance.SpendLimit {
			if coin.Denom == types2.Denom && coin.Amount.GTE(allowanceAmount) {
				return
			}
		}
		txHash, err = gnfdClient.RevokeAllowance(ctx, bundlerAcc.GetAddress().String(), types2.TxOption{})
		if err != nil {
			util.Logger.Warnf("revoke fee allowance failed: %v", err)
		}
		_, err = gnfdClient.WaitForTx(ctx, txHash)
		if err != nil {
			util.Logger.Warnf("wait for revoke allowance tx failed: %v", err)
		}
	}

	txHash, err = gnfdClient.GrantBasicAllowance(ctx, bundlerAcc.GetAddress().String(), allowanceAmount, nil, types2.TxOption{})
	if err != nil {
		util.Logger.Fatalf("grant fee allowance failed: %v", err)
	}
	_, err = gnfdClient.WaitForTx(ctx, txHash)
	if err != nil {
		util.Logger.Fatalf("wait for grant fee tx failed: %v", err)
	}
}

func setupDatabaseRecords(db *gorm.DB, config *util.ServerConfig) string {
	ctx := context.Background()
	accInfo, err := gnfdClient.GetAccount(ctx, ownerAcc.GetAddress().String())
	if err != nil {
		util.Logger.Fatalf("get owner account failed: %v", err)
	}

	bundleName := fmt.Sprintf("bundle-%d", accInfo.GetSequence())
	objectName := "test-object.txt"
	bundle := database.Bundle{
		Owner:           ownerAcc.GetAddress().String(),
		Bucket:          bucketName,
		Name:            bundleName,
		BundlerAccount:  bundlerAcc.GetAddress().String(),
		Status:          database.BundleStatusBundling,
		MaxFiles:        100,
		MaxSize:         1024 * 1024 * 1024,
		MaxFinalizeTime: 10, // Seconds
	}
	result := db.Create(&bundle)
	if result.Error != nil {
		util.Logger.Fatalf("create bundle error, err=%s", result.Error.Error())
	}

	objectPath := storage.GetObjectPath(config.BundleConfig.LocalStoragePath, bucketName, bundleName, objectName)
	err = os.MkdirAll(filepath.Dir(objectPath), os.ModePerm)
	if err != nil {
		util.Logger.Fatalf("mkdir for object %s failed: %v", objectPath, err)
	}
	objectFile, err := os.Create(objectPath)
	if err != nil {
		util.Logger.Fatalf("create object %s failed: %v", objectPath, err)
	}
	size, err := objectFile.Write([]byte(objectContent))
	if err != nil {
		util.Logger.Fatalf("write object %s failed: %v", objectPath, err)
	}

	object := database.Object{
		Bucket:     bucketName,
		BundleName: bundleName,
		ObjectName: objectName,
		Owner:      ownerAcc.GetAddress().String(),
		Size:       int64(size),
	}
	result = db.Create(&object)
	if result.Error != nil {
		util.Logger.Fatalf("create object error, err=%s", result.Error.Error())
	}

	return bundleName
}

func TestBundler(t *testing.T) {
	config := util.ParseServerConfigFromFile("../config/server/dev.json")

	util.InitLogger(config.LogConfig)

	db, err := database.ConnectDBWithConfig(config.DBConfig)
	if err != nil {
		util.Logger.Fatalf("connect database error, err=%s", err.Error())
	}

	if len(config.BundleConfig.BundlerPrivateKeys) == 0 {
		util.Logger.Fatalf("no bundler account available")
	}

	setupAccounts(config)
	setupBucket()
	grantFeesAndPermission()
	bundleName := setupDatabaseRecords(db, config)

	bundler, err := bundler.NewBundler(config, db)
	go bundler.Run()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		objDetails, err := gnfdClient.HeadObject(context.Background(), bucketName, bundleName)
		if err != nil {
			continue
		}

		util.Logger.Infof("bundle on-chain success, name=%s, status=%v", objDetails.ObjectInfo.ObjectName, objDetails.ObjectInfo.ObjectStatus)
		if objDetails.ObjectInfo.ObjectStatus == storageTypes.OBJECT_STATUS_SEALED {
			_, _, err := gnfdClient.GetObject(context.Background(), bucketName, bundleName, types.GetObjectOptions{})
			if err != nil {
				util.Logger.Fatalf("get bundled object failed, err=%v", err.Error())
			}

			return
		}
	}
}
