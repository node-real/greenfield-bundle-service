package service

import (
	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

type UserBundlerAccount interface {
	GetOrCreateUserBundlerAccount(user string) (database.UserBundlerAccount, error)
}

type UserBundlerAccountService struct {
	userBundlerDao    dao.UserBundlerAccountDao
	bundlerAccountDao dao.BundlerAccountDao
}

// NewUserBundlerAccountService returns a new UserBundlerAccountService
func NewUserBundlerAccountService(userBundlerDao dao.UserBundlerAccountDao, bundlerAccountDao dao.BundlerAccountDao) UserBundlerAccount {
	return &UserBundlerAccountService{
		userBundlerDao:    userBundlerDao,
		bundlerAccountDao: bundlerAccountDao,
	}
}

// GetOrCreateUserBundlerAccount returns the bundler account for the specified user, if not exists, create one
func (s *UserBundlerAccountService) GetOrCreateUserBundlerAccount(user string) (database.UserBundlerAccount, error) {
	userBundlerAccount, err := s.userBundlerDao.GetUserBundlerAccount(user)
	if err != nil {
		util.Logger.Errorf("get user bundler account error, user=%s, err=%s", user, err.Error())
		return database.UserBundlerAccount{}, err
	}

	// if the user bundler account exists, return it
	if userBundlerAccount.Id != 0 {
		return userBundlerAccount, nil
	}

	// get bundler account for the user
	bundlerAccount, err := s.bundlerAccountDao.GetBundlerAccountForUser(user)
	if err != nil {
		util.Logger.Errorf("get bundler account for user error, user=%s, err=%s", user, err.Error())
		return database.UserBundlerAccount{}, err
	}
	if bundlerAccount.AccountAddress == "" {
		util.Logger.Errorf("bundler account not found for user, user=%s", user)
		return database.UserBundlerAccount{}, err
	}

	// create user bundler account
	userBundlerAccount, err = s.userBundlerDao.CreateUserBundlerAccount(database.UserBundlerAccount{
		UserAddress:    user,
		BundlerAddress: bundlerAccount.AccountAddress,
	})
	if err != nil {
		return database.UserBundlerAccount{}, err
	}

	return userBundlerAccount, nil
}
