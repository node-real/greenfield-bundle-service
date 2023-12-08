package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/bnb-chain/greenfield-bundle-sdk/bundle"
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"

	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	LocalPathBundlePrefix = "bundle"
	LocalPathObjectPrefix = "object"
)

func GetBundlePath(storagePath, bucket, bundle string) string {
	return filepath.Join(storagePath, LocalPathBundlePrefix, bucket, bundle)
}

func GetObjectPath(storagePath, bucket, bundle, object string) string {
	return filepath.Join(storagePath, LocalPathObjectPrefix, bucket, bundle, object)
}

type FileManager struct {
	config     *util.ServerConfig
	gnfdClient client.IClient
}

func NewFileManager(config *util.ServerConfig, gnfdClient client.IClient) *FileManager {
	return &FileManager{
		config:     config,
		gnfdClient: gnfdClient,
	}
}

func (f *FileManager) GetObject(bucket string, bundle string, object string) (io.ReadCloser, int64, error) {
	return f.GetObjectLocal(bucket, bundle, object)
}

func (f *FileManager) GetObjectFromBundle(bucket string, bundleName string, object string) (io.ReadCloser, int64, error) {
	bundleFilePath := GetBundlePath(f.config.BundleConfig.LocalStoragePath, bucket, bundleName)

	// check if bundleName exists
	if _, err := os.Stat(bundleFilePath); os.IsNotExist(err) {
		// bundleName not exists, get from gnfd
		objectFile, _, err := f.gnfdClient.GetObject(context.Background(), bucket, bundleName, types.GetObjectOptions{})
		if err != nil {
			return nil, 0, err
		}

		bundleFilePath, _, err = f.StoreBundleFileLocal(bucket, bundleName, objectFile)
		if err != nil {
			return nil, 0, err
		}
	}

	bundleFile, err := bundle.NewBundleFromFile(bundleFilePath)
	if err != nil {
		return nil, 0, err
	}

	return bundleFile.GetObject(object)
}

// GetObjectLocal returns the object file from local storage
func (f *FileManager) GetObjectLocal(bucket string, bundle string, object string) (io.ReadCloser, int64, error) {
	filePath := GetObjectPath(f.config.BundleConfig.LocalStoragePath, bucket, bundle, object)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, 0, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	return file, fileInfo.Size(), nil
}

func (f *FileManager) StoreObject(bucket string, bundle string, object string, in io.ReadCloser) (string, int64, error) {
	return f.StoreObjectLocal(bucket, bundle, object, in)
}

// StoreObjectLocal stores the object file to local storage
func (f *FileManager) StoreObjectLocal(bucket string, bundle string, object string, in io.ReadCloser) (string, int64, error) {
	filePath := GetObjectPath(f.config.BundleConfig.LocalStoragePath, bucket, bundle, object)

	return f.storeLocalFile(filePath, in)
}

// StoreBundleFileLocal stores the bundle file to local storage
func (f *FileManager) StoreBundleFileLocal(bucket string, bundle string, in io.ReadCloser) (string, int64, error) {
	filePath := GetBundlePath(f.config.BundleConfig.LocalStoragePath, bucket, bundle)

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
