package internal

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger returns a structured logrus entry logger
// with fields and configuration pulled from environment variables
func NewLogger() *logrus.Entry {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		if reqLogLevel, err := logrus.ParseLevel(level); err == nil {
			logger.SetLevel(reqLogLevel)
		}
	}
	hostname, _ := os.Hostname()
	environment := os.Getenv("ENVIRONMENT")

	return logger.WithFields(logrus.Fields{
		"host":        hostname,
		"environment": environment,
	})
}
