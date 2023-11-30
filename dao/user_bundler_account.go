package dao

import (
	"errors"

	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/database"
)

type UserBundlerAccountDao interface {
	GetUserBundlerAccount(user string) (database.UserBundlerAccount, error)
	CreateUserBundlerAccount(userBundlerAccount database.UserBundlerAccount) (database.UserBundlerAccount, error)
}

type dbUserBundlerAccountDao struct {
	db *gorm.DB
}

// NewUserBundlerAccountDao returns a new UserBundlerAccountDao
func NewUserBundlerAccountDao(db *gorm.DB) UserBundlerAccountDao {
	return &dbUserBundlerAccountDao{
		db: db,
	}
}

// GetUserBundlerAccount returns the bundler account for the specified user
func (s *dbUserBundlerAccountDao) GetUserBundlerAccount(user string) (database.UserBundlerAccount, error) {
	var userBundlerAccount database.UserBundlerAccount
	err := s.db.Where("user_address = ?", user).Take(&userBundlerAccount).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return userBundlerAccount, err
	}
	return userBundlerAccount, nil
}

// CreateUserBundlerAccount creates a new bundler account for the specified user
func (s *dbUserBundlerAccountDao) CreateUserBundlerAccount(userBundlerAccount database.UserBundlerAccount) (database.UserBundlerAccount, error) {
	if err := s.db.Create(&userBundlerAccount).Error; err != nil {
		return database.UserBundlerAccount{}, err
	}

	return userBundlerAccount, nil
}
