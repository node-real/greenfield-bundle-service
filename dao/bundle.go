package dao

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/node-real/greenfield-bundle-service/database"
)

type BundleDao interface {
	CreateBundleIfNotBundlingExist(newBundle database.Bundle) (database.Bundle, error)
	QueryBundleWithMaxNonce(bucket string) (*database.Bundle, error)
	QueryBundle(bucket string, name string) (*database.Bundle, error)
	UpdateBundle(bundle database.Bundle) (*database.Bundle, error)
	GetBundlingBundle(bucket string) (database.Bundle, error)
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

func (s *dbBundleDao) UpdateBundle(bundle database.Bundle) (*database.Bundle, error) {
	bundle.UpdatedAt = time.Now()
	err := s.db.Save(&bundle).Error
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

// GetBundlingBundle returns the bundling bundle for the bucket if it exists
func (s *dbBundleDao) GetBundlingBundle(bucket string) (database.Bundle, error) {
	var bundle database.Bundle
	err := s.db.Where("bucket = ? AND status = ?", bucket, database.BundleStatusBundling).Take(&bundle).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Bundle{}, err
	}
	return bundle, nil
}

func (s *dbBundleDao) QueryBundle(bucket string, name string) (*database.Bundle, error) {
	var bundle database.Bundle
	err := s.db.Where("bucket = ? AND name = ?", bucket, name).Take(&bundle).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &bundle, nil
}

func (s *dbBundleDao) QueryBundleWithMaxNonce(bucket string) (*database.Bundle, error) {
	var bundle database.Bundle
	err := s.db.Where("bucket = ?", bucket).Order("nonce desc").Take(&bundle).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &bundle, nil
}

func (s *dbBundleDao) CreateBundleIfNotBundlingExist(newBundle database.Bundle) (database.Bundle, error) {
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
		return database.Bundle{}, err
	}

	return newBundle, nil
}
