package dao

import (
	"errors"

	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

type BundlerAccountDao interface {
	GetBundlerAccountForUser(user string) (database.BundlerAccount, error)
	GetBundlerAccount(bundler string) (database.BundlerAccount, error)
	CreateBundlerAccount(bundlerAccount database.BundlerAccount) error
}

type dbBundlerAccountDao struct {
	db *gorm.DB
}

// NewBundlerAccountDao returns a new BundlerAccountDao
func NewBundlerAccountDao(db *gorm.DB) BundlerAccountDao {
	return &dbBundlerAccountDao{
		db: db,
	}
}

// GetBundlerAccount returns the bundler account for the specified bundler
func (s *dbBundlerAccountDao) GetBundlerAccount(bundler string) (database.BundlerAccount, error) {
	var bundlerAccount database.BundlerAccount
	err := s.db.Where("account_address = ?", bundler).Take(&bundlerAccount).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Errorf("get bundler account error, err=%s", err.Error())
		return database.BundlerAccount{}, err
	}
	return bundlerAccount, nil
}

// CreateBundlerAccount creates a new bundler account
func (s *dbBundlerAccountDao) CreateBundlerAccount(bundlerAccount database.BundlerAccount) error {
	err := s.db.Create(&bundlerAccount).Error
	if err != nil {
		util.Logger.Errorf("create bundler account error, err=%s", err.Error())
		return err
	}
	return nil
}

// GetBundlerAccountForUser returns the bundler account for the specified user
func (s *dbBundlerAccountDao) GetBundlerAccountForUser(user string) (database.BundlerAccount, error) {
	allBundlers, err := s.GetAllBundlerAccounts()
	if err != nil {
		util.Logger.Errorf("get all bundler accounts error, err=%s", err.Error())
		return database.BundlerAccount{}, err
	}

	bundlerForUser, err := types.PickBundlerIndexForAccount(len(allBundlers), user)
	if err != nil {
		util.Logger.Errorf("pick bundler index for account error, err=%s", err.Error())
		return database.BundlerAccount{}, err
	}
	return allBundlers[bundlerForUser], nil
}

func (s *dbBundlerAccountDao) GetAllBundlerAccounts() ([]database.BundlerAccount, error) {
	var bundlers []database.BundlerAccount
	err := s.db.Find(&bundlers).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return bundlers, nil
}
