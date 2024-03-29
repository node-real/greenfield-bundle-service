package database

import "time"

type BundleStatus uint

const (
	BundleStatusBundling       BundleStatus = 0
	BundleStatusFinalized      BundleStatus = 1
	BundleStatusCreatedOnChain BundleStatus = 2
	BundleStatusSealedOnChain  BundleStatus = 3
	BundleStatusExpired        BundleStatus = 4
)

var (
	maxRetryInterval = 2 * time.Hour
	retryIntervals   = []time.Duration{time.Minute, 10 * time.Minute, 30 * time.Minute, time.Hour, maxRetryInterval}
)

// Bundle is used to store the bundle information
type Bundle struct {
	Id              int64        `json:"id" gorm:"primaryKey"`
	Owner           string       `json:"owner" gorm:"size:64"`
	Bucket          string       `json:"bucket" gorm:"size:64;index:idx_bundle_name,priority:1,unique"`
	Name            string       `json:"name" gorm:"size:128;index:idx_bundle_name,priority:2,unique"`
	BundlerAccount  string       `json:"bundler_account" gorm:"size:64"`
	Status          BundleStatus `json:"status"`
	Files           int64        `json:"files"`
	Size            int64        `json:"size"`
	MaxFiles        int64        `json:"max_files"`
	MaxSize         int64        `json:"max_size"`
	MaxFinalizeTime int64        `json:"max_finalize_time"`
	Nonce           int64        `json:"nonce"`     // nonce is used to generate bundle name for auto generated bundle
	ObjectId        uint64       `json:"object_id"` // object_id is used to record the bundled object id on Greenfield
	TxHash          string       `json:"tx_hash"`   // tx_hash is used to record the tx hash on Greenfield
	RetryCounter    int          `json:"retry_counter"`
	ErrMessage      string       `json:"err_message"`
	CreatedAt       time.Time    `json:"created_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"NOT NULL;type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
}

func (b *Bundle) IsTimeToRetry() bool {
	if b.RetryCounter == 0 {
		return true
	}

	index := b.RetryCounter - 1
	if index >= len(retryIntervals) {
		index = len(retryIntervals) - 1
	}

	interval := retryIntervals[index]
	return time.Since(b.UpdatedAt) >= interval
}
