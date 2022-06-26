package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

var root = rootLogger{}

type rootLogger struct {
	done uint32
	m    sync.Mutex
	l    *Logger
}

func (rl *rootLogger) verify() {
	if atomic.LoadUint32(&root.done) == 0 {
		rl.setDefault()
	}
}

func (rl *rootLogger) setDefault() {
	rl.m.Lock()
	defer rl.m.Unlock()
	if rl.done == 0 {
		defer atomic.StoreUint32(&rl.done, 1)
		var err error
		rl.l, err = getLogger(Logging{
			Env:   "dev",
			Level: "debug",
		})
		if err != nil {
			panic(err)
		}
	}
}

func (rl *rootLogger) set(cfg Logging) error {
	rl.m.Lock()
	defer rl.m.Unlock()
	var err error
	rl.l, err = getLogger(cfg)
	if err != nil {
		return err
	}
	atomic.StoreUint32(&rl.done, 1)
	return nil
}

// GetLogger return logger with a scope
func GetLogger(scope ...string) *Logger {
	root.verify()
	if len(scope) < 1 {
		return root.l
	}
	module := strings.Join(scope, ".")
	subLogger := root.l.Logger.With().Str("module", module).Logger()
	return &Logger{module: module, Logger: &subLogger}
}

// Init initializes a rs/zerolog logger from user config
func Init(cfg Logging) (err error) {
	if err != root.set(cfg) {
		return err
	}
	return nil
}

// getLogger initializes a root logger
func getLogger(cfg Logging) (*Logger, error) {
	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	var w io.Writer
	switch cfg.Env {
	case "dev":
		cw := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		cw.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		cw.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("***%s****", i)
		}
		cw.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		cw.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}
		w = io.Writer(cw)
	default:
		w = os.Stderr
	}
	l := zerolog.New(w).Level(lvl).With().Timestamp().Logger()
	return &Logger{module: "root", Logger: &l}, nil
}

