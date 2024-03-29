package dao

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/node-real/greenfield-bundle-service/database"
)

type ObjectDao interface {
	CreateObjectForBundling(object database.Object) (database.Object, error)
	UpdateObject(object database.Object) (*database.Object, error)
	GetObject(bucket string, bundle string, object string) (database.Object, error)
	GetBundleObjects(bucket string, bundle string) ([]*database.Object, error)
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

func (s *dbObjectDao) UpdateObject(object database.Object) (*database.Object, error) {
	err := s.db.Save(&object).Error
	if err != nil {
		return nil, err
	}
	return &object, nil
}

// GetObject gets an object
func (s *dbObjectDao) GetObject(bucket string, bundle string, object string) (database.Object, error) {
	var obj database.Object
	err := s.db.Where("bucket = ? AND bundle_name = ? AND object_name = ?", bucket, bundle, object).First(&obj).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return obj, err
	}
	return obj, nil
}

func (s *dbObjectDao) GetBundleObjects(bucket string, bundle string) ([]*database.Object, error) {
	var objs []*database.Object
	err := s.db.Where("bucket = ? AND bundle_name = ?", bucket, bundle).Find(&objs).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return objs, nil
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

		// update Files and Size in the bundle
		bundle.Files++
		bundle.Size += object.Size

		// save the updated bundle
		if err := tx.Save(&bundle).Error; err != nil {
			return err
		}

		return tx.Create(&object).Error
	})

	if err != nil {
		return database.Object{}, err
	}

	return object, nil
}
