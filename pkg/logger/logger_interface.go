package logger

import "github.com/nessi-dev/nessi/pkg/errorcode"

// These methods implement the errorcode.Logger interface

// SimpleDebug logs a debug message with the simple interface required by errorcode.Logger
func (l Logger) SimpleDebug(msg string) {
	l.Debug(msg)
}

// SimpleInfo logs an info message with the simple interface required by errorcode.Logger
func (l Logger) SimpleInfo(msg string) {
	l.Info(msg)
}

// SimpleWarn logs a warning message with the simple interface required by errorcode.Logger
func (l Logger) SimpleWarn(msg string) {
	l.Warn(msg)
}

// SimpleError logs an error message with the simple interface required by errorcode.Logger
func (l Logger) SimpleError(msg string) {
	l.Logger.Error().Msg(msg)
}

// SimpleWithField returns a new logger with a field added with the simple interface required by errorcode.Logger
func (l Logger) SimpleWithField(key string, value interface{}) errorcode.Logger {
	return Logger{
		Logger: l.Logger.With().Interface(key, value).Logger(),
		ctx:    l.ctx,
	}
}
