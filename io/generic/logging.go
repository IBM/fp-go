package generic

import (
	"log"

	Logging "github.com/ibm/fp-go/logging"
)

func Logger[GA ~func() any, A any](loggers ...*log.Logger) func(string) func(A) GA {
	_, right := Logging.LoggingCallbacks(loggers...)
	return func(prefix string) func(A) GA {
		return func(a A) GA {
			return FromImpure[GA](func() {
				right("%s: %v", prefix, a)
			})
		}
	}
}

func Logf[GA ~func() any, A any](loggers ...*log.Logger) func(string) func(A) GA {
	_, right := Logging.LoggingCallbacks(loggers...)
	return func(prefix string) func(A) GA {
		return func(a A) GA {
			return FromImpure[GA](func() {
				right(prefix, a)
			})
		}
	}
}
