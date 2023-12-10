package bundler

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"io"
	"time"

	"github.com/bnb-chain/greenfield-bundle-sdk/bundle"
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/storage"
	"github.com/node-real/greenfield-bundle-service/util"
)

type Bundler struct {
	config *util.ServerConfig

	objectDao         dao.ObjectDao
	bundleDao         dao.BundleDao
	bundlerAccountDao dao.BundlerAccountDao
	fileManager       *storage.FileManager
}

func NewBundler(config *util.ServerConfig, db *gorm.DB) (*Bundler, error) {
	objectDao := dao.NewObjectDao(db)
	bundleDao := dao.NewBundleDao(db)
	bundlerAccountDao := dao.NewBundlerAccountDao(db)

	gnfdClient, err := client.New(config.GnfdConfig.ChainId, config.GnfdConfig.RpcUrl, client.Option{})
	if err != nil {
		util.Logger.Fatalf("unable to new greenfield client, %v", err)
	}

	fileManager := storage.NewFileManager(config, gnfdClient)
	return &Bundler{
		config:            config,
		objectDao:         objectDao,
		bundleDao:         bundleDao,
		bundlerAccountDao: bundlerAccountDao,
		fileManager:       fileManager,
	}, nil
}

func (b *Bundler) Run() {
	b.startSubmitLoops()
	b.timeOutLoop()
}

func (b *Bundler) timeOutLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bundles, err := b.bundleDao.GetTimeOutBundlingBundles()
			if err != nil {
				util.Logger.Errorf("get time out bundling bundles failed, err=%v", err.Error())
				continue
			}

			// TODO: Batch update
			for _, bundle := range bundles {
				bundle.Status = database.BundleStatusFinalized
				_, err := b.bundleDao.UpdateBundle(*bundle)
				if err != nil {
					util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
				}
			}
		}
	}
}

func (b *Bundler) startSubmitLoops() {
	bundlerAccounts, err := b.bundlerAccountDao.GetAllBundlerAccounts()
	if err != nil {
		util.Logger.Fatalf("get bundler accounts failed, err=%s", err.Error())
	}

	for _, account := range bundlerAccounts {
		go b.submitLoop(account)
	}
}

func (b *Bundler) submitLoop(bundler database.BundlerAccount) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// TODO: Get account private key
	privateKey := ""
	account, err := types.NewAccountFromPrivateKey("bundler-account", privateKey)
	if err != nil {
		util.Logger.Fatalf("create bundler account failed, account=%s, err=%v", bundler.AccountAddress, err.Error())
	}

	client, err := client.New(b.config.GnfdConfig.ChainId, b.config.GnfdConfig.RpcUrl, client.Option{DefaultAccount: account})
	if err != nil {
		util.Logger.Fatalf("create bundler account failed, account=%s, err=%v", bundler.AccountAddress, err.Error())
	}

	for {
		select {
		case <-ticker.C:
			bundles, err := b.bundleDao.GetFinalizedBundlesByBundlerAccount(bundler.AccountAddress)
			if err != nil {
				util.Logger.Errorf("get finalized bundles by bundler account failed, bundler=%s, err=%v", bundler.AccountAddress, err.Error())
				continue
			}

			for _, bundle := range bundles {
				bundledObject, size, err := b.assembleBundleObject(bundle)
				if err != nil {
					util.Logger.Errorf("assemble bundle object failed, bundle=%s, err=%v", bundle.Bucket+bundle.Name, err.Error())
					continue
				}

				err = b.submitBundledObject(client, bundle, bundledObject, size)
				if err != nil {
					util.Logger.Errorf("submit bundle object failed, bundle=%s, err=%v", bundle.Bucket+bundle.Name, err.Error())
					continue
				}

				bundle.Status = database.BundleStatusCreatedOnChain
				_, err = b.bundleDao.UpdateBundle(*bundle)
				if err != nil {
					util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
				}
			}
		}
	}
}

func (b *Bundler) assembleBundleObject(bundleRecord *database.Bundle) (io.ReadSeekCloser, int64, error) {
	newBundle, err := bundle.NewBundle()
	if err != nil {
		return nil, 0, fmt.Errorf("new bundle failed: %v", err)
	}

	objects, err := b.objectDao.GetBundleObjects(bundleRecord.Bucket, bundleRecord.Name)
	if err != nil {
		return nil, 0, fmt.Errorf("get bundle objects failed: %v", err)
	}

	for _, object := range objects {
		objectReader, size, err := b.fileManager.GetObject(bundleRecord.Bucket, bundleRecord.Name, object.ObjectName)
		if err != nil {
			return nil, 0, fmt.Errorf("get object failed, object=%s, err=%v", object.ObjectName, err)
		}

		_, err = newBundle.AppendObject(object.ObjectName, size, objectReader, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("append object to bundle object failed, object=%s, err=%v", object.ObjectName, err)
		}
	}

	bundledObject, size, err := newBundle.FinalizeBundle()
	if err != nil {
		return nil, 0, fmt.Errorf("finalize bundle failed, err=%v", err)
	}

	return bundledObject, size, nil
}

func (b *Bundler) submitBundledObject(client client.IClient, bundle *database.Bundle, object io.ReadSeekCloser, size int64) error {
	if size == 0 {
		return fmt.Errorf("invalid bundle size")
	}

	_, err := client.HeadObject(context.Background(), bundle.Bucket, bundle.Name)
	if err != nil {
		// Create object on chain
		opts := types.CreateObjectOptions{
			Visibility:  storageTypes.VISIBILITY_TYPE_PUBLIC_READ,
			ContentType: "bundle",
		}

		_, err := client.CreateObject(context.Background(), bundle.Bucket, bundle.Name, object, opts)
		if err != nil {
			return fmt.Errorf("create bundle object failed, bucket=%s, bundle=%s, err=%v", bundle.Bucket, bundle.Name, err)
		}

		object.Seek(0, 0)
	}

	opts := types.PutObjectOptions{
		ContentType: "bundle",
	}
	return client.PutObject(context.Background(), bundle.Bucket, bundle.Name, size, object, opts)
}
