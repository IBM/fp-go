// Package circuitbreaker provides a functional implementation of the circuit breaker pattern.
// A circuit breaker prevents cascading failures by temporarily blocking requests to a failing service,
// allowing it time to recover before retrying.
//
// # State machine
//
// The breaker moves between three states:
//
//	closed    --- failure threshold exceeded --->  open
//	open      --- resetAt reached ------------->   half-open (one canary request)
//	half-open --- canary succeeds ------------->   closed (failure tracking reset)
//	half-open --- canary fails ---------------->   open (extended resetAt, canary rearmed)
//	half-open --- canary deadline passed ----->   half-open (a new canary replaces the old one)
//
// In the closed state every request is let through and its outcome is recorded in a
// [ClosedState]. Once the [ClosedState] reports that the failure threshold is exceeded,
// the circuit opens and every request is rejected with a [CircuitBreakerError] until the
// reset time calculated from the [retry.RetryPolicy] has passed.
//
// The first request after the reset time becomes the canary: it is let through while all
// other requests keep being rejected. A successful canary closes the circuit and resets
// the failure tracking; a failed canary keeps the circuit open with a longer delay and
// rearms the canary for the new reset time.
//
// Starting a canary also sets a canary deadline, which is the same delay that governs the
// open period. A canary normally replaces the state long before that deadline, but should
// its computation never report back, for example because it panicked, the deadline makes
// the next request a new canary instead of leaving the circuit open forever.
//
// # Which errors count
//
// Not every error means the downstream service is unhealthy. The checkError filter passed
// to [MakeCircuitBreaker] decides which ones do: errors it maps to None are handed to the
// caller unchanged but are treated like a success by the breaker. [AnyError] counts every
// error, [InfrastructureError] counts only network, TLS and 5xx failures.
//
// # Metrics
//
// Every request emits exactly one admission metric ([Metrics.Accept], [Metrics.Canary] or
// [Metrics.Reject]) and every state transition is reported exactly once through
// [Metrics.Open] or [Metrics.Close], also when several requests fail concurrently. See
// [MakeCircuitBreaker] for the details.
//
// # Choosing a failure detector
//
//   - [MakeClosedStateCounter] opens the circuit after n consecutive failures
//   - [MakeClosedStateHistory] opens the circuit after n failures inside a sliding time window
//
// # Getting started
//
// Instantiating [MakeCircuitBreaker] requires a handful of operations of the monad that is
// being protected. Prefer the ready made bindings, for example
// github.com/IBM/fp-go/v2/context/readerioresult.MakeSingletonBreaker, over calling
// [MakeCircuitBreaker] directly.
//
// # Thread Safety
//
// All data structures in this package are immutable except for IORef[BreakerState].
// The IORef provides thread-safe mutable state through atomic operations.
//
// Immutable types (safe for concurrent use):
//   - BreakerState (Either[openState, ClosedState])
//   - openState
//   - ClosedState implementations (closedStateWithErrorCount, closedStateWithHistory)
//   - All function types and readers
//
// Mutable types (thread-safe through atomic operations):
//   - IORef[BreakerState] - provides atomic read/write/modify operations
//
// ClosedState implementations must be thread-safe. The recommended approach is to
// return new copies for all operations (Empty, AddError, AddSuccess, Check), which
// provides automatic thread safety through immutability.
package circuitbreaker

import (
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/IBM/fp-go/v2/state"
)

type (
	// Ord is a type alias for ord.Ord, representing a total ordering on type A.
	// Used for comparing values in a consistent way.
	Ord[A any] = ord.Ord[A]

	// Option is a type alias for option.Option, representing an optional value.
	// It can be either Some(value) or None, used for safe handling of nullable values.
	Option[A any] = option.Option[A]

	// Endomorphism is a type alias for endomorphism.Endomorphism, representing a function from A to A.
	// Used for transformations that preserve the type.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// IO is a type alias for io.IO, representing a lazy computation that produces a value of type T.
	// Used for side-effectful operations that are deferred until execution.
	IO[T any] = io.IO[T]

	// Pair is a type alias for pair.Pair, representing a tuple of two values.
	// Used for grouping related values together.
	Pair[L, R any] = pair.Pair[L, R]

	// IORef is a type alias for ioref.IORef, representing a mutable reference to a value of type T.
	// Used for managing mutable state in a functional way with IO operations.
	IORef[T any] = ioref.IORef[T]

	// State is a type alias for state.State, representing a stateful computation.
	// It transforms a state of type T and produces a result of type R.
	State[T, R any] = state.State[T, R]

	// Either is a type alias for either.Either, representing a value that can be one of two types.
	// Left[E] represents an error or alternative path, Right[A] represents the success path.
	Either[E, A any] = either.Either[E, A]

	// Predicate is a type alias for predicate.Predicate, representing a function that tests a value.
	// Returns true if the value satisfies the predicate condition, false otherwise.
	Predicate[A any] = predicate.Predicate[A]

	// Reader is a type alias for reader.Reader, representing a computation that depends on an environment R
	// and produces a value of type A. Used for dependency injection and configuration.
	Reader[R, A any] = reader.Reader[R, A]

	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// openState represents the internal state when the circuit breaker is open.
	// In the open state, requests are blocked to give the failing service time to recover.
	// The circuit breaker will transition to a half-open state (canary request) after resetAt.
	openState struct {
		// openedAt is the time when the circuit breaker opened the circuit
		openedAt time.Time

		// resetAt is the time when the circuit breaker should attempt a canary request
		// to test if the service has recovered. Calculated based on the retry policy.
		resetAt time.Time

		// retryStatus tracks the current retry attempt information, including the number
		// of retries and the delay between attempts. Used by the retry policy to calculate
		// exponential backoff or other retry strategies.
		retryStatus retry.RetryStatus

		// canaryRequest indicates whether the circuit is in half-open state, allowing
		// a single test request (canary) to check if the service has recovered.
		// If true, one request is allowed through to test the service.
		// If the canary succeeds, the circuit closes; if it fails, the circuit remains open
		// with an extended reset time.
		canaryRequest bool
	}

	// BreakerState represents the current state of the circuit breaker.
	// It is an Either type where:
	//   - Left[openState] represents an open circuit (requests are blocked)
	//   - Right[ClosedState] represents a closed circuit (requests are allowed through)
	//
	// State Transitions:
	//   - Closed -> Open: When failure threshold is exceeded in ClosedState
	//   - Open -> Half-Open: When resetAt is reached (canaryRequest = true)
	//   - Half-Open -> Closed: When canary request succeeds
	//   - Half-Open -> Open: When canary request fails (with extended resetAt)
	BreakerState = Either[openState, ClosedState]

	Void = function.Void
)
