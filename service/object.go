package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/util"
)

type Object interface {
	CreateObjectForBundling(newObject database.Object) (database.Object, error)
	StoreObjectFile(bundleName string, params bundle.UploadObjectParams) (string, int64, error)
}

type ObjectService struct {
	config         *util.ServerConfig
	userBundlerDao dao.UserBundlerAccountDao
	bundleDao      dao.BundleDao
	objectDao      dao.ObjectDao
}

// NewObjectService returns a new ObjectService
func NewObjectService(config *util.ServerConfig, bundleDao dao.BundleDao, objectDao dao.ObjectDao, userBundlerDao dao.UserBundlerAccountDao) Object {
	return &ObjectService{
		config:         config,
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
	fileName := fmt.Sprintf("%s-%s-%s", params.XBundleBucketName, bundleName, params.XBundleFileName)

	localFilePath := filepath.Join(s.config.BundleConfig.LocalStoragePath, fileName)
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return "", 0, err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, params.File)
	if err != nil {
		return "", 0, err
	}

	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		return "", 0, err
	}

	return localFilePath, fileInfo.Size(), nil
}
