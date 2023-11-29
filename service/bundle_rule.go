package service

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
)

type BundleRule interface {
	CreateOrUpdateBundleRule(userAddress common.Address, bucketName string, maxFiles int64, maxSize int64, maxFinalizeTime int64) (database.BundleRule, error)
}

type BundleRuleService struct {
	bundleRuleDao dao.BundleRuleDao
}

// NewBundleRuleService returns a new BundleRuleService
func NewBundleRuleService(bundleRuleDao dao.BundleRuleDao) BundleRule {
	bs := BundleRuleService{
		bundleRuleDao: bundleRuleDao,
	}
	return &bs
}

func (s *BundleRuleService) QueryBundleRule(userAddress string, bucketName string) (database.BundleRule, error) {
	bundleRule, err := s.bundleRuleDao.Get(userAddress, bucketName)
	if err != nil {
		return database.BundleRule{}, err
	}

	return bundleRule, nil
}

// CreateOrUpdateBundleRule creates or updates bundle rule
func (s *BundleRuleService) CreateOrUpdateBundleRule(userAddress common.Address, bucketName string, maxFiles int64, maxSize int64, maxFinalizeTime int64) (database.BundleRule, error) {
	bundleRule, err := s.bundleRuleDao.Get(userAddress.String(), bucketName)
	if err != nil {
		return database.BundleRule{}, err
	}

	if bundleRule.Id == 0 {
		bundleRule, err = s.bundleRuleDao.Create(database.BundleRule{
			Owner:           userAddress.String(),
			Bucket:          bucketName,
			MaxFiles:        maxFiles,
			MaxSize:         maxSize,
			MaxFinalizeTime: maxFinalizeTime,
		})
		if err != nil {
			return database.BundleRule{}, err
		}
	} else {
		bundleRule, err = s.bundleRuleDao.Update(database.BundleRule{
			Id:              bundleRule.Id,
			Owner:           userAddress.String(),
			Bucket:          bucketName,
			MaxFiles:        maxFiles,
			MaxSize:         maxSize,
			MaxFinalizeTime: maxFinalizeTime,
		})
		if err != nil {
			return database.BundleRule{}, err
		}
	}

	return bundleRule, nil
}
