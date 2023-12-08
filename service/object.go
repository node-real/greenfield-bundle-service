package service

import (
	"io"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/storage"
	"github.com/node-real/greenfield-bundle-service/util"
)

type Object interface {
	CreateObjectForBundling(newObject database.Object) (database.Object, error)
	StoreObjectFile(bundleName string, params bundle.UploadObjectParams) (string, int64, error)
	GetObject(bucket string, bundle string, object string) (database.Object, error)
	GetObjectFile(bucket string, bundle string, object string) (io.ReadCloser, int64, error)
}

type ObjectService struct {
	config         *util.ServerConfig
	fileManager    *storage.FileManager
	userBundlerDao dao.UserBundlerAccountDao
	bundleDao      dao.BundleDao
	objectDao      dao.ObjectDao
}

// NewObjectService returns a new ObjectService
func NewObjectService(config *util.ServerConfig, fileManager *storage.FileManager, bundleDao dao.BundleDao, objectDao dao.ObjectDao, userBundlerDao dao.UserBundlerAccountDao) Object {
	return &ObjectService{
		config:         config,
		fileManager:    fileManager,
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

// StoreObjectFile stores the object file to local storage
func (s *ObjectService) StoreObjectFile(bundleName string, params bundle.UploadObjectParams) (string, int64, error) {
	return s.fileManager.StoreObjectLocal(params.XBundleBucketName, bundleName, params.XBundleFileName, params.File)
}

// GetObjectFile gets the object file
func (s *ObjectService) GetObjectFile(bucket string, bundle string, object string) (io.ReadCloser, int64, error) {
	return s.fileManager.GetObject(bucket, bundle, object)
}

// GetObject gets an object from database
func (s *ObjectService) GetObject(bucket string, bundle string, object string) (database.Object, error) {
	return s.objectDao.GetObject(bucket, bundle, object)
}
