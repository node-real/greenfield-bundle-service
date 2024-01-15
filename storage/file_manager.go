package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bnb-chain/greenfield-bundle-sdk/bundle"
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	LocalPathObjectPrefix = "object"
	LocalPathBundlePrefix = "bundle"
)

func GetObjectPath(storagePath, bucket, bundle, object string) string {
	return filepath.Join(storagePath, LocalPathObjectPrefix, bucket, bundle, object)
}
func GetBundlePath(storagePath, bucket, bundle string) string {
	return filepath.Join(storagePath, LocalPathBundlePrefix, bucket, bundle)
}

func GetObjectKeyInOss(bucket, bundle, object string) string {
	return fmt.Sprintf("%s/%s/%s", bucket, bundle, object)
}

func GetBundleKeyInOss(bucket string, bundle string) string {
	return fmt.Sprintf("%s/%s", bucket, bundle)
}

type FileManager struct {
	config          *util.ServerConfig
	useLocalStorage bool
	ossStore        *OssStore
	objectDao       dao.ObjectDao
	bundleDao       dao.BundleDao
	gnfdClient      client.IClient
}

func NewFileManager(config *util.ServerConfig, objectDao dao.ObjectDao, bundleDao dao.BundleDao, gnfdClient client.IClient) *FileManager {
	fileManager := &FileManager{
		config:     config,
		objectDao:  objectDao,
		bundleDao:  bundleDao,
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

// GetObject returns the object file
func (f *FileManager) GetObject(bucket string, bundle string, object string) (io.ReadCloser, error) {
	if f.useLocalStorage {
		return f.GetObjectLocal(bucket, bundle, object)
	}
	return f.GetObjectFromOss(bucket, bundle, object)
}

// GetBundle returns the bundle file
func (f *FileManager) GetBundle(bucket string, bundle string) (io.ReadCloser, error) {
	if f.useLocalStorage {
		return f.GetBundleLocal(bucket, bundle)
	}
	return f.GetBundleFromOss(bucket, bundle)
}

// GetBundleLocal returns the bundle file from local storage
func (f *FileManager) GetBundleLocal(bucket string, bundle string) (io.ReadCloser, error) {
	filePath := GetBundlePath(f.config.BundleConfig.LocalStoragePath, bucket, bundle)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetBundleFromOss returns the bundle file from oss
func (f *FileManager) GetBundleFromOss(bucket string, bundle string) (io.ReadCloser, error) {
	bundleKey := GetBundleKeyInOss(bucket, bundle)

	bundleFile, err := f.ossStore.GetObject(context.Background(), bundleKey, 0, 0)
	if err != nil {
		return nil, err
	}
	return bundleFile, nil
}

// GetObjectFromGnfdBundle returns the object file from gnfd
func (f *FileManager) GetObjectFromGnfdBundle(bucket string, bundle string, object string) (io.ReadCloser, error) {
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

	objectFile, _, err := f.gnfdClient.GetObject(context.Background(), bucket, bundle, getObjectOption)
	if err != nil {
		return nil, err
	}

	return objectFile, nil
}

// GetObjectFromOss returns the object file from oss
func (f *FileManager) GetObjectFromOss(bucket string, bundle string, object string) (io.ReadCloser, error) {
	objectKey := GetObjectKeyInOss(bucket, bundle, object)

	startTime := time.Now()
	objectFile, err := f.ossStore.GetObject(context.Background(), objectKey, 0, 0)
	if err != nil {
		if IsNoSuchKey(err) {
			queriedBundle, err := f.bundleDao.QueryBundle(bucket, bundle)
			if err != nil {
				return nil, err
			}
			if queriedBundle.Id == 0 {
				return nil, fmt.Errorf("bundle not found, bucket=%s, bundle=%s", bucket, bundle)
			}

			var objectFile io.ReadCloser

			startTime = time.Now()
			if queriedBundle.Status == database.BundleStatusFinalized {
				objectFile, err = f.GetObjectFromOssBundle(bucket, bundle, object)
				if err != nil {
					util.Logger.Errorf("failed to get object from oss bundle, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
					return nil, err
				}
				util.Logger.Infof("get object from oss bundle, bucket=%s, bundle=%s, object=%s, time=%s", bucket, bundle, object, time.Since(startTime).String())
			} else {
				objectFile, err = f.GetObjectFromGnfdBundle(bucket, bundle, object)
				if err != nil {
					util.Logger.Errorf("failed to get object from gnfd bundle, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
					return nil, err
				}
				util.Logger.Infof("get object from gnfd bundle, bucket=%s, bundle=%s, object=%s, time=%s", bucket, bundle, object, time.Since(startTime).String())
			}

			var buf bytes.Buffer
			tempFile := io.TeeReader(objectFile, &buf)

			// store object to oss
			err = f.ossStore.PutObject(context.Background(), objectKey, tempFile)
			if err != nil {
				util.Logger.Errorf("failed to store object to oss, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
				return nil, err
			}

			return io.NopCloser(&buf), nil
		}
		return nil, err
	}
	util.Logger.Infof("get object from oss, bucket=%s, bundle=%s, object=%s, time=%s", bucket, bundle, object, time.Since(startTime).String())
	return objectFile, nil
}

// GetObjectFromOssBundle returns the object file from oss bundle
func (f *FileManager) GetObjectFromOssBundle(bucket string, bundle string, object string) (io.ReadCloser, error) {
	// get object from database
	dbObject, err := f.objectDao.GetObject(bucket, bundle, object)
	if err != nil {
		return nil, err
	}
	if dbObject.Id == 0 {
		return nil, fmt.Errorf("object not found, bucket=%s, bundle=%s, object=%s", bucket, bundle, object)
	}

	bundleKey := GetBundleKeyInOss(bucket, bundle)

	objectFile, err := f.ossStore.GetObject(context.Background(), bundleKey, dbObject.OffsetInBundle, dbObject.Size)
	if err != nil {
		return nil, err
	}
	return objectFile, nil
}

// GetObjectFromLocalBundle returns the object file from local bundle
func (f *FileManager) GetObjectFromLocalBundle(bucket string, bundleName string, object string) (io.ReadCloser, error) {
	bundlePath := GetBundlePath(f.config.BundleConfig.LocalStoragePath, bucket, bundleName)
	if _, err := os.Stat(bundlePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("bundle file does not exist: %v", err)
	}

	bdl, err := bundle.NewBundleFromFile(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create bundle from file: %v", err)
	}

	obj, _, err := bdl.GetObject(object)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from bundle: %v", err)
	}

	return obj, nil
}

// GetObjectLocal returns the object file from local storage
func (f *FileManager) GetObjectLocal(bucket string, bundle string, object string) (io.ReadCloser, error) {
	filePath := GetObjectPath(f.config.BundleConfig.LocalStoragePath, bucket, bundle, object)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		queriedBundle, err := f.bundleDao.QueryBundle(bucket, bundle)
		if err != nil {
			return nil, err
		}
		if queriedBundle.Id == 0 {
			return nil, fmt.Errorf("bundle not found, bucket=%s, bundle=%s", bucket, bundle)
		}

		var objectFile io.ReadCloser
		if queriedBundle.Status == database.BundleStatusFinalized {
			objectFile, err = f.GetObjectFromLocalBundle(bucket, bundle, object)
			if err != nil {
				return nil, err
			}
		} else {
			objectFile, err = f.GetObjectFromGnfdBundle(bucket, bundle, object)
			if err != nil {
				return nil, err
			}
		}

		// store object to local
		_, _, err = f.storeLocalFile(filePath, objectFile)
		if err != nil {
			util.Logger.Errorf("failed to store object to local, bucket=%s, bundle=%s, object=%s, err=%s", bucket, bundle, object, err.Error())
			return nil, err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// StoreObject stores the object file
func (f *FileManager) StoreObject(bucket string, bundle string, object string, in io.ReadCloser) (string, int64, error) {
	if f.useLocalStorage {
		util.Logger.Infof("store object to local, bucket=%s, bundle=%s, object=%s", bucket, bundle, object)
		return f.StoreObjectLocal(bucket, bundle, object, in)
	}
	util.Logger.Infof("store object to oss, bucket=%s, bundle=%s, object=%s", bucket, bundle, object)
	return f.StoreObjectToOss(bucket, bundle, object, in)
}

// StoreBundle stores the bundle file
func (f *FileManager) StoreBundle(bucket string, bundle string, in io.ReadCloser) (string, int64, error) {
	if f.useLocalStorage {
		util.Logger.Infof("store bundle to local, bundle=%s", bundle)
		return f.StoreBundleLocal(bucket, bundle, in)
	}
	util.Logger.Infof("store bundle to oss, bundle=%s", bundle)
	return f.StoreBundleToOss(bucket, bundle, in)
}

// StoreBundleLocal stores the bundle file to local storage
func (f *FileManager) StoreBundleLocal(bucket, bundle string, in io.ReadCloser) (string, int64, error) {
	filePath := GetBundlePath(f.config.BundleConfig.LocalStoragePath, bucket, bundle)

	return f.storeLocalFile(filePath, in)
}

// StoreBundleToOss stores the bundle file to oss
func (f *FileManager) StoreBundleToOss(bucket string, bundle string, in io.ReadCloser) (string, int64, error) {
	bundleKey := GetBundleKeyInOss(bucket, bundle)

	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, in)
	if err != nil {
		return "", 0, err
	}

	err = f.ossStore.PutObject(context.Background(), bundleKey, buf)
	if err != nil {
		return "", 0, err
	}

	return bundleKey, size, nil
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
