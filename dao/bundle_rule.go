package dao

import (
	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/database"
)

type BundleRuleDao interface {
	Get(userAddress string, bucketName string) (database.BundleRule, error)
}

type dbBundleRuleDao struct {
	db *gorm.DB
}

// NewBundleRuleDao returns a new BundleRuleDao
func NewBundleRuleDao(db *gorm.DB) BundleRuleDao {
	return &dbBundleRuleDao{
		db: db,
	}
}

func (dao *dbBundleRuleDao) Get(userAddress string, bucketName string) (database.BundleRule, error) {
	var rule database.BundleRule
	err := dao.db.Where("owner = ? AND bucket = ?", userAddress, bucketName).Take(&rule).Error
	if err != nil {
		return rule, err
	}
	return rule, nil
}
