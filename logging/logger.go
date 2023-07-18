package Logging

import (
	"log"
)

func LoggingCallbacks(loggers ...*log.Logger) (func(string, ...any), func(string, ...any)) {
	switch len(loggers) {
	case 0:
		def := log.Default()
		return def.Printf, def.Printf
	case 1:
		log0 := loggers[0]
		return log0.Printf, log0.Printf
	default:
		return loggers[0].Printf, loggers[1].Printf
	}
}
