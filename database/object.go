package database

import (
	"time"

	"github.com/bnb-chain/greenfield-bundle-sdk/types"
)

// Object is used to store the object information
type Object struct {
	Id             int64          `json:"id" gorm:"primaryKey"`
	Bucket         string         `json:"bucket" gorm:"size:64;index:idx_object_name,priority:1,unique"`
	BundleName     string         `json:"bundle_name" gorm:"size:128;index:idx_object_name,priority:2,unique"`
	ObjectName     string         `json:"object_name" gorm:"size:1024;index:idx_object_name,priority:3,unique"`
	ContentType    string         `json:"content_type" gorm:"size:64"`
	HashAlgo       types.HashAlgo `json:"hash_algo"`
	Hash           []byte         `json:"hash"`
	Owner          string         `json:"owner" gorm:"size:64"`
	Size           int64          `json:"size"`
	OffsetInBundle int64          `json:"offset_in_bundle"`
	Tags           string         `json:"tags"`
	CreatedAt      time.Time      `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
}
