package util

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
