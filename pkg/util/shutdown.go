package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	quit         chan os.Signal
	shutdownOnce sync.Once
)

func NewShutdownListener() chan os.Signal {
	shutdownOnce.Do(func() {
		quit = make(chan os.Signal)
		// pipe sigint and sigterm to quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	})
	return quit
}

func Shutdown() {
	if quit != nil {
		quit <- syscall.SIGINT
	}
}
