package dao

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/node-real/greenfield-bundle-service/util"

	"github.com/node-real/greenfield-bundle-service/database"
)

type BundleDao interface {
	CreateBundleIfNotBundlingExist(newBundle database.Bundle) (database.Bundle, error)
	QueryBundleWithMaxNonce(bucket string) (*database.Bundle, error)
	QueryBundle(bucket string, name string) (*database.Bundle, error)
	UpdateBundle(bundle database.Bundle) (*database.Bundle, error)
	GetBundlingBundle(bucket string) (database.Bundle, error)
	DeleteBundle(bucket string, name string) error
	GetBundlingBundles() ([]*database.Bundle, error)
	GetFinalizedBundlesByBundlerAccount(account string) ([]*database.Bundle, error)
	GetCreatedOnChainBundlesByBundlerAccount(account string) ([]*database.Bundle, error)
	InsertObjectsInOneTransaction(bundle database.Bundle, objects []database.Object) (database.Bundle, error)
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

// UpdateBundle updates a bundle
func (s *dbBundleDao) UpdateBundle(bundle database.Bundle) (*database.Bundle, error) {
	bundle.UpdatedAt = time.Now()
	err := s.db.Save(&bundle).Error
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

// DeleteBundle deletes a bundle and related objects
func (s *dbBundleDao) DeleteBundle(bucket string, name string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// todo(igor): This function can be optimized by a offline job if there is performance issue. we can mark the
		// bundle as deleted and delete the bundle and related objects in the offline job

		// Delete the Object records that have the specified bucket and bundleName
		if err := tx.Where("bucket = ? AND bundle_name = ?", bucket, name).Delete(&database.Object{}).Error; err != nil {
			return err
		}

		// Delete the Bundle record that has the specified bucket and name
		if err := tx.Where("bucket = ? AND name = ?", bucket, name).Delete(&database.Bundle{}).Error; err != nil {
			return err
		}

		return nil
	})
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
			util.Logger.Errorf("count bundle error, err=%s", err.Error())
			return err
		}

		// Check if a bundle with the same bucket and 'Bundling' status already exists
		if count > 0 {
			util.Logger.Errorf("a bundling bundle for the bucket already exists")
			// A bundle with the same bucket and 'Bundling' status already exists
			return errors.New("a bundling bundle for the bucket already exists")
		}

		newBundle.Status = database.BundleStatusBundling

		// No existing bundle found, safe to create a new one
		if err := tx.Create(&newBundle).Error; err != nil {
			util.Logger.Errorf("create bundle error, err=%s", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return database.Bundle{}, err
	}

	return newBundle, nil
}

func (s *dbBundleDao) GetBundlingBundles() ([]*database.Bundle, error) {
	var bundles []*database.Bundle
	err := s.db.Where("status = ?", database.BundleStatusBundling).Find(&bundles).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return bundles, nil
}

func (s *dbBundleDao) GetFinalizedBundlesByBundlerAccount(account string) ([]*database.Bundle, error) {
	var bundles []*database.Bundle
	err := s.db.Where("status = ? AND bundler_account = ?", database.BundleStatusFinalized, account).Find(&bundles).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return bundles, nil
}

func (s *dbBundleDao) GetCreatedOnChainBundlesByBundlerAccount(account string) ([]*database.Bundle, error) {
	var bundles []*database.Bundle
	err := s.db.Where("status = ? AND bundler_account = ?", database.BundleStatusCreatedOnChain, account).Find(&bundles).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return bundles, nil
}

// InsertObjectsInOneTransaction inserts objects in one transaction
func (s *dbBundleDao) InsertObjectsInOneTransaction(bundle database.Bundle, objects []database.Object) (database.Bundle, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Insert the bundle into the database
		if err := tx.Create(&bundle).Error; err != nil {
			return err
		}

		sql := "INSERT INTO objects (bucket, bundle_name, object_name, content_type, hash_algo, hash, owner, size, offset_in_bundle, tags, created_at, updated_at) VALUES "
		var values []interface{}

		for _, object := range objects {
			sql += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
			values = append(values, object.Bucket, object.BundleName, object.ObjectName, object.ContentType, object.HashAlgo, object.Hash, object.Owner, object.Size, object.OffsetInBundle, object.Tags, object.CreatedAt, object.UpdatedAt)
		}

		// trim the last ","
		sql = sql[0 : len(sql)-1]

		return tx.Exec(sql, values...).Error
	})

	if err != nil {
		return database.Bundle{}, err
	}

	return bundle, nil
}
