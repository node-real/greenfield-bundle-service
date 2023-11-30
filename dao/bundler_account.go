package dao

import (
	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/util"

	"github.com/node-real/greenfield-bundle-service/database"
)

type BundlerAccountDao interface {
	GetBundlerAccountForUser(user string) (database.BundlerAccount, error)
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

// GetBundlerAccountForUser returns the bundler account for the specified user
func (s *dbBundlerAccountDao) GetBundlerAccountForUser(user string) (database.BundlerAccount, error) {
	// todo: get bundler account for user
	_, account, err := util.GenerateRandomAccount()
	if err != nil {
		util.Logger.Error("generate random account error, err=%s", err.Error())
		return database.BundlerAccount{}, err
	}

	return database.BundlerAccount{
		Id:             1,
		AccountAddress: account.String(),
	}, nil
}
