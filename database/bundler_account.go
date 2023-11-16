package database

import "time"

type BundleAccountStatus uint

const (
	BundleAccountStatusLiving BundleAccountStatus = 0
)

// BundlerAccount is used to store the bundler account information
type BundlerAccount struct {
	Id             int64               `json:"id" gorm:"primaryKey"`
	AccountAddress string              `json:"account_address" gorm:"64"`
	Status         BundleAccountStatus `json:"status"`
	CreatedAt      time.Time           `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt      time.Time           `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
}
