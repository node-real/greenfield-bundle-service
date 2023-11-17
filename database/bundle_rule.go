package database

import "time"

// BundleRule is used to store the bundle rule information
type BundleRule struct {
	Id              int64     `json:"id" gorm:"primaryKey"`
	Owner           string    `json:"owner" gorm:"64"`
	Bucket          string    `json:"bucket" gorm:"64"`
	MaxFiles        int64     `json:"max_files"`
	MaxSize         int64     `json:"max_size"`
	MaxFinalizeTime int64     `json:"max_finalize_time"`
	CreatedAt       time.Time `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
}
