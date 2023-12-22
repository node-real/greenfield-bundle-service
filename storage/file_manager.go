package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	LocalPathObjectPrefix = "object"
)

func GetObjectPath(storagePath, bucket, bundle, object string) string {
	return filepath.Join(storagePath, LocalPathObjectPrefix, bucket, bundle, object)
}

func GetObjectKeyInOss(bucket, bundle, object string) string {
	return fmt.Sprintf("%s/%s/%s", bucket, bundle, object)
}

type FileManager struct {
	config          *util.ServerConfig
	useLocalStorage bool
	ossStore        *OssStore
	objectDao       dao.ObjectDao
	gnfdClient      client.IClient
}

func NewFileManager(config *util.ServerConfig, objectDao dao.ObjectDao, gnfdClient client.IClient) *FileManager {
	fileManager := &FileManager{
		config:     config,
		objectDao:  objectDao,
		gnfdClient: gnfdClient,
	}

	if config.BundleConfig.LocalStoragePath != "" {
		fileManager.useLocalStorage = true
	}

	if config.BundleConfig.OssBucketUrl != "" {
		util.Logger.Infof("use oss storage, bucketUrl=%s", config.BundleConfig.OssBucketUrl)
		fileManager.useLocalStorage = false
		ossStore, err := NewOssStoreFromEnv(config.BundleConfig.OssBucketUrl)
		if err != nil {
			panic(err)
		}
		fileManager.ossStore = ossStore
	}

	util.Logger.Infof("use local storage %v", fileManager.useLocalStorage)

	return fileManager
}

func (f *FileManager) GetObject(bucket string, bundle string, object string) (io.ReadCloser, error) {
	if f.useLocalStorage {
		return f.GetObjectLocal(bucket, bundle, object)
	}
	return f.GetObjectFromOss(bucket, bundle, object)
}

func (f *FileManager) GetObjectFromGnfd(bucket string, bundle string, object string) (io.ReadCloser, error) {
	// get object from database
	dbObject, err := f.objectDao.GetObject(bucket, bundle, object)
	if err != nil {
		return nil, err
	}
	if dbObject.Id == 0 {
		return nil, fmt.Errorf("object not found, bucket=%s, bundle=%s, object=%s", bucket, bundle, object)
	}
	// query object from gnfd
	getObjectOption := types.GetObjectOptions{}
	err = getObjectOption.SetRange(dbObject.OffsetInBundle, dbObject.OffsetInBundle+dbObject.Size-1) // [start, end]
	if err != nil {
		util.Logger.Errorf("failed to set range for object, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
		return nil, err
	}

	objectFile, _, err := f.gnfdClient.GetObject(context.Background(), bucket, bundle, types.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return objectFile, nil
}

func (f *FileManager) GetObjectFromOss(bucket string, bundle string, object string) (io.ReadCloser, error) {
	objectKey := GetObjectKeyInOss(bucket, bundle, object)

	objectFile, err := f.ossStore.GetObject(context.Background(), objectKey, 0, 0)
	if err != nil {
		if IsNoSuchKey(err) {
			objectFile, err := f.GetObjectFromGnfd(bucket, bundle, object)
			if err != nil {
				return nil, err
			}

			// store object to oss
			err = f.ossStore.PutObject(context.Background(), objectKey, objectFile)
			if err != nil {
				util.Logger.Errorf("failed to store object to oss, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
				return nil, err
			}

			return objectFile, nil
		}
		return nil, err
	}
	return objectFile, nil
}

// GetObjectLocal returns the object file from local storage
func (f *FileManager) GetObjectLocal(bucket string, bundle string, object string) (io.ReadCloser, error) {
	filePath := GetObjectPath(f.config.BundleConfig.LocalStoragePath, bucket, bundle, object)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		objectFile, err := f.GetObjectFromGnfd(bucket, bundle, object)
		if err != nil {
			return nil, err
		}

		// store object to local
		_, _, err = f.storeLocalFile(filePath, objectFile)
		if err != nil {
			util.Logger.Errorf("failed to store object to local, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
			return nil, err
		}
		return objectFile, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *FileManager) StoreObject(bucket string, bundle string, object string, in io.ReadCloser) (string, int64, error) {
	if f.useLocalStorage {
		util.Logger.Infof("store object to local, bucket=%s, bundle=%s, object=%s", bucket, bundle, object)
		return f.StoreObjectLocal(bucket, bundle, object, in)
	}
	util.Logger.Infof("store object to oss, bucket=%s, bundle=%s, object=%s", bucket, bundle, object)
	return f.StoreObjectToOss(bucket, bundle, object, in)
}

// StoreObjectToOss stores the object file to oss
func (f *FileManager) StoreObjectToOss(bucket string, bundle string, object string, in io.ReadCloser) (string, int64, error) {
	objectKey := GetObjectKeyInOss(bucket, bundle, object)

	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, in)
	if err != nil {
		return "", 0, err
	}

	err = f.ossStore.PutObject(context.Background(), objectKey, buf)
	if err != nil {
		return "", 0, err
	}

	return objectKey, size, nil
}

// StoreObjectLocal stores the object file to local storage
func (f *FileManager) StoreObjectLocal(bucket string, bundle string, object string, in io.ReadCloser) (string, int64, error) {
	filePath := GetObjectPath(f.config.BundleConfig.LocalStoragePath, bucket, bundle, object)

	return f.storeLocalFile(filePath, in)
}

// storeLocalFile stores the file to local storage
func (f *FileManager) storeLocalFile(filePath string, in io.ReadCloser) (string, int64, error) {
	if !f.localFileExists(filePath) {
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return "", 0, err
		}

		file, err := os.Create(filePath)
		if err != nil {
			return "", 0, err
		}
		defer file.Close()

		size, err := io.Copy(file, in)
		if err != nil {
			return "", 0, err
		}

		util.Logger.Infof("file stored to local, path=%s, size=%d", filePath, size)
		if err := in.Close(); err != nil {
			// don't return error here
			util.Logger.Errorf("failed to close file: %v", err)
		}
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", 0, err
	}

	return filePath, fileInfo.Size(), nil
}

func (f *FileManager) localFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
