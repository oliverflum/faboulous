package util

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm/logger"
)

func GetLogLevelFiber() log.Level {
	env := os.Getenv("FAB_ENV")
	if env == "development" {
		return log.LevelTrace
	} else if env == "test" {
		return log.LevelDebug
	} else {
		return log.LevelInfo
	}
}

func GetLogLevelGorm() logger.LogLevel {
	env := os.Getenv("FAB_ENV")
	if env == "development" {
		return logger.Info
	} else if env == "test" {
		return logger.Warn
	} else {
		return logger.Error
	}
}
