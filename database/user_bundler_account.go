package database

import "time"

// UserBundlerAccount is used to store the relationship between user and bundler account
type UserBundlerAccount struct {
	Id             int64     `json:"id" gorm:"primaryKey"`
	UserAddress    string    `json:"user_address" gorm:"64"`
	BundlerAddress string    `json:"account_address" gorm:"64"`
	CreatedAt      time.Time `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
}
