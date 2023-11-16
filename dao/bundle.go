package dao

import "gorm.io/gorm"

type BundleDao interface {
}

type dbAccountDao struct {
	db *gorm.DB
}

// NewBundleDao returns a new BundleDao
func NewBundleDao(db *gorm.DB) BundleDao {
	return &dbAccountDao{
		db: db,
	}
}
