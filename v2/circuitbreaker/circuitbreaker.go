package circuitbreaker

import (
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
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

// handleSuccessOnClosed handles a successful request when the circuit breaker is in closed state.
// It updates the closed state by recording the success and returns an IO operation that
// modifies the breaker state.
//
// This function is part of the circuit breaker's state management for the closed state.
// When a request succeeds in closed state:
//  1. The current time is obtained
//  2. The addSuccess function is called with the current time to update the ClosedState
//  3. The updated ClosedState is wrapped in a Right (closed) BreakerState
//  4. The breaker state is modified with the new state
//
// Parameters:
//   - currentTime: An IO operation that provides the current time
//   - addSuccess: A Reader that takes a time and returns an endomorphism for ClosedState,
//     typically resetting failure counters or history
//
// Returns:
//   - An io.Kleisli that takes another io.Kleisli and chains them together.
//     The outer Kleisli takes an Endomorphism[BreakerState] and returns BreakerState.
//     This allows composing the success handling with other state modifications.
//
// Thread Safety: This function creates IO operations that will atomically modify the
// IORef[BreakerState] when executed. The state modifications are thread-safe.
//
// Type signature:
//
//	io.Kleisli[io.Kleisli[Endomorphism[BreakerState], BreakerState], BreakerState]
//
// Usage Context:
//   - Called when a request succeeds while the circuit is closed
//   - Resets failure tracking (counter or history) in the ClosedState
//   - Keeps the circuit in closed state
func handleSuccessOnClosed(
	currentTime IO[time.Time],
	addSuccess Reader[time.Time, Endomorphism[ClosedState]],
) io.Kleisli[io.Kleisli[Endomorphism[BreakerState], BreakerState], BreakerState] {
	return F.Flow2(
		io.Chain,
		identity.Flap[IO[BreakerState]](F.Pipe1(
			currentTime,
			io.Map(F.Flow2(
				addSuccess,
				either.Map[openState],
			)))),
	)
}

// handleFailureOnClosed handles a failed request when the circuit breaker is in closed state.
// It updates the closed state by recording the failure and checks if the circuit should open.
//
// This function is part of the circuit breaker's state management for the closed state.
// When a request fails in closed state:
//  1. The current time is obtained
//  2. The addError function is called to record the failure in the ClosedState
//  3. The checkClosedState function is called to determine if the failure threshold is exceeded
//  4. If the threshold is exceeded (Check returns None):
//     - The circuit transitions to open state using openCircuit
//     - A new openState is created with resetAt time calculated from the retry policy
//  5. If the threshold is not exceeded (Check returns Some):
//     - The circuit remains closed with the updated failure tracking
//
// Parameters:
//   - currentTime: An IO operation that provides the current time
//   - addError: A Reader that takes a time and returns an endomorphism for ClosedState,
//     recording a failure (incrementing counter or adding to history)
//   - checkClosedState: A Reader that takes a time and returns an option.Kleisli that checks
//     if the ClosedState should remain closed. Returns Some if circuit stays closed, None if it should open.
//   - openCircuit: A Reader that takes a time and returns an openState with calculated resetAt time
//
// Returns:
//   - An io.Kleisli that takes another io.Kleisli and chains them together.
//     The outer Kleisli takes an Endomorphism[BreakerState] and returns BreakerState.
//     This allows composing the failure handling with other state modifications.
//
// Thread Safety: This function creates IO operations that will atomically modify the
// IORef[BreakerState] when executed. The state modifications are thread-safe.
//
// Type signature:
//
//	io.Kleisli[io.Kleisli[Endomorphism[BreakerState], BreakerState], BreakerState]
//
// State Transitions:
//   - Closed -> Closed: When failure threshold is not exceeded (Some from checkClosedState)
//   - Closed -> Open: When failure threshold is exceeded (None from checkClosedState)
//
// Usage Context:
//   - Called when a request fails while the circuit is closed
//   - Records the failure in the ClosedState (counter or history)
//   - May trigger transition to open state if threshold is exceeded
func handleFailureOnClosed(
	currentTime IO[time.Time],
	addError Reader[time.Time, Endomorphism[ClosedState]],
	checkClosedState Reader[time.Time, option.Kleisli[ClosedState, ClosedState]],
	openCircuit Reader[time.Time, openState],
) io.Kleisli[io.Kleisli[Endomorphism[BreakerState], BreakerState], BreakerState] {

	return F.Flow2(
		io.Chain,
		identity.Flap[IO[BreakerState]](F.Pipe1(
			currentTime,
			io.Map(func(ct time.Time) either.Operator[openState, ClosedState, ClosedState] {
				return either.Chain(F.Flow3(
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
					),
				))
			}))),
	)

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

	handleSuccess := handleSuccessOnClosed(currentTime, addSuccess)
	handleFailure := handleFailureOnClosed(currentTime, addError, checkClosedState, openCircuit)

	onClosed := func(modify io.Kleisli[Endomorphism[BreakerState], BreakerState]) Operator {

		return F.Flow2(
			// error case
			chainFirstLeftIOK(F.Flow3(
				checkError,
				option.Fold(
					// the error is not applicable, handle as success
					F.Pipe2(
						modify,
						handleSuccess,
						lazy.Of,
					),
					// the error is relevant, record it
					F.Pipe2(
						modify,
						handleFailure,
						reader.Of[E],
					),
				),
				// metering
				io.ChainFirst(either.Fold(
					F.Flow2(
						openedAtLens.Get,
						metrics.Open,
					),
					func(c ClosedState) IO[Void] {
						return io.Of(function.VOID)
					},
				)),
			)),
			// good case
			chainFirstIOK(F.Pipe2(
				modify,
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
