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
	"github.com/IBM/fp-go/v2/readerio"
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
		"openState.OpenedAt",
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
	//
	// Thread Safety: The returned IORef[BreakerState] is thread-safe. It uses atomic
	// operations for all read/write/modify operations. The BreakerState itself is immutable.
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

	// modifyV creates a Reader that sequences an IORef modification operation.
	// It takes an IORef[BreakerState] and returns a Reader that, when given an endomorphism
	// (a function from BreakerState to BreakerState), produces an IO operation that modifies
	// the IORef and returns the new state.
	//
	// This is used internally to create state modification operations that can be composed
	// with other Reader-based operations in the circuit breaker logic.
	//
	// Thread Safety: The IORef modification is atomic. Multiple concurrent calls will be
	// serialized by the IORef's atomic operations.
	//
	// Type signature: Reader[IORef[BreakerState], IO[Endomorphism[BreakerState]]]
	modifyV = reader.Sequence(ioref.Modify[BreakerState])

	initialRetry = retry.DefaultRetryStatus

	// testCircuit sets the canaryRequest flag to true in an openState.
	// This is used to mark that the circuit breaker is in half-open state,
	// allowing a single test request (canary) to check if the service has recovered.
	//
	// When canaryRequest is true:
	//   - One request is allowed through to test the service
	//   - If the canary succeeds, the circuit closes
	//   - If the canary fails, the circuit remains open with an extended reset time
	//
	// Thread Safety: This is a pure function that returns a new openState; it does not
	// modify its input. Safe for concurrent use.
	//
	// Type signature: Endomorphism[openState]
	testCircuit = canaryRequestLens.Set(true)
)

// makeOpenCircuitFromPolicy creates a function that constructs an openState from a retry policy.
// This is a curried function that takes a retry policy and returns a function that takes a retry status
// and current time to produce an openState with calculated reset time.
//
// The function applies the retry policy to determine the next retry delay and calculates
// the resetAt time by adding the delay to the current time. If no previous delay exists
// (first failure), the resetAt is set to the current time.
//
// Parameters:
//   - policy: The retry policy that determines backoff strategy (e.g., exponential backoff)
//
// Returns:
//   - A curried function that takes:
//     1. rs (retry.RetryStatus): The current retry status containing retry count and previous delay
//     2. ct (time.Time): The current time when the circuit is opening
//     And returns an openState with:
//   - openedAt: Set to the current time (ct)
//   - resetAt: Current time plus the delay from the retry policy
//   - retryStatus: The updated retry status from applying the policy
//   - canaryRequest: false (will be set to true when reset time is reached)
//
// Thread Safety: This is a pure function that creates new openState instances.
// Safe for concurrent use.
//
// Example:
//
//	policy := retry.ExponentialBackoff(1*time.Second, 2.0, 10)
//	makeOpen := makeOpenCircuitFromPolicy(policy)
//	openState := makeOpen(retry.DefaultRetryStatus)(time.Now())
//	// openState.resetAt will be approximately 1 second from now
func makeOpenCircuitFromPolicy(policy retry.RetryPolicy) func(rs retry.RetryStatus) func(ct time.Time) openState {

	return func(rs retry.RetryStatus) func(ct time.Time) openState {

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
}

// extendOpenCircuitFromMakeCircuit creates a function that extends the open state of a circuit breaker
// when a canary request fails. It takes a circuit maker function and returns a function that,
// given the current time, produces an endomorphism that updates an openState.
//
// This function is used when a canary request (test request in half-open state) fails.
// It extends the circuit breaker's open period by:
//  1. Extracting the current retry status from the open state
//  2. Using the makeCircuit function to calculate a new open state with updated retry status
//  3. Applying the current time to get the new state
//  4. Setting the canaryRequest flag to true to allow another test request later
//
// Parameters:
//   - makeCircuit: A function that creates an openState from a retry status and current time.
//     This is typically created by makeOpenCircuitFromPolicy.
//
// Returns:
//   - A curried function that takes:
//     1. ct (time.Time): The current time when extending the circuit
//     And returns an Endomorphism[openState] that:
//   - Increments the retry count
//   - Calculates a new resetAt time based on the retry policy (typically with exponential backoff)
//   - Sets canaryRequest to true for the next test attempt
//
// Thread Safety: This is a pure function that returns new openState instances.
// Safe for concurrent use.
//
// Usage Context:
//   - Called when a canary request fails in the half-open state
//   - Extends the open period with increased backoff delay
//   - Prepares the circuit for another canary attempt at the new resetAt time
func extendOpenCircuitFromMakeCircuit(
	makeCircuit func(rs retry.RetryStatus) func(ct time.Time) openState,
) func(time.Time) Endomorphism[openState] {
	return func(ct time.Time) Endomorphism[openState] {
		return F.Flow4(
			retryStatusLens.Get,
			makeCircuit,
			identity.Flap[openState](ct),
			testCircuit,
		)
	}
}

// isResetTimeExceeded checks if the reset time for an open circuit has been exceeded.
// This is used to determine if the circuit breaker should transition from open to half-open state
// by allowing a canary request.
//
// The function returns an option.Kleisli that succeeds (returns Some) only when:
//  1. The circuit is not already in canary mode (canaryRequest is false)
//  2. The current time is after the resetAt time
//
// Parameters:
//   - ct: The current time to compare against the reset time
//
// Returns:
//   - An option.Kleisli[openState, openState] that:
//   - Returns Some(openState) if the reset time has been exceeded and no canary is active
//   - Returns None if the reset time has not been exceeded or a canary request is already active
//
// Thread Safety: This is a pure function that does not modify its input.
// Safe for concurrent use.
//
// Usage Context:
//   - Called when the circuit is open to check if it's time to attempt a canary request
//   - If this returns Some, the circuit transitions to half-open state (canary mode)
//   - If this returns None, the circuit remains fully open and requests are blocked
func isResetTimeExceeded(ct time.Time) option.Kleisli[openState, openState] {
	return option.FromPredicate(func(open openState) bool {
		return !open.canaryRequest && ct.After(resetAtLens.Get(open))
	})
}

// handleSuccessOnClosed creates a Reader that handles successful requests when the circuit is closed.
// This function is used to update the circuit breaker state after a successful operation completes
// while the circuit is in the closed state.
//
// The function takes a Reader that adds a success record to the ClosedState and lifts it to work
// with BreakerState by mapping over the Right (closed) side of the Either type. This ensures that
// success tracking only affects the closed state and leaves any open state unchanged.
//
// Parameters:
//   - addSuccess: A Reader that takes the current time and returns an Endomorphism that updates
//     the ClosedState by recording a successful operation. This typically increments a success
//     counter or updates a success history.
//
// Returns:
//   - A Reader[time.Time, Endomorphism[BreakerState]] that, when given the current time, produces
//     an endomorphism that updates the BreakerState by applying the success update to the closed
//     state (if closed) or leaving the state unchanged (if open).
//
// Thread Safety: This is a pure function that creates new state instances. The returned
// endomorphism is safe for concurrent use as it does not mutate its input.
//
// Usage Context:
//   - Called after a successful request completes while the circuit is closed
//   - Updates success metrics/counters in the ClosedState
//   - Does not affect the circuit state if it's already open
//   - Part of the normal operation flow when the circuit breaker is functioning properly
func handleSuccessOnClosed(
	addSuccess Reader[time.Time, Endomorphism[ClosedState]],
) Reader[time.Time, Endomorphism[BreakerState]] {
	return F.Flow2(
		addSuccess,
		either.Map[openState],
	)
}

// handleFailureOnClosed creates a Reader that handles failed requests when the circuit is closed.
// This function manages the critical logic for determining whether a failure should cause the
// circuit breaker to open (transition from closed to open state).
//
// The function orchestrates three key operations:
//  1. Records the failure in the ClosedState using addError
//  2. Checks if the failure threshold has been exceeded using checkClosedState
//  3. If threshold exceeded, opens the circuit; otherwise, keeps it closed with updated error count
//
// The decision flow is:
//   - Add the error to the closed state's error tracking
//   - Check if the updated closed state exceeds the failure threshold
//   - If threshold exceeded (checkClosedState returns None):
//   - Create a new openState with calculated reset time based on retry policy
//   - Transition the circuit to open state (Left side of Either)
//   - If threshold not exceeded (checkClosedState returns Some):
//   - Keep the circuit closed with the updated error count
//   - Continue allowing requests through
//
// Parameters:
//   - addError: A Reader that takes the current time and returns an Endomorphism that updates
//     the ClosedState by recording a failed operation. This typically increments an error
//     counter or adds to an error history.
//   - checkClosedState: A Reader that takes the current time and returns an option.Kleisli that
//     validates whether the ClosedState is still within acceptable failure thresholds.
//     Returns Some(ClosedState) if threshold not exceeded, None if threshold exceeded.
//   - openCircuit: A Reader that takes the current time and creates a new openState with
//     appropriate reset time calculated from the retry policy. Used when transitioning to open.
//
// Returns:
//   - A Reader[time.Time, Endomorphism[BreakerState]] that, when given the current time, produces
//     an endomorphism that either:
//   - Keeps the circuit closed with updated error tracking (if threshold not exceeded)
//   - Opens the circuit with calculated reset time (if threshold exceeded)
//
// Thread Safety: This is a pure function that creates new state instances. The returned
// endomorphism is safe for concurrent use as it does not mutate its input.
//
// Usage Context:
//   - Called after a failed request completes while the circuit is closed
//   - Implements the core circuit breaker logic for opening the circuit
//   - Determines when to stop allowing requests through to protect the failing service
//   - Critical for preventing cascading failures in distributed systems
//
// State Transition:
//   - Closed (under threshold) -> Closed (with incremented error count)
//   - Closed (at/over threshold) -> Open (with reset time for recovery attempt)
func handleFailureOnClosed(
	addError Reader[time.Time, Endomorphism[ClosedState]],
	checkClosedState Reader[time.Time, option.Kleisli[ClosedState, ClosedState]],
	openCircuit Reader[time.Time, openState],
) Reader[time.Time, Endomorphism[BreakerState]] {
	return F.Pipe2(
		F.Pipe1(
			addError,
			reader.ApS(reader.Map[ClosedState], checkClosedState),
		),
		reader.Chain(F.Flow2(
			reader.Map[ClosedState](option.Fold(
				F.Pipe2(
					openCircuit,
					reader.Map[time.Time](createOpenCircuit),
					lazy.Of,
				),
				F.Flow2(
					createClosedCircuit,
					reader.Of[time.Time],
				),
			)),
			reader.Sequence,
		)),
		reader.Map[time.Time](either.Chain[openState, ClosedState, ClosedState]),
	)
}

func handleErrorOnClosed2[E any](
	checkError option.Kleisli[E, E],
	onSuccess Reader[time.Time, Endomorphism[BreakerState]],
	onFailure Reader[time.Time, Endomorphism[BreakerState]],
) reader.Kleisli[time.Time, E, Endomorphism[BreakerState]] {
	return F.Flow3(
		checkError,
		option.MapTo[E](onFailure),
		option.GetOrElse(lazy.Of(onSuccess)),
	)
}

func stateModifier(
	modify io.Kleisli[Endomorphism[BreakerState], BreakerState],
) reader.Operator[time.Time, Endomorphism[BreakerState], IO[BreakerState]] {
	return reader.Map[time.Time](modify)
}

func reportOnClose2(
	onClosed ReaderIO[time.Time, Void],
	onOpened ReaderIO[time.Time, Void],
) readerio.Operator[time.Time, BreakerState, Void] {
	return readerio.Chain(either.Fold(
		reader.Of[openState](onOpened),
		reader.Of[ClosedState](onClosed),
	))
}

func applyAndReportClose2(
	currentTime IO[time.Time],
	metrics readerio.Operator[time.Time, BreakerState, Void],
) func(io.Kleisli[Endomorphism[BreakerState], BreakerState]) func(Reader[time.Time, Endomorphism[BreakerState]]) IO[Void] {
	return func(modify io.Kleisli[Endomorphism[BreakerState], BreakerState]) func(Reader[time.Time, Endomorphism[BreakerState]]) IO[Void] {
		return F.Flow3(
			reader.Map[time.Time](modify),
			metrics,
			readerio.ReadIO[Void](currentTime),
		)
	}
}

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
// Thread Safety: The returned State monad creates operations that are thread-safe when
// executed. The IORef[BreakerState] uses atomic operations for all state modifications.
// Multiple concurrent requests will be properly serialized at the IORef level.
//
// Returns:
//   - A State monad that transforms a pair of (IORef[BreakerState], HKTT) into HKTT,
//     applying circuit breaker logic to the computation
func MakeCircuitBreaker[E, T, HKTT, HKTOP, HKTHKTT any](

	left func(E) HKTT,
	chainFirstIOK func(io.Kleisli[T, BreakerState]) func(HKTT) HKTT,
	chainFirstLeftIOK func(io.Kleisli[E, BreakerState]) func(HKTT) HKTT,

	chainFirstIOK2 func(io.Kleisli[Either[E, T], Void]) func(HKTT) HKTT,

	fromIO func(IO[func(HKTT) HKTT]) HKTOP,
	flap func(HKTT) func(HKTOP) HKTHKTT,
	flatten func(HKTHKTT) HKTT,

	currentTime IO[time.Time],
	closedState ClosedState,
	makeError Reader[time.Time, E],
	checkError option.Kleisli[E, E],
	policy retry.RetryPolicy,
	metrics Metrics,
) State[Pair[IORef[BreakerState], HKTT], HKTT] {

	type Operator = func(HKTT) HKTT

	addSuccess := reader.From1(ClosedState.AddSuccess)
	addError := reader.From1(ClosedState.AddError)
	checkClosedState := reader.From1(ClosedState.Check)

	closedCircuit := createClosedCircuit(closedState.Empty())
	makeOpenCircuit := makeOpenCircuitFromPolicy(policy)

	openCircuit := F.Pipe1(
		initialRetry,
		makeOpenCircuit,
	)

	extendOpenCircuit := extendOpenCircuitFromMakeCircuit(makeOpenCircuit)

	failWithError := F.Flow4(
		resetAtLens.Get,
		makeError,
		left,
		reader.Of[HKTT],
	)

	handleSuccess2 := handleSuccessOnClosed(addSuccess)
	handleFailure2 := handleFailureOnClosed(addError, checkClosedState, openCircuit)

	handleError2 := handleErrorOnClosed2(checkError, handleSuccess2, handleFailure2)

	metricsClose2 := reportOnClose2(metrics.Accept, metrics.Open)
	apply2 := applyAndReportClose2(currentTime, metricsClose2)

	onClosed := func(modify io.Kleisli[Endomorphism[BreakerState], BreakerState]) Operator {
		return chainFirstIOK2(F.Flow2(
			either.Fold(
				handleError2,
				reader.Of[T](handleSuccess2),
			),
			apply2(modify),
		))
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
// Thread Safety: The returned function is thread-safe. The internal IORef[BreakerState]
// uses atomic operations to manage state. Multiple concurrent calls to the returned function
// will be properly serialized at the state modification level.
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
