package types

import "github.com/node-real/greenfield-bundle-service/models"

const (
	DefaultMaxFileSize = 10 * 1024 * 1024 // 10MB

	DefaultMaxBundleFiles  = 100
	DefaultMaxBundleSize   = 1024 * 1024 * 1024 // 1GB
	DefaultMaxFinalizeTime = 60 * 60 * 24

	MaxBundleFiles  = 10000
	MaxBundleSize   = 2 * 1024 * 1024 * 1024 // 2GB
	MaxFinalizeTime = 60 * 60 * 24 * 7       // 7 days
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
	ErrorInvalidFileName = &models.Error{
		Code:    10008,
		Message: "Invalid file name",
	}
	ErrorInvalidContentType = &models.Error{
		Code:    10009,
		Message: "Invalid content type",
	}
	ErrorInvalidFileContent = &models.Error{
		Code:    10010,
		Message: "Invalid file content",
	}
	ErrorInvalidFileSha256 = &models.Error{
		Code:    10011,
		Message: "Invalid file sha256",
	}
	ErrorObjectExist = &models.Error{
		Code:    10012,
		Message: "Object exists",
	}
	ErrorBundleNotExist = &models.Error{
		Code:    10013,
		Message: "Bundle does not exist",
	}
	ErrorInvalidBundleStatus = &models.Error{
		Code:    10014,
		Message: "Invalid bundle status",
	}
	ErrorObjectAlreadyExists = &models.Error{
		Code:    10015,
		Message: "Object already exists",
	}
	ErrorInvalidTags = &models.Error{
		Code:    10016,
		Message: "Invalid tags",
	}
	ErrorBundleSizeExceedsLimit = &models.Error{
		Code:    10017,
		Message: "Bundle size or files exceeds limit",
	}
	ErrorInvalidBundleRuleParams = &models.Error{
		Code:    10018,
		Message: "Invalid bundle rule params",
	}
)

func InternalErrorWithError(err error) *models.Error {
	return &models.Error{
		Code:    500,
		Message: err.Error(),
	}
}

func InvalidFileContentErrorWithError(err error) *models.Error {
	return &models.Error{
		Code:    10010,
		Message: err.Error(),
	}
}
