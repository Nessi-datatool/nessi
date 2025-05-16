package datalake

import (
	"io"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"github.com/sirupsen/logrus"
)

// initMockLogger initializes a mock logger for testing
func initMockLogger() {
	// Get the global logger instance
	logger := logging.GetLogger()
	
	// Configure it to discard all output
	logger.SetOutput(io.Discard)
	
	// Set level to panic to minimize any logging
	logger.SetLevel(logrus.PanicLevel)
}

