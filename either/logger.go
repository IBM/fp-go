package either

import (
	"log"

	F "github.com/ibm/fp-go/function"
	L "github.com/ibm/fp-go/logging"
)

func _log[E, A any](left func(string, ...any), right func(string, ...any), prefix string) func(Either[E, A]) Either[E, A] {
	return Fold(
		func(e E) Either[E, A] {
			left("%s: %v", prefix, e)
			return Left[E, A](e)
		},
		func(a A) Either[E, A] {
			right("%s: %v", prefix, a)
			return Right[E](a)
		})
}

func Logger[E, A any](loggers ...*log.Logger) func(string) func(Either[E, A]) Either[E, A] {
	left, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) func(Either[E, A]) Either[E, A] {
		delegate := _log[E, A](left, right, prefix)
		return func(ma Either[E, A]) Either[E, A] {
			return F.Pipe1(
				delegate(ma),
				ChainTo[E, A](ma),
			)
		}
	}
}
