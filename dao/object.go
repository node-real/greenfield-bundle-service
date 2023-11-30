package dao

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/node-real/greenfield-bundle-service/database"
)

type ObjectDao interface {
	CreateObjectForBundling(object database.Object) (database.Object, error)
}

type dbObjectDao struct {
	db *gorm.DB
}

// NewObjectDao returns a new ObjectDao
func NewObjectDao(db *gorm.DB) ObjectDao {
	return &dbObjectDao{
		db: db,
	}
}

// CreateObjectForBundling creates a new object for bundling
func (s *dbObjectDao) CreateObjectForBundling(object database.Object) (database.Object, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// find and lock the bundle with the specified bucket name and status BundleStatusBundling
		var bundle database.Bundle
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("bucket = ? AND status = ?", object.Bucket, database.BundleStatusBundling).First(&bundle).Error; err != nil {
			return err
		}

		// check the bundle name, it should be the same as the bundle name in the object
		if bundle.Name != object.BundleName {
			return errors.New("bundle name mismatch")
		}

		return tx.Create(&object).Error
	})

	if err != nil {
		return database.Object{}, err
	}

	return object, nil
}
