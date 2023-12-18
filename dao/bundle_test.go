package dao_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

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
