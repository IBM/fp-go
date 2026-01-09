package readerioresult

import (
	"time"

	"github.com/IBM/fp-go/v2/circuitbreaker"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/retry"
)

type (
	ClosedState = circuitbreaker.ClosedState

	Env[T any] = Pair[IORef[circuitbreaker.BreakerState], ReaderIOResult[T]]

	CircuitBreaker[T any] = State[Env[T], ReaderIOResult[T]]
)

func MakeCircuitBreaker[T any](
	currentTime IO[time.Time],
	closedState ClosedState,
	checkError option.Kleisli[error, error],
	policy retry.RetryPolicy,
	logger io.Kleisli[string, string],
) CircuitBreaker[T] {
	return circuitbreaker.MakeCircuitBreaker[error, T](
		Left,
		ChainFirstIOK,
		ChainFirstLeftIOK,
		FromIO,
		Flap,
		Flatten,

		currentTime,
		closedState,
		circuitbreaker.MakeCircuitBreakerError,
		checkError,
		policy,
		logger,
	)
}

func MakeSingletonBreaker[T any](
	currentTime IO[time.Time],
	closedState ClosedState,
	checkError option.Kleisli[error, error],
	policy retry.RetryPolicy,
	logger io.Kleisli[string, string],
) Operator[T, T] {
	return circuitbreaker.MakeSingletonBreaker(
		MakeCircuitBreaker[T](
			currentTime,
			closedState,
			checkError,
			policy,
			logger,
		),
		closedState,
	)
}
