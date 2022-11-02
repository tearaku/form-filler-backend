package logger

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

var log logr.Logger

// TODO: maybe make this customizable?
func init() {
	zl := zerolog.New(os.Stdout)
	zl.With().Caller().Timestamp().Logger()
	log = zerologr.New(&zl)
}

func GetLogger() *logr.Logger {
	return &log
}
