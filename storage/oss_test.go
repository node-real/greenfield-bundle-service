package storage

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOssStore_PutObject(t *testing.T) {
	bucketURL := "https://greenfield-stagenet-bundle-nodereal.oss-ap-northeast-1.aliyuncs.com"
	os.Setenv(OSSAccessId, "xxx")
	os.Setenv(OSSSecretKey, "xxx")

	store, err := NewOssStoreFromEnv(bucketURL)
	if err != nil {
		t.Fatalf("Failed to create OssStore: %v", err)
	}

	tempFile, err := os.CreateTemp("", "oss_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString("This is some test data")
	_, _ = tempFile.Seek(0, 0)

	err = store.PutObject(context.Background(), "testfile", tempFile)
	assert.NoError(t, err, "PutObject should not return an error")

	err = store.DeleteObject(context.Background(), GetObjectKeyInOss("bundle-test", "bundle-1", "test.html"))
	if err != nil {
		println(err.Error())
		return
	}

	println("delete success")

	// Get the object that was just put
	object, err := store.GetObject(context.Background(), "xxx", 0, 0)
	if err != nil {
		println("is no such key", IsNoSuchKey(err))

		t.Fatalf("Failed to get object: %s", err.Error())
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		t.Fatalf("Failed to read object content: %v", err)
	}
	t.Logf("Object content: %s", content)

}
