package util

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	Logger = logrus.StandardLogger()
)

// InitLogger initialises the logger.
func InitLogger(config *LogConfig) {

	logrus.SetReportCaller(true)

	var writers []io.Writer

	if config.UseConsoleLogger {
		writers = append(writers, os.Stdout)
	}

	if config.UseFileLogger {
		writers = append(writers, &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxFileSizeInMB,              // MaxSize is the maximum size in megabytes of the log file
			MaxBackups: config.MaxBackupsOfLogFiles,         // MaxBackups is the maximum number of old log files to retain
			MaxAge:     config.MaxAgeToRetainLogFilesInDays, // MaxAge is the maximum number of days to retain old log files
			Compress:   config.Compress,
		})
	}

	multiWriter := io.MultiWriter(writers...)
	Logger.SetOutput(multiWriter)
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		panic(err)
	}
	Logger.SetLevel(level)
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: false,
	})
}
