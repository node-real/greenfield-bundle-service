package e2e

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"

	"github.com/node-real/greenfield-bundle-service/bundler"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/storage"
	"github.com/node-real/greenfield-bundle-service/util"
)

var (
	objectContent = "Hello, Greenfield Bundle! Welcome to test and use it!"
)

func prepareBundleAndObject(db *gorm.DB, config *util.ServerConfig, client client.IClient, account *types.Account) (string, string, error) {
	ctx := context.Background()
	accInfo, err := client.GetAccount(context.Background(), account.GetAddress().String())
	if err != nil {
		return "", "", fmt.Errorf("get account failed: %v", err)
	}

	bucketName := "bundler-test"
	bundleName := fmt.Sprintf("bundle-%d", accInfo.GetSequence())
	objectName := "object1.txt"

	_, err = client.HeadBucket(ctx, bucketName)
	if err != nil {
		_, err = client.CreateBucket(ctx, bucketName, "0xbb0cac4668d970db8b51b6c1f80e6627f2bd3c02", types.CreateBucketOptions{Visibility: storageTypes.VISIBILITY_TYPE_PUBLIC_READ})
		if err != nil {
			return bucketName, bundleName, fmt.Errorf("create bucket %s failed: %v", bucketName, err)
		}
	}

	bundle := database.Bundle{
		Owner:           account.GetAddress().String(),
		Bucket:          bucketName,
		Name:            bundleName,
		BundlerAccount:  account.GetAddress().String(),
		Status:          database.BundleStatusBundling,
		MaxFiles:        100,
		MaxSize:         1024 * 1024 * 1024,
		MaxFinalizeTime: 10, // Seconds
	}

	result := db.Create(&bundle)
	if result.Error != nil {
		return bucketName, bundleName, fmt.Errorf("create bundle error, err=%s", result.Error.Error())
	}

	objectPath := storage.GetObjectPath(config.BundleConfig.LocalStoragePath, bucketName, bundleName, objectName)
	if err := os.MkdirAll(filepath.Dir(objectPath), os.ModePerm); err != nil {
		return bucketName, bundleName, fmt.Errorf("mkdir for object %s failed: %v", objectPath, err)
	}
	objectFile, err := os.Create(objectPath)
	if err != nil {
		return bucketName, bundleName, fmt.Errorf("create object %s failed: %v", objectPath, err)
	}
	size, err := objectFile.Write([]byte(objectContent))
	if err != nil {
		return bucketName, bundleName, fmt.Errorf("write object %s failed: %v", objectPath, err)
	}

	object := database.Object{
		Bucket:     bucketName,
		BundleName: bundleName,
		ObjectName: objectName,
		Owner:      account.GetAddress().String(),
		Size:       int64(size),
	}

	result = db.Create(&object)
	if result.Error != nil {
		return bucketName, bundleName, fmt.Errorf("create object error, err=%s", result.Error.Error())
	}

	return bucketName, bundleName, nil
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

	account, err := types.NewAccountFromPrivateKey("bundler-account", config.BundleConfig.BundlerPrivateKeys[0])
	if err != nil {
		util.Logger.Fatalf("create bundler account failed, err=%v", err.Error())
	}

	client, err := client.New(config.GnfdConfig.ChainId, config.GnfdConfig.RpcUrl, client.Option{DefaultAccount: account})
	if err != nil {
		util.Logger.Fatalf("create greenfield client failed, account=%s, err=%v", account.GetAddress().String(), err.Error())
	}

	bucketName, bundleName, err := prepareBundleAndObject(db, config, client, account)
	if err != nil {
		util.Logger.Fatalf("prepare failed, err=%v", err.Error())
	}

	bundler, err := bundler.NewBundler(config, db)
	go bundler.Run()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			objDetails, err := client.HeadObject(context.Background(), bucketName, bundleName)
			if err == nil {
				util.Logger.Infof("bundle on-chain success, name=%s, status=%v", objDetails.ObjectInfo.ObjectName, objDetails.ObjectInfo.ObjectStatus)
				if objDetails.ObjectInfo.ObjectStatus == storageTypes.OBJECT_STATUS_SEALED {
					_, _, err := client.GetObject(context.Background(), bucketName, bundleName, types.GetObjectOptions{})
					if err != nil {
						util.Logger.Fatalf("get bundled object failed, err=%v", err.Error())
					}

					return
				}
			}
		}
	}
}
