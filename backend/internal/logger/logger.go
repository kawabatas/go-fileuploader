package logger

import (
	"sync"
)

type Logger interface {
	Debugf(string, ...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}

var (
	logger Logger
	logMu  sync.RWMutex
)

func init() {
	SetupStdLogger()
}

func SetupStdLogger() {
	logMu.Lock()
	logger = newStdLogger()
	logMu.Unlock()
}

func Debugf(format string, v ...any) {
	logMu.RLock()
	logger.Debugf(format, v...)
	logMu.RUnlock()
}

func Infof(format string, v ...any) {
	logMu.RLock()
	logger.Infof(format, v...)
	logMu.RUnlock()
}

func Errorf(format string, v ...any) {
	logMu.RLock()
	logger.Errorf(format, v...)
	logMu.RUnlock()
}
