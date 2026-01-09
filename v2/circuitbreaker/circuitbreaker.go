package circuitbreaker

import (
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/retry"
)

var (
	canaryRequestLens = lens.MakeLensWithName(
		func(os openState) bool { return os.canaryRequest },
		func(os openState, flag bool) openState {
			os.canaryRequest = flag
			return os
		},
		"openState.CanaryRequest",
	)

	retryStatusLens = lens.MakeLensWithName(
		func(os openState) retry.RetryStatus { return os.retryStatus },
		func(os openState, status retry.RetryStatus) openState {
			os.retryStatus = status
			return os
		},
		"openState.RetryStatus",
	)

	resetAtLens = lens.MakeLensWithName(
		func(os openState) time.Time { return os.resetAt },
		func(os openState, tm time.Time) openState {
			os.resetAt = tm
			return os
		},
		"openState.ResetAt",
	)

	openedAtLens = lens.MakeLensWithName(
		func(os openState) time.Time { return os.openedAt },
		func(os openState, tm time.Time) openState {
			os.openedAt = tm
			return os
		},
		"openState.OpendAt",
	)

	createClosedCircuit = either.Right[openState, ClosedState]
	createOpenCircuit   = either.Left[ClosedState, openState]

	// MakeClosedIORef creates an IORef containing a closed circuit breaker state.
	// It wraps the provided ClosedState in a Right (closed) BreakerState and creates
	// a mutable reference to it.
	//
	// Parameters:
	//   - closedState: The initial closed state configuration
	//
	// Returns:
	//   - An IO operation that creates an IORef[BreakerState] initialized to closed state
	MakeClosedIORef = F.Flow2(
		createClosedCircuit,
		ioref.MakeIORef,
	)

	// IsOpen checks if a BreakerState is in the open state.
	// Returns true if the circuit breaker is open (blocking requests), false otherwise.
	IsOpen = either.IsLeft[openState, ClosedState]

	// IsClosed checks if a BreakerState is in the closed state.
	// Returns true if the circuit breaker is closed (allowing requests), false otherwise.
	IsClosed = either.IsRight[openState, ClosedState]
)

// MakeCircuitBreaker creates a circuit breaker implementation for a higher-kinded type.
//
// This is a generic circuit breaker factory that works with any monad-like type (HKTT).
// It implements the circuit breaker pattern by wrapping operations and managing state transitions
// between closed, open, and half-open states based on failure rates and retry policies.
//
// Type Parameters:
//   - E: The error type
//   - T: The success value type
//   - HKTT: The higher-kinded type representing the computation (e.g., IO[T], ReaderIO[R, T])
//   - HKTOP: The higher-kinded type for operators (e.g., IO[func(HKTT) HKTT])
//   - HKTHKTT: The nested higher-kinded type (e.g., IO[IO[T]])
//
// Parameters:
//   - left: Constructs an error result in HKTT from an error value
//   - chainFirstIOK: Chains an IO operation that runs after success, preserving the original value
//   - chainFirstLeftIOK: Chains an IO operation that runs after error, preserving the original error
//   - fromIO: Lifts an IO operation into HKTOP
//   - flap: Applies a value to a function wrapped in a higher-kinded type
//   - flatten: Flattens nested higher-kinded types (join operation)
//   - currentTime: IO operation that provides the current time
//   - closedState: The initial closed state configuration
//   - makeError: Creates an error from a reset time when the circuit is open
//   - checkError: Predicate to determine if an error should trigger circuit breaker logic
//   - policy: Retry policy for determining reset times when circuit opens
//   - logger: Logging function for circuit breaker events
//
// Returns:
//   - A State monad that transforms a pair of (IORef[BreakerState], HKTT) into HKTT,
//     applying circuit breaker logic to the computation
func MakeCircuitBreaker[E, T, HKTT, HKTOP, HKTHKTT any](

	left func(E) HKTT,
	chainFirstIOK func(io.Kleisli[T, BreakerState]) func(HKTT) HKTT,
	chainFirstLeftIOK func(io.Kleisli[E, BreakerState]) func(HKTT) HKTT,

	fromIO func(IO[func(HKTT) HKTT]) HKTOP,
	flap func(HKTT) func(HKTOP) HKTHKTT,
	flatten func(HKTHKTT) HKTT,

	currentTime IO[time.Time],
	closedState ClosedState,
	makeError Reader[time.Time, E],
	checkError option.Kleisli[E, E],
	policy retry.RetryPolicy,
	logger io.Kleisli[string, string],
) State[Pair[IORef[BreakerState], HKTT], HKTT] {

	type Operator = func(HKTT) HKTT

	addSuccess := reader.From1(ClosedState.AddSuccess)
	addError := reader.From1(ClosedState.AddError)
	checkClosedState := reader.From1(ClosedState.Check)

	initialRetry := retry.DefaultRetryStatus

	modifyV := reader.Sequence(ioref.Modify[BreakerState])

	closedCircuit := createClosedCircuit(closedState.Empty())

	makeOpenCircuit := func(rs retry.RetryStatus) func(ct time.Time) openState {

		retryStatus := retry.ApplyPolicy(policy, rs)

		return func(ct time.Time) openState {

			resetTime := F.Pipe2(
				retryStatus,
				retry.PreviousDelayLens.Get,
				option.Fold(
					F.Pipe1(
						ct,
						lazy.Of,
					),
					ct.Add,
				),
			)

			return openState{openedAt: ct, resetAt: resetTime, retryStatus: retryStatus}
		}
	}

	openCircuit := F.Pipe1(
		initialRetry,
		makeOpenCircuit,
	)

	testCircuit := canaryRequestLens.Set(true)

	extendOpenCircuit := func(ct time.Time) Endomorphism[openState] {
		return F.Flow4(
			retryStatusLens.Get,
			makeOpenCircuit,
			identity.Flap[openState](ct),
			testCircuit,
		)
	}

	isResetTimeExceeded := func(ct time.Time) option.Kleisli[openState, openState] {
		return option.FromPredicate(func(open openState) bool {
			return !open.canaryRequest && ct.After(resetAtLens.Get(open))
		})
	}

	failWithError := func(open openState) Operator {
		return F.Pipe4(
			open,
			resetAtLens.Get,
			makeError,
			left,
			reader.Of[HKTT],
		)
	}

	onClosed := func(modify io.Kleisli[Endomorphism[BreakerState], BreakerState]) Operator {

		handleSuccess := F.Pipe2(
			currentTime,
			io.Map(F.Flow2(
				addSuccess,
				either.Map[openState],
			)),
			io.Chain(modify),
		)

		return F.Flow2(
			// error case
			chainFirstLeftIOK(F.Flow2(
				checkError,
				option.Fold(
					// the error is not applicable, handle as success
					F.Pipe1(
						handleSuccess,
						lazy.Of,
					),
					// the error is relevant, record it
					F.Pipe1(
						F.Pipe1(
							currentTime,
							io.Chain(func(ct time.Time) io.IO[BreakerState] {
								return F.Pipe1(
									either.Chain(F.Flow3(
										addError(ct),
										checkClosedState(ct),
										option.Fold(
											F.Pipe2(
												ct,
												lazy.Of,
												lazy.Map(F.Flow2(
													openCircuit,
													createOpenCircuit,
												)),
											),
											createClosedCircuit,
										)),
									),
									modify,
								)
							}),
						),
						reader.Of[E],
					),
				),
			)),
			// good case
			chainFirstIOK(F.Pipe1(
				handleSuccess,
				reader.Of[T],
			)),
		)
	}

	onCanary := func(modify io.Kleisli[Endomorphism[BreakerState], BreakerState]) Operator {

		handleSuccess := F.Pipe2(
			closedCircuit,
			reader.Of[BreakerState],
			modify,
		)

		return F.Flow2(
			// the canary request fails
			chainFirstLeftIOK(F.Flow2(
				checkError,
				option.Fold(
					// the canary request succeeds, we close the circuit
					F.Pipe1(
						handleSuccess,
						lazy.Of,
					),
					// the canary request fails, we extend the circuit
					F.Pipe1(
						F.Pipe1(
							currentTime,
							io.Chain(func(ct time.Time) IO[BreakerState] {
								return F.Pipe1(
									F.Flow2(
										either.Fold(
											extendOpenCircuit(ct),
											F.Pipe1(
												openCircuit(ct),
												reader.Of[ClosedState],
											),
										),
										createOpenCircuit,
									),
									modify,
								)
							}),
						),
						reader.Of[E],
					),
				),
			)),
			// the canary request succeeds, we'll close the circuit
			chainFirstIOK(F.Pipe1(
				handleSuccess,
				reader.Of[T],
			)),
		)
	}

	onOpen := func(ref IORef[BreakerState]) Operator {

		modify := modifyV(ref)

		return F.Pipe3(
			currentTime,
			io.Chain(func(ct time.Time) IO[Operator] {
				return F.Pipe1(
					ref,
					ioref.ModifyWithResult(either.Fold(
						func(open openState) Pair[BreakerState, Operator] {
							return option.Fold(
								func() Pair[BreakerState, Operator] {
									return pair.MakePair(createOpenCircuit(open), failWithError(open))
								},
								func(open openState) Pair[BreakerState, Operator] {
									return pair.MakePair(createOpenCircuit(testCircuit(open)), onCanary(modify))
								},
							)(isResetTimeExceeded(ct)(open))
						},
						func(closed ClosedState) Pair[BreakerState, Operator] {
							return pair.MakePair(createClosedCircuit(closed), onClosed(modify))
						},
					)),
				)
			}),
			fromIO,
			func(src HKTOP) Operator {
				return func(rdr HKTT) HKTT {
					return F.Pipe2(
						src,
						flap(rdr),
						flatten,
					)
				}
			},
		)
	}

	return func(e Pair[IORef[BreakerState], HKTT]) Pair[Pair[IORef[BreakerState], HKTT], HKTT] {
		return pair.MakePair(e, onOpen(pair.Head(e))(pair.Tail(e)))
	}
}

// MakeSingletonBreaker creates a singleton circuit breaker operator for a higher-kinded type.
//
// This function creates a circuit breaker that maintains its own internal state reference.
// It's called "singleton" because it creates a single, self-contained circuit breaker instance
// with its own IORef for state management. The returned function can be used to wrap
// computations with circuit breaker protection.
//
// Type Parameters:
//   - HKTT: The higher-kinded type representing the computation (e.g., IO[T], ReaderIO[R, T])
//
// Parameters:
//   - cb: The circuit breaker State monad created by MakeCircuitBreaker
//   - closedState: The initial closed state configuration for the circuit breaker
//
// Returns:
//   - A function that wraps a computation (HKTT) with circuit breaker logic.
//     The circuit breaker state is managed internally and persists across invocations.
//
// Example Usage:
//
//	// Create a circuit breaker for IO operations
//	breaker := MakeSingletonBreaker(
//	    MakeCircuitBreaker(...),
//	    MakeClosedStateCounter(3),
//	)
//
//	// Use it to wrap operations
//	protectedOp := breaker(myIOOperation)
func MakeSingletonBreaker[HKTT any](
	cb State[Pair[IORef[BreakerState], HKTT], HKTT],
	closedState ClosedState,
) func(HKTT) HKTT {
	return F.Flow3(
		F.Pipe3(
			closedState,
			MakeClosedIORef,
			io.Run,
			pair.FromHead[HKTT],
		),
		cb,
		pair.Tail,
	)
}
