package database

import "time"

// BundleRule is used to store the bundle rule information
type BundleRule struct {
	Id              int64     `json:"id" gorm:"primaryKey"`
	Owner           string    `json:"owner" gorm:"size:64;index:idx_bundle_rule,priority:1,unique"`
	Bucket          string    `json:"bucket" gorm:"size:64;index:idx_bundle_rule,priority:2,unique"`
	MaxFiles        int64     `json:"max_files"`
	MaxSize         int64     `json:"max_size"`
	MaxFinalizeTime int64     `json:"max_finalize_time"`
	CreatedAt       time.Time `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
}
