// Package log provides a logging interface for the application.
package log

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger represents a logger instance.
type Logger struct {
	isDebug bool
	logger *zerolog.Logger
	ctx    context.Context
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
	return &Logger{
		isDebug: isDebug,
		logger: &l,
		ctx: context.Background(),
	}
}

// GetIsDebug returns true if the logger is configured to log at the debug level
func (l *Logger) GetIsDebug() bool {
	if l == nil {
		return false
	}
	return l.isDebug
}

// GetCtx returns the context for the logger
func (l *Logger) GetCtx() context.Context {
	return l.ctx
}

// WithCtx set the context for the logger
func (l *Logger) WithCtx(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	l.ctx = ctx
	return l
}

// With returns a new logger instance with additional context.
// The args parameter should be a series of key-value pairs, where each key is a string.
func (l *Logger) With(args ...interface{}) *Logger {
	ctx := l.logger.With()
	totalArgs := len(args)

	if totalArgs % 2 != 0 {
		l.Panic("args must be even")
		return nil
	}

	for i := 0; i < totalArgs; i += 2 {
		k, ok := args[i].(string);
		if !ok {
			l.Panic("key must be string")
			return nil
		}
		ctx = ctx.Interface(k, args[i+1])
	}

	nlog := ctx.Logger()

	return &Logger{
		isDebug: l.GetIsDebug(),
		logger: &nlog,
		ctx: l.GetCtx(),
	}
}

// Debug logs a debug-level message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug().Ctx(l.ctx).Msgf(msg, args...)
}

// Info logs an info-level message.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info().Ctx(l.ctx).Msgf(msg, args...)
}

// Warn logs a warn-level message.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Ctx(l.ctx).Msgf(msg, args...)
}

// Error logs an error-level message.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error().Ctx(l.ctx).Msgf(msg, args...)
}

// Fatal logs a fatal-level message and exits the program.
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatal().Ctx(l.ctx).Msgf(msg, args...)
}

// Panic logs a panic-level message and panics.
func (l *Logger) Panic(msg string, args ...interface{}) {
	l.logger.Panic().Ctx(l.ctx).Msgf(msg, args...)
}

// Trace logs a trace-level message.
func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logger.Trace().Ctx(l.ctx).Msgf(msg, args...)
}
