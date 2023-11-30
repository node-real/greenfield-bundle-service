package service

import (
	"fmt"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	BundleNameFormat = "bundle-%d"
)

type Bundle interface {
	CreateBundle(newBundle database.Bundle) (database.Bundle, error)
	QueryBundle(bucketName string, bundleName string) (*database.Bundle, error)
	FinalizeBundle(bucketName string, bundleName string) (*database.Bundle, error)
	GetBundlingBundle(bucketName string) (database.Bundle, error)
}

type BundleService struct {
	bundleDao      dao.BundleDao
	bundleRuleDao  dao.BundleRuleDao
	userBundlerDao dao.UserBundlerAccountDao
}

// NewBundleService returns a new BundleService
func NewBundleService(bundleDao dao.BundleDao, bundleRuleDao dao.BundleRuleDao, userBundlerDao dao.UserBundlerAccountDao) Bundle {
	return &BundleService{
		bundleDao:      bundleDao,
		bundleRuleDao:  bundleRuleDao,
		userBundlerDao: userBundlerDao,
	}
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

// QueryBundle returns the bundle for the bucket if it exists
func (s *BundleService) QueryBundle(bucketName string, bundleName string) (*database.Bundle, error) {
	bundle, err := s.bundleDao.QueryBundle(bucketName, bundleName)
	if err != nil {
		util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", bucketName, bundleName, err.Error())
		return nil, err
	}

	return bundle, nil
}

// CreateBundle creates a new bundle for the bucket if it does not exist
func (s *BundleService) CreateBundle(newBundle database.Bundle) (database.Bundle, error) {
	bundleRule, err := s.bundleRuleDao.Get(newBundle.Owner, newBundle.Bucket)
	if err != nil {
		util.Logger.Errorf("get bundle rule error, owner=%s, bucket=%s, err=%s", newBundle.Owner, newBundle.Bucket, err.Error())
		return database.Bundle{}, err
	}

	// set bundle rule
	newBundle.MaxFiles = bundleRule.MaxFiles
	newBundle.MaxSize = bundleRule.MaxSize
	newBundle.MaxFinalizeTime = bundleRule.MaxFinalizeTime

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
