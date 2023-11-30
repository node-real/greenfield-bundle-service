package service

import (
	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

type Object interface {
	CreateObjectForBundling(newObject database.Object) (database.Object, error)
}

type ObjectService struct {
	userBundlerDao dao.UserBundlerAccountDao
	bundleDao      dao.BundleDao
	objectDao      dao.ObjectDao
}

// NewObjectService returns a new ObjectService
func NewObjectService(bundleDao dao.BundleDao, objectDao dao.ObjectDao, userBundlerDao dao.UserBundlerAccountDao) Object {
	return &ObjectService{
		bundleDao:      bundleDao,
		objectDao:      objectDao,
		userBundlerDao: userBundlerDao,
	}
}

// CreateObjectForBundling creates a new object for bundling
func (s *ObjectService) CreateObjectForBundling(newObject database.Object) (database.Object, error) {
	// create object
	object, err := s.objectDao.CreateObjectForBundling(newObject)
	if err != nil {
		util.Logger.Errorf("create object error, object=%+v, err=%s", newObject, err.Error())
		return database.Object{}, err
	}

	return object, nil
}
