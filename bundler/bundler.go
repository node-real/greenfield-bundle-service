package bundler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/bnb-chain/greenfield-bundle-sdk/bundle"
	bundleTypes "github.com/bnb-chain/greenfield-bundle-sdk/types"
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	gnfdsdktypes "github.com/bnb-chain/greenfield/sdk/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/storage"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	MaxSealOnChainTime = 60 * 60 * 24 // 1 day
	EmptyErrMessage    = ""
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
	b.finalizeLoop()
}

func (b *Bundler) finalizeLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		bundles, err := b.bundleDao.GetBundlingBundles()
		if err != nil {
			util.Logger.Errorf("get time out bundling bundles failed, err=%v", err.Error())
			continue
		}

		for _, bundle := range bundles {
			if bundle.Size >= bundle.MaxSize || bundle.Files >= bundle.MaxFiles ||
				time.Since(bundle.CreatedAt).Seconds() >= float64(bundle.MaxFinalizeTime) {
				bundle.Status = database.BundleStatusFinalized

				// mark the objects as expired if the bundle is empty
				if bundle.Files == 0 {
					bundle.Status = database.BundleStatusExpired
				}

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

		// register bundler account
		b.registerBundler(account)

		go b.submitLoop(account)
	}
}

func (b *Bundler) registerBundler(account *types.Account) {
	accountAddr := account.GetAddress().String()
	bundlerAccount, err := b.bundlerAccountDao.GetBundlerAccount(accountAddr)
	if err != nil {
		util.Logger.Errorf("get bundler account failed, bundler=%s, err=%v", accountAddr, err.Error())
		return
	}

	if bundlerAccount.Id != 0 {
		return
	}

	bundlerAccount = database.BundlerAccount{
		AccountAddress: accountAddr,
	}
	err = b.bundlerAccountDao.CreateBundlerAccount(bundlerAccount)
	if err != nil {
		util.Logger.Errorf("create bundler account failed, bundler=%s, err=%v", accountAddr, err.Error())
		return
	}
}

func (b *Bundler) submitLoop(account *types.Account) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	sealTicker := time.NewTicker(30 * time.Second)
	defer sealTicker.Stop()

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
				if !bundle.IsTimeToRetry() {
					continue
				}

				bundledObject, size, err := b.assembleBundleObject(bundle)
				if err != nil {
					util.Logger.Errorf("assemble bundle object failed, bundle=%s, err=%v", bundle.Bucket+bundle.Name, err.Error())
					bundle.RetryCounter++
					bundle.ErrMessage = fmt.Sprintf("assemble bundle failed: %v", err)
					_, err = b.bundleDao.UpdateBundle(*bundle)
					if err != nil {
						util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
					}
					continue
				}

				txHash, objectDetail, err := b.submitBundledObject(client, bundle, bundledObject, size)
				if err != nil {
					util.Logger.Errorf("submit bundle object failed, bundle=%s, err=%v", bundle.Bucket+bundle.Name, err.Error())
					bundle.RetryCounter++
					bundle.ErrMessage = fmt.Sprintf("submit bundle failed: %v", err)
					_, err = b.bundleDao.UpdateBundle(*bundle)
					if err != nil {
						util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
					}
					continue
				}

				bundle.Status = database.BundleStatusCreatedOnChain
				bundle.TxHash = txHash
				bundle.ObjectId = objectDetail.ObjectInfo.Id.Uint64()
				bundle.RetryCounter = 0
				bundle.ErrMessage = EmptyErrMessage
				_, err = b.bundleDao.UpdateBundle(*bundle)
				if err != nil {
					util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
				}
			}

		case <-sealTicker.C:
			bundles, err := b.bundleDao.GetCreatedOnChainBundlesByBundlerAccount(accountAddr)
			if err != nil {
				util.Logger.Errorf("get created onchain bundles by bundler account failed, bundler=%s, err=%v", accountAddr, err.Error())
				continue
			}

			for _, bundle := range bundles {
				sealed := b.checkBundleSealed(client, bundle)
				if sealed {
					bundle.Status = database.BundleStatusSealedOnChain
					_, err = b.bundleDao.UpdateBundle(*bundle)
					if err != nil {
						util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
					}
					continue
				}

				if bundle.RetryCounter == 0 && time.Since(bundle.UpdatedAt).Seconds() < MaxSealOnChainTime {
					continue
				}

				// Cancel create for sealing timeout bundled object.
				if !bundle.IsTimeToRetry() {
					continue
				}
				err = b.cancelCreateBundle(client, bundle)
				if err != nil {
					util.Logger.Errorf("cancel create timeout bundle error, bundle=%+v, err=%s", bundle, err.Error())
					bundle.RetryCounter++
					bundle.ErrMessage = fmt.Sprintf("seal timeout, but cancel bundle failed: %v", err)
					_, err = b.bundleDao.UpdateBundle(*bundle)
					if err != nil {
						util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
					}
					continue
				}

				// Change the bundle status to "finalized" to trigger resubmission.
				bundle.Status = database.BundleStatusFinalized
				bundle.RetryCounter = 0
				bundle.ErrMessage = EmptyErrMessage
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

		var tags map[string]string
		err = json.Unmarshal([]byte(object.Tags), &tags)
		if err != nil {
			util.Logger.Warnf("unmarshal tags failed, tags=%s, err=%v", object.Tags, err.Error())
			tags = nil
		}
		objectMeta, err := newBundle.AppendObject(object.ObjectName, objectReader, &bundleTypes.AppendObjectOptions{
			HashAlgo:    object.HashAlgo,
			Hash:        object.Hash,
			ContentType: object.ContentType,
			Tags:        tags,
		})
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

func (b *Bundler) submitBundledObject(client client.IClient, bundle *database.Bundle, object io.ReadSeekCloser, size int64) (string, *types.ObjectDetail, error) {
	if size == 0 {
		return "", nil, fmt.Errorf("invalid bundle size")
	}

	owner, err := sdk.AccAddressFromHexUnsafe(bundle.Owner)
	if err != nil {
		return "", nil, fmt.Errorf("invalid owner address, owner=%s, err=%v", bundle.Owner, err)
	}

	var txHash string
	objectDetail, err := client.HeadObject(context.Background(), bundle.Bucket, bundle.Name)
	if err != nil {
		opts := types.CreateObjectOptions{
			Visibility:  storageTypes.VISIBILITY_TYPE_PUBLIC_READ,
			ContentType: "bundle",
			TxOpts:      &gnfdsdktypes.TxOption{FeeGranter: owner},
		}

		txHash, err = client.CreateObject(context.Background(), bundle.Bucket, bundle.Name, object, opts)
		if err != nil {
			return "", nil, fmt.Errorf("create bundle object failed, bucket=%s, bundle=%s, err=%v", bundle.Bucket, bundle.Name, err)
		}

		objectDetail, err = client.HeadObject(context.Background(), bundle.Bucket, bundle.Name)
		if err != nil {
			return "", nil, fmt.Errorf("head bundle object failed, bucket=%s, bundle=%s, err=%v", bundle.Bucket, bundle.Name, err)
		}

		_, _ = object.Seek(0, 0)
	}

	opts := types.PutObjectOptions{
		ContentType: "bundle",
	}
	return txHash, objectDetail, client.PutObject(context.Background(), bundle.Bucket, bundle.Name, size, object, opts)
}

func (b *Bundler) checkBundleSealed(client client.IClient, bundle *database.Bundle) bool {
	objectDetail, err := client.HeadObjectByID(context.Background(), strconv.FormatUint(bundle.ObjectId, 10))
	if err != nil {
		util.Logger.Errorf("head bundle object failed, bundle=%s, objectId = %d, err=%v", bundle.Bucket+bundle.Name, bundle.ObjectId, err)
		return false
	}

	return objectDetail.ObjectInfo.ObjectStatus == storageTypes.OBJECT_STATUS_SEALED
}

func (b *Bundler) cancelCreateBundle(client client.IClient, bundle *database.Bundle) error {
	owner, err := sdk.AccAddressFromHexUnsafe(bundle.Owner)
	if err != nil {
		return fmt.Errorf("invalid owner address, owner=%s, err=%v", bundle.Owner, err)
	}

	_, err = client.CancelCreateObject(context.Background(), bundle.Bucket, bundle.Name, types.CancelCreateOption{
		TxOpts: &gnfdsdktypes.TxOption{FeeGranter: owner}},
	)

	return err
}
