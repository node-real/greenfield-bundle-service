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

	fileManager := storage.NewFileManager(config, objectDao, gnfdClient)
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
			bundles, err := b.bundleDao.GetBundlingBundles()
			if err != nil {
				util.Logger.Errorf("get time out bundling bundles failed, err=%v", err.Error())
				continue
			}

			cur := time.Now()
			for _, bundle := range bundles {
				if cur.Sub(bundle.CreatedAt).Seconds() < float64(bundle.MaxFinalizeTime) {
					continue
				}

				bundle.Status = database.BundleStatusFinalized
				_, err := b.bundleDao.UpdateBundle(*bundle)
				if err != nil {
					util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
					continue
				}
			}
		}
	}
}

func (b *Bundler) startSubmitLoops() {
	if len(b.config.BundleConfig.BundlerPrivateKeys) == 0 {
		util.Logger.Fatal("no bundler account available")
	}

	for _, privateKey := range b.config.BundleConfig.BundlerPrivateKeys {
		account, err := types.NewAccountFromPrivateKey("bundler-account", privateKey)
		if err != nil {
			util.Logger.Fatalf("create bundler account failed, err=%v", err.Error())
		}

		go b.submitLoop(account)
	}
}

func (b *Bundler) submitLoop(account *types.Account) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	accountAddr := account.GetAddress().String()
	client, err := client.New(b.config.GnfdConfig.ChainId, b.config.GnfdConfig.RpcUrl, client.Option{DefaultAccount: account})
	if err != nil {
		util.Logger.Fatalf("create greenfield client failed, account=%s, err=%v", accountAddr, err.Error())
	}

	for {
		select {
		case <-ticker.C:
			bundles, err := b.bundleDao.GetFinalizedBundlesByBundlerAccount(accountAddr)
			if err != nil {
				util.Logger.Errorf("get finalized bundles by bundler account failed, bundler=%s, err=%v", accountAddr, err.Error())
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
		objectReader, err := b.fileManager.GetObject(bundleRecord.Bucket, bundleRecord.Name, object.ObjectName)
		if err != nil {
			return nil, 0, fmt.Errorf("get object failed, object=%s, err=%v", object.ObjectName, err)
		}

		// TODO: Fix zero size
		objectMeta, err := newBundle.AppendObject(object.ObjectName, 0, objectReader, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("append object to bundle object failed, object=%s, err=%v", object.ObjectName, err)
		}

		object.OffsetInBundle = int64(objectMeta.Offset)
		_, err = b.objectDao.UpdateObject(*object)
		if err != nil {
			return nil, 0, fmt.Errorf("update object error, object=%+v, err=%s", object, err.Error())
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
