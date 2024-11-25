package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes the Logrus logger
func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
