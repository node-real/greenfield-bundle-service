package service

import (
	"context"
	"fmt"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	sdktypes "github.com/bnb-chain/greenfield-go-sdk/types"
	gnfdtypes "github.com/bnb-chain/greenfield/x/storage/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/node-real/greenfield-bundle-service/auth"
	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	BundleNamePrefix = "bundle-"
	BundleNameFormat = "bundle-%d"
)

type Bundle interface {
	CreateBundle(newBundle database.Bundle) (database.Bundle, error)
	QueryBundle(bucketName string, bundleName string) (*database.Bundle, error)
	FinalizeBundle(bucketName string, bundleName string) (*database.Bundle, error)
	GetBundlingBundle(bucketName string) (database.Bundle, error)
	QueryBucketFromGndf(bucketName string) (*gnfdtypes.BucketInfo, error)
	HeadObjectFromGnfd(bucketName string, objectName string) (*sdktypes.ObjectDetail, error)
	DeleteBundle(bucketName, bundleName string) error
}

type BundleService struct {
	gndfClient     client.IClient
	authManager    *auth.AuthManager
	bundleDao      dao.BundleDao
	bundleRuleDao  dao.BundleRuleDao
	userBundlerDao dao.UserBundlerAccountDao
}

// NewBundleService returns a new BundleService
func NewBundleService(gndfClient client.IClient, authManager *auth.AuthManager, bundleDao dao.BundleDao, bundleRuleDao dao.BundleRuleDao, userBundlerDao dao.UserBundlerAccountDao) Bundle {
	bs := BundleService{
		gndfClient:     gndfClient,
		authManager:    authManager,
		bundleDao:      bundleDao,
		bundleRuleDao:  bundleRuleDao,
		userBundlerDao: userBundlerDao,
	}
	return &bs
}

// GetBundlingBundle returns the bundling bundle for the bucket if it exists
func (s *BundleService) GetBundlingBundle(bucketName string) (database.Bundle, error) {
	bundle, err := s.bundleDao.GetBundlingBundle(bucketName)
	if err != nil {
		util.Logger.Errorf("get bundling bundle error, bucket=%s, err=%s", bucketName, err.Error())
		return database.Bundle{}, err
	}

	return bundle, nil
}

// QueryBucketFromGndf queries the bucket info from gndf
func (s *BundleService) QueryBucketFromGndf(bucketName string) (*gnfdtypes.BucketInfo, error) {
	bucket, err := s.gndfClient.HeadBucket(context.Background(), bucketName)
	if err != nil {
		util.Logger.Errorf("query bucket error, bucket=%s, err=%s", bucketName, err.Error())
		return nil, err
	}

	return bucket, nil
}

// DeleteBundle deletes the bundle for the bucket
func (s *BundleService) DeleteBundle(bucketName, bundleName string) error {
	err := s.bundleDao.DeleteBundle(bucketName, bundleName)
	if err != nil {
		util.Logger.Errorf("delete bundle error, bucket=%s, bundle=%s, err=%s", bucketName, bundleName, err.Error())
		return err
	}
	return nil
}

// HeadObjectFromGnfd queries the object info from gndf
func (s *BundleService) HeadObjectFromGnfd(bucketName string, objectName string) (*sdktypes.ObjectDetail, error) {
	object, err := s.gndfClient.HeadObject(context.Background(), bucketName, objectName)
	if err != nil {
		util.Logger.Errorf("query object error, bucket=%s, object=%s, err=%s", bucketName, objectName, err.Error())
		return nil, err
	}

	return object, nil
}

// QueryBundle returns the bundle for the bucket if it exists
func (s *BundleService) QueryBundle(bucketName string, bundleName string) (*database.Bundle, error) {
	bundle, err := s.bundleDao.QueryBundle(bucketName, bundleName)
	if err != nil {
		util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", bucketName, bundleName, err.Error())
		return nil, err
	}

	return bundle, nil
}

// CreateBundle creates a new bundle for the bucket if it does not exist, it also checks the permission for the bucket
// of the bundler
func (s *BundleService) CreateBundle(newBundle database.Bundle) (database.Bundle, error) {
	// check permission for the bucket
	isPermissionGranted, err := s.authManager.IsBucketPermissionGranted(common.HexToAddress(newBundle.BundlerAccount), newBundle.Bucket)
	if err != nil {
		util.Logger.Errorf("check bucket permission error, bucket=%s, err=%s", newBundle.Bucket, err.Error())
		return database.Bundle{}, err
	}

	if !isPermissionGranted {
		util.Logger.Errorf("bucket(%s) permission not granted for bundler(%s)", newBundle.Bucket, newBundle.BundlerAccount)
		return database.Bundle{}, fmt.Errorf("bucket(%s) permission not granted for bundler(%s)", newBundle.Bucket, newBundle.BundlerAccount)
	}

	// get bundle rule for the bucket
	bundleRule, err := s.bundleRuleDao.Get(newBundle.Owner, newBundle.Bucket)
	if err != nil {
		util.Logger.Errorf("get bundle rule error, owner=%s, bucket=%s, err=%s", newBundle.Owner, newBundle.Bucket, err.Error())
		return database.Bundle{}, err
	}

	// set bundle rule for the new bundle, if not exist, use default
	if bundleRule.Id == 0 {
		newBundle.MaxFiles = types.DefaultMaxBundleFiles
		newBundle.MaxSize = types.DefaultMaxBundleSize
		newBundle.MaxFinalizeTime = types.DefaultMaxFinalizeTime
	} else {
		newBundle.MaxFiles = bundleRule.MaxFiles
		newBundle.MaxSize = bundleRule.MaxSize
		newBundle.MaxFinalizeTime = bundleRule.MaxFinalizeTime
	}

	// set nonce
	previousBundle, err := s.bundleDao.QueryBundleWithMaxNonce(newBundle.Bucket)
	if err != nil {
		util.Logger.Errorf("get bundle with max nonce error, bucket=%s, err=%s", newBundle.Bucket, err.Error())
		return database.Bundle{}, err
	}
	if previousBundle == nil {
		newBundle.Nonce = 0
	} else {
		newBundle.Nonce = previousBundle.Nonce + 1
	}

	// set bundle name if not specified
	if newBundle.Name == "" {
		newBundle.Name = fmt.Sprintf(BundleNameFormat, newBundle.Nonce)
	}

	createdBundle, err := s.bundleDao.CreateBundleIfNotBundlingExist(newBundle)
	if err != nil {
		util.Logger.Errorf("create bundle error, bundle=%+v, err=%s", newBundle, err.Error())
		return database.Bundle{}, err
	}

	return createdBundle, nil
}

// FinalizeBundle finalizes the bundle for the bucket if it exists
func (s *BundleService) FinalizeBundle(bucketName string, bundleName string) (*database.Bundle, error) {
	bundle, err := s.bundleDao.QueryBundle(bucketName, bundleName)
	if err != nil {
		util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", bucketName, bundleName, err.Error())
		return nil, err
	}

	if bundle == nil {
		return nil, fmt.Errorf("bundle not found")
	}

	if bundle.Status != database.BundleStatusBundling {
		return nil, fmt.Errorf("bundle status is not bundling")
	}

	bundle.Status = database.BundleStatusFinalized

	updatedBundle, err := s.bundleDao.UpdateBundle(*bundle)
	if err != nil {
		util.Logger.Errorf("update bundle error, bundle=%+v, err=%s", bundle, err.Error())
		return nil, err
	}

	return updatedBundle, nil
}
