package easymongo

import (
	"github.com/sirupsen/logrus"
)

// Logger represents a common interface that is consumed by easymongo. Add wrappers around
//     e.g. ConnectWith().Logger(Logger).Connect()
type Logger interface {
	// WithField(msg string, arg interface{}) Logger
	// WithFields(fields map[string]interface{}) Logger
	Debugf(format string, args ...interface{})
	// Debugf(format string, args ...interface{})
	// Info(args ...interface{})
	// Infof(format string, args ...interface{})
	// Warn(args ...interface{})
	// Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	// Errorf(format string, args ...interface{})
}

type DefaultLogger struct {
	logger *logrus.Entry
}

// NewDefaultLogger returns a DefaultLogger, which implements the Logger interface
func NewDefaultLogger() *DefaultLogger {
	l := logrus.New().WithField("src", "easymongo")
	l.Logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote:   true,
		DisableSorting: true,
	})
	l.Logger.SetLevel(logrus.DebugLevel)
	return &DefaultLogger{
		logger: l,
	}
}

func (logger *DefaultLogger) Debugf(format string, args ...interface{}) {
	logger.logger.Debugf(format, args...)
}

func (logger *DefaultLogger) Errorf(format string, args ...interface{}) {
	logger.logger.Errorf(format, args...)
}
