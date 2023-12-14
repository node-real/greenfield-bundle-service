package util

import (
	"encoding/json"
	"os"
)

type BundleConfig struct {
	LocalStoragePath string `json:"local_storage_path"`
	OssBucketUrl     string `json:"oss_bucket_url"`
}

type GnfdConfig struct {
	ChainId string `json:"chain_id"`
	RpcUrl  string `json:"rpc_url"`
}

type DBConfig struct {
	DBDialect     string `json:"db_dialect"`
	DBPath        string `json:"db_path"`
	Password      string `json:"password"`
	Username      string `json:"username"`
	MaxIdleConns  int    `json:"max_idle_conns"`
	MaxOpenConns  int    `json:"max_open_conns"`
	AWSRegion     string `json:"aws_region"`
	AWSSecretName string `json:"aws_secret_name"`
}

type LogConfig struct {
	Level                        string `json:"level"`
	Filename                     string `json:"filename"`
	MaxFileSizeInMB              int    `json:"max_file_size_in_mb"`
	MaxBackupsOfLogFiles         int    `json:"max_backups_of_log_files"`
	MaxAgeToRetainLogFilesInDays int    `json:"max_age_to_retain_log_files_in_days"`
	UseConsoleLogger             bool   `json:"use_console_logger"`
	UseFileLogger                bool   `json:"use_file_logger"`
	Compress                     bool   `json:"compress"`
}

type ServerConfig struct {
	DBConfig     *DBConfig     `json:"db_config"`
	BundleConfig *BundleConfig `json:"bundle_config"`
	GnfdConfig   *GnfdConfig   `json:"gnfd_config"`
	LogConfig    *LogConfig    `json:"log_config"`
}

func ParseServerConfigFromFile(filePath string) *ServerConfig {
	bz, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var config ServerConfig
	if err := json.Unmarshal(bz, &config); err != nil {
		panic(err)
	}

	if config.DBConfig.Username == "" || config.DBConfig.Password == "" { // read password from AWS secret
		config.DBConfig.Username, config.DBConfig.Password = GetDBUsernamePassword(config.DBConfig)
	}

	return &config
}

func GetDBUsernamePassword(cfg *DBConfig) (string, string) {
	result, err := GetSecret(cfg.AWSSecretName, cfg.AWSRegion)
	if err != nil {
		panic(err)
	}
	type DBPass struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var dbPassword DBPass
	err = json.Unmarshal([]byte(result), &dbPassword)
	if err != nil {
		panic(err)
	}
	return dbPassword.Username, dbPassword.Password
}
