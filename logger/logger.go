package logger

import (
	"context"
	"strings"
)

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var ContextKey = contextKey{}
var ErrNoLoggerInContext = errors.New("no logger in context")

type contextKey struct{}

// Logging is the config info
type Logging struct {
	Env   string
	Level string
}

// Logger is wrapper for rs/zerolog logger with module, it is singleton.
type Logger struct {
	module string
	*zerolog.Logger
}

func (l *Logger) Named(name string) *Logger {
	module := strings.Join([]string{l.module, name}, ".")
	subLogger := root.l.With().Str("module", module).Logger()
	return &Logger{module: module, Logger: &subLogger}
}

// Loggable indicates the implement supports logging
type Loggable interface {
	SetLogger(*Logger)
}

func Fetch(ctx context.Context, name string) *Logger {
	parentLogger := ctx.Value(ContextKey)
	if parentLogger == nil {
		return GetLogger(name)
	}
	if pl, ok := parentLogger.(*Logger); ok {
		return pl.Named(name)
	}
	return GetLogger(name)
}
