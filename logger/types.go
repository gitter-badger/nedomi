package logger

// This file is generated with go generate. Any changes to it will be lost after
// subsequent generates.
//
// If you want to edit it go to types.go.template

import (
	"github.com/ironsmile/nedomi/config"
)

type newLoggerFunc func(cfg config.LoggerSection) (Logger, error)

var loggerTypes map[string]newLoggerFunc = map[string]newLoggerFunc{}
