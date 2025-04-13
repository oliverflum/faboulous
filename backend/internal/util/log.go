package util

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
)

func GetLogLevel() log.Level {
	env := os.Getenv("FAB_ENV")
	if env == "development" {
		return log.LevelTrace
	} else if env == "test" {
		return log.LevelDebug
	} else {
		return log.LevelInfo
	}
}
