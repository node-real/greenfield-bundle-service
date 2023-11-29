package dao

import (
	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/database"
)

type BundleRuleDao interface {
	Get(userAddress string, bucketName string) (database.BundleRule, error)
	Create(rule database.BundleRule) (database.BundleRule, error)
	Update(rule database.BundleRule) (database.BundleRule, error)
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

// Get gets a bundle rule
func (dao *dbBundleRuleDao) Get(userAddress string, bucketName string) (database.BundleRule, error) {
	var rule database.BundleRule
	err := dao.db.Where("owner = ? AND bucket = ?", userAddress, bucketName).Take(&rule).Error
	if err != nil {
		return rule, err
	}
	return rule, nil
}

// Create creates a new bundle rule
func (dao *dbBundleRuleDao) Create(rule database.BundleRule) (database.BundleRule, error) {
	err := dao.db.Create(&rule).Error
	if err != nil {
		return rule, err
	}
	return rule, nil
}

// Update updates a bundle rule
func (dao *dbBundleRuleDao) Update(rule database.BundleRule) (database.BundleRule, error) {
	err := dao.db.Save(&rule).Error
	if err != nil {
		return rule, err
	}
	return rule, nil
}
