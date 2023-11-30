package database

import "time"

// UserBundlerAccount is used to store the relationship between user and bundler account
type UserBundlerAccount struct {
	Id             int64     `json:"id" gorm:"primaryKey"`
	UserAddress    string    `json:"user_address" gorm:"size:64;index:idx_user_bundler_account,priority:1,unique"`
	BundlerAddress string    `json:"bundler_address" gorm:"size:64;index:idx_user_bundler_account,priority:2,unique"`
	CreatedAt      time.Time `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
}
