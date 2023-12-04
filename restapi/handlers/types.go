package handlers

import (
	"github.com/node-real/greenfield-bundle-service/models"
)

const (
	TimestampExpireTime = 60 * 5 // 5 min

	DefaultMaxBundleFiles  = 100
	DefaultMaxBundleSize   = 1024 * 1024 * 1024 // 1GB
	DefaultMaxFinalizeTime = 60 * 60 * 24       // 1 day
)

var (
	ErrorInternalError = &models.Error{
		Code:    500,
		Message: "Internal error",
	}

	ErrorInvalidExpiryTimestamp = &models.Error{
		Code:    10000,
		Message: "Invalid expiry timestamp",
	}
	ErrorInvalidSignature = &models.Error{
		Code:    10001,
		Message: "Invalid signature",
	}
	ErrorInvalidBucketName = &models.Error{
		Code:    10002,
		Message: "Invalid bucket name",
	}
	ErrorInvalidMaxBundleFiles = &models.Error{
		Code:    10003,
		Message: "Invalid max bundle files",
	}
	ErrorInvalidMaxBundleSize = &models.Error{
		Code:    10004,
		Message: "Invalid max bundle size",
	}
	ErrorInvalidMaxFinalizeTime = &models.Error{
		Code:    10005,
		Message: "Invalid max finalize time",
	}
	ErrorInvalidBundleName = &models.Error{
		Code:    10006,
		Message: "Invalid bundle name",
	}
	ErrorInvalidBundleOwner = &models.Error{
		Code:    10007,
		Message: "Invalid bundle owner",
	}
)
