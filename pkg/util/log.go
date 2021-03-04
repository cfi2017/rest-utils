package util

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

var (
	logger  *zap.Logger
	logOnce = sync.Once{}
)

func GetL() *zap.Logger {
	logOnce.Do(func() {
		debug := viper.GetBool("log.debug")
		var (
			l   *zap.Logger
			err error
		)
		if debug {
			l, err = zap.NewDevelopment()
		} else {
			l, err = zap.NewProduction()
		}
		if err != nil {
			panic(err)
		}
		logger = l
	})
	return logger
}

func GetS() *zap.SugaredLogger {
	return GetL().Sugar()
}

func InitialiseLoggerFlags() {
	// logging flags
	pflag.String("log.level", "info", "log level")
	pflag.Bool("log.debug", false, "debug logging")
}
