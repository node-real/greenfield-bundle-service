package dao

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/node-real/greenfield-bundle-service/database"
)

type BundleDao interface {
	CreateBundleIfNotBundlingExist(newBundle database.Bundle) (*database.Bundle, error)
	QueryBundleWithMaxNonce(bucket string) (*database.Bundle, error)
}

type dbBundleDao struct {
	db *gorm.DB
}

// NewBundleDao returns a new BundleDao
func NewBundleDao(db *gorm.DB) BundleDao {
	return &dbBundleDao{
		db: db,
	}
}

func (s *dbBundleDao) QueryBundleWithMaxNonce(bucket string) (*database.Bundle, error) {
	var bundle database.Bundle
	err := s.db.Where("bucket = ?", bucket).Order("nonce desc").Take(&bundle).Error
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (s *dbBundleDao) CreateBundleIfNotBundlingExist(newBundle database.Bundle) (*database.Bundle, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var count int64

		// Lock the rows with a "SELECT FOR UPDATE" to prevent concurrent inserts
		// for the same bucket and status.
		if err := tx.Model(&database.Bundle{}).
			Where("bucket = ? AND status = ?", newBundle.Bucket, database.BundleStatusBundling).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Count(&count).Error; err != nil {
			return err
		}

		// Check if a bundle with the same bucket and 'Bundling' status already exists
		if count > 0 {
			// A bundle with the same bucket and 'Bundling' status already exists
			return errors.New("a bundling bundle for the bucket already exists")
		}

		newBundle.Status = database.BundleStatusBundling

		// No existing bundle found, safe to create a new one
		if err := tx.Create(&newBundle).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newBundle, nil
}
