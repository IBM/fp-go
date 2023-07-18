package option

import (
	"log"

	F "github.com/IBM/fp-go/function"
	L "github.com/IBM/fp-go/logging"
)

func _log[A any](left func(string, ...any), right func(string, ...any), prefix string) func(Option[A]) Option[A] {
	return Fold(
		func() Option[A] {
			left("%s", prefix)
			return None[A]()
		},
		func(a A) Option[A] {
			right("%s: %v", prefix, a)
			return Some(a)
		})
}

func Logger[A any](loggers ...*log.Logger) func(string) func(Option[A]) Option[A] {
	left, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) func(Option[A]) Option[A] {
		delegate := _log[A](left, right, prefix)
		return func(ma Option[A]) Option[A] {
			return F.Pipe1(
				delegate(ma),
				ChainTo[A](ma),
			)
		}
	}
}
