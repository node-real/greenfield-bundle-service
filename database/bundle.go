package database

import "time"

type BundleStatus uint

const (
	BundleStatusBundling       BundleStatus = 0
	BundleStatusBundled        BundleStatus = 1
	BundleStatusCreatedOnChain BundleStatus = 2
	BundleStatusDeleted        BundleStatus = 3
)

// Bundle is used to store the bundle information
type Bundle struct {
	Id             int64        `json:"id" gorm:"primaryKey"`
	Owner          string       `json:"owner" gorm:"64"`
	Bucket         string       `json:"bucket" gorm:"64"`
	Name           string       `json:"name" gorm:"size:128"`
	BundlerAccount string       `json:"bundler_account" gorm:"64"`
	Status         BundleStatus `json:"status"`
	Files          int64        `json:"files"`
	Sizes          int64        `json:"size"`
	MaxFiles       int64        `json:"max_files"`
	MaxSize        int64        `json:"max_size"`
	MaxBundleTime  int64        `json:"max_bundle_time"`
	Nonce          int64        `json:"nonce"`
	CreatedAt      time.Time    `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
}
