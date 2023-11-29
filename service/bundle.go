package service

import (
	"fmt"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

type Bundle interface {
	CreateBundle(newBundle database.Bundle) (*database.Bundle, error)
	QueryBundle(bucketName string, bundleName string) (*database.Bundle, error)
	FinalizeBundle(bucketName string, bundleName string) (*database.Bundle, error)
}

type BundleService struct {
	bundleDao     dao.BundleDao
	bundleRuleDao dao.BundleRuleDao
}

func NewBundleService(bundleDao dao.BundleDao, bundleRuleDao dao.BundleRuleDao) Bundle {
	bs := BundleService{
		bundleDao:     bundleDao,
		bundleRuleDao: bundleRuleDao,
	}
	return &bs
}

func (s *BundleService) QueryBundle(bucketName string, bundleName string) (*database.Bundle, error) {
	bundle, err := s.bundleDao.QueryBundle(bucketName, bundleName)
	if err != nil {
		util.Logger.Errorf("query bundle error, bucket=%s, bundle=%s, err=%s", bucketName, bundleName, err.Error())
		return nil, err
	}

	return bundle, nil
}

func (s *BundleService) CreateBundle(newBundle database.Bundle) (*database.Bundle, error) {
	bundleRule, err := s.bundleRuleDao.Get(newBundle.Owner, newBundle.Bucket)
	if err != nil {
		util.Logger.Errorf("get bundle rule error, owner=%s, bucket=%s, err=%s", newBundle.Owner, newBundle.Bucket, err.Error())
		return nil, err
	}

	// set bundle rule
	newBundle.MaxFiles = bundleRule.MaxFiles
	newBundle.MaxSize = bundleRule.MaxSize
	newBundle.MaxFinalizeTime = bundleRule.MaxFinalizeTime

	// set nonce
	previousBundle, err := s.bundleDao.QueryBundleWithMaxNonce(newBundle.Bucket)
	if err != nil {
		util.Logger.Errorf("get bundle with max nonce error, bucket=%s, err=%s", newBundle.Bucket, err.Error())
		return nil, err
	}
	if previousBundle == nil {
		newBundle.Nonce = 0
	} else {
		newBundle.Nonce = previousBundle.Nonce + 1
	}

	createdBundle, err := s.bundleDao.CreateBundleIfNotBundlingExist(newBundle)
	if err != nil {
		util.Logger.Errorf("create bundle error, bundle=%+v, err=%s", newBundle, err.Error())
		return nil, err
	}

	return createdBundle, nil
}

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
