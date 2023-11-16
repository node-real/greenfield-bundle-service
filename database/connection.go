package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/node-real/greenfield-bundle-service/util"
)

func ConnectDBWithConfig(config *util.DBConfig) (*gorm.DB, error) {
	if config.DBDialect == "sqlite3" {
		db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
		return db.Debug(), err
	} else if config.DBDialect == "mysql" {
		dbPath := fmt.Sprintf("%s:%s@%s", config.Username, config.Password, config.DBPath)
		db, err := gorm.Open(mysql.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		dbConfig, err := db.DB()
		if err != nil {
			panic(err)
		}
		dbConfig.SetMaxIdleConns(config.MaxIdleConns)
		dbConfig.SetMaxOpenConns(config.MaxOpenConns)

		if err = db.AutoMigrate(&Bundle{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&Object{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&BundleRule{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&BundlerAccount{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&UserBundlerAccount{}); err != nil {
			panic(err)
		}
		return db.Debug(), nil
	} else {
		return nil, fmt.Errorf("dialect %s not supported", config.DBDialect)
	}
}
