package readerioresult

import (
	"time"

	"github.com/IBM/fp-go/v2/circuitbreaker"
	"github.com/IBM/fp-go/v2/context/readerio"
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
	metrics circuitbreaker.Metrics,
) CircuitBreaker[T] {
	return circuitbreaker.MakeCircuitBreaker[error, T](
		Left,
		ChainFirstIOK,
		ChainFirstLeftIOK,

		readerio.ChainFirstIOK,

		FromIO,
		Flap,
		Flatten,

		currentTime,
		closedState,
		circuitbreaker.MakeCircuitBreakerError,
		checkError,
		policy,
		metrics,
	)
}

func MakeSingletonBreaker[T any](
	currentTime IO[time.Time],
	closedState ClosedState,
	checkError option.Kleisli[error, error],
	policy retry.RetryPolicy,
	metrics circuitbreaker.Metrics,
) Operator[T, T] {
	return circuitbreaker.MakeSingletonBreaker(
		MakeCircuitBreaker[T](
			currentTime,
			closedState,
			checkError,
			policy,
			metrics,
		),
		closedState,
	)
}
