// Package log provides a logging interface for the application.
package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger represents a logger instance.
type Logger struct {
	logger *zerolog.Logger
}

// NewLogger returns a new logger instance.
// If isDebug is true, the logger will be configured to log at the debug level.
func NewLogger(isDebug bool) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	l := log.With().Caller().Logger()
	return &Logger{logger: &l}
}

// With returns a new logger instance with additional context.
// The args parameter should be a series of key-value pairs, where each key is a string.
func (l *Logger) With(args ...interface{}) *Logger {
	// ...
}

// Debug logs a debug-level message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}

// Info logs an info-level message.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

// Warn logs a warn-level message.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

// Error logs an error-level message.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

// Fatal logs a fatal-level message and exits the program.
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatal().Msgf(msg, args...)
}

// Panic logs a panic-level message and panics.
func (l *Logger) Panic(msg string, args ...interface{}) {
	l.logger.Panic().Msgf(msg, args...)
}

// Trace logs a trace-level message.
func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logger.Trace().Msgf(msg, args...)
}
