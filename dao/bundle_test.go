package dao_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/util"
)

func TestCreateBundleIfNotBundlingExist_Concurrent(t *testing.T) {
	config := &util.DBConfig{
		DBDialect: "mysql",
		DBPath:    "tcp(localhost:3306)/test",
		Username:  "root",
		Password:  "12345678",
	}

	db, err := database.ConnectDBWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	// Empty the tables
	db.Exec("DELETE FROM bundles")
	db.Exec("DELETE FROM objects")
	db.Exec("DELETE FROM bundle_rules")
	db.Exec("DELETE FROM bundler_accounts")
	db.Exec("DELETE FROM user_bundler_accounts")

	bundleDao := dao.NewBundleDao(db)

	bundle := database.Bundle{
		Bucket: "testBucket",
		Name:   "testBundle",
	}

	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err1 = bundleDao.CreateBundleIfNotBundlingExist(bundle)
	}()

	go func() {
		defer wg.Done()
		_, err2 = bundleDao.CreateBundleIfNotBundlingExist(bundle)
	}()

	wg.Wait()

	assert.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil))
}

func TestInsertObjects(t *testing.T) {
	config := &util.DBConfig{
		DBDialect: "mysql",
		DBPath:    "tcp(localhost:3306)/test",
		Username:  "root",
		Password:  "12345678",
	}

	db, err := database.ConnectDBWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	// Empty the tables
	db.Exec("DELETE FROM objects")

	startTime := time.Now()

	err = db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < 1000; i++ {
			idx := i
			object := database.Object{
				Bucket:         "testBucket",
				BundleName:     "testBundle",
				ObjectName:     "testObject" + strconv.Itoa(idx),
				ContentType:    "application/octet-stream",
				HashAlgo:       1,
				Hash:           []byte("hash" + strconv.Itoa(idx)),
				Owner:          "owner" + strconv.Itoa(idx),
				Size:           int64(idx),
				OffsetInBundle: int64(idx),
				Tags:           "tag" + strconv.Itoa(idx),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := tx.Create(&object).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to insert objects: %v", err)
	}

	endTime := time.Now()

	timeCost := endTime.Sub(startTime).Milliseconds()

	t.Logf("Time cost: %d ms", timeCost)

	assert.LessOrEqual(t, timeCost, int64(10000), "Inserting 1000 objects should take less than 10000 milliseconds")

}

func TestInsertObjectsInOne(t *testing.T) {
	config := &util.DBConfig{
		DBDialect: "mysql",
		DBPath:    "tcp(localhost:3306)/test",
		Username:  "root",
		Password:  "12345678",
	}

	db, err := database.ConnectDBWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	// Empty the tables
	db.Exec("DELETE FROM objects")

	var objects []database.Object

	for i := 0; i < 1000; i++ {
		object := database.Object{
			Bucket:         "testBucket",
			BundleName:     "testBundle",
			ObjectName:     "testObject" + strconv.Itoa(i),
			ContentType:    "application/octet-stream",
			HashAlgo:       1,
			Hash:           []byte("hash" + strconv.Itoa(i)),
			Owner:          "owner" + strconv.Itoa(i),
			Size:           int64(i),
			OffsetInBundle: int64(i),
			Tags:           "tag" + strconv.Itoa(i),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		objects = append(objects, object)
	}

	startTime := time.Now()

	err = db.Transaction(func(tx *gorm.DB) error {
		sql := "INSERT INTO objects (id, bucket, bundle_name, object_name, content_type, hash_algo, hash, owner, size, offset_in_bundle, tags, created_at, updated_at) VALUES "
		values := []interface{}{}

		for _, object := range objects {
			sql += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
			values = append(values, object.Id, object.Bucket, object.BundleName, object.ObjectName, object.ContentType, object.HashAlgo, object.Hash, object.Owner, object.Size, object.OffsetInBundle, object.Tags, object.CreatedAt, object.UpdatedAt)
		}

		// trim the last ","
		sql = sql[0 : len(sql)-1]

		return tx.Exec(sql, values...).Error
	})

	if err != nil {
		t.Fatalf("Failed to insert objects: %v", err)
	}

	endTime := time.Now()

	timeCost := endTime.Sub(startTime).Milliseconds()

	t.Logf("Time cost: %d ms", timeCost)

	assert.LessOrEqual(t, timeCost, int64(10000), "Inserting 1000 objects should take less than 10000 milliseconds")
}
