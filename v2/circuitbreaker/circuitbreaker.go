package circuitbreaker

import (
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/retry"
	TU "github.com/IBM/fp-go/v2/tuple"
)

type (
	// decision is what the circuit breaker decided to do with an incoming request.
	// It pairs the metric to emit for that request with the operator that has to be
	// applied to the protected computation.
	//
	// The metric is deliberately kept as an unexecuted [IO] so that it can be emitted
	// outside of the critical section that inspects and updates the breaker state.
	decision[OP any] = Pair[IO[Void], OP]

	// transition is the result of inspecting the breaker state for an incoming request:
	// the state the breaker moves to, together with the [decision] for that request.
	//
	// This is exactly the shape expected by [ioref.ModifyWithResult], which makes the
	// inspect-decide-update cycle a single atomic operation.
	transition[OP any] = Pair[BreakerState, decision[OP]]
)

var (
	canaryRequestLens = common.MakeLensWithName(
		func(os openState) bool { return os.canaryRequest },
		func(os openState, flag bool) openState {
			os.canaryRequest = flag
			return os
		},
		"openState.CanaryRequest",
	)

	retryStatusLens = common.MakeLensWithName(
		func(os openState) retry.RetryStatus { return os.retryStatus },
		func(os openState, status retry.RetryStatus) openState {
			os.retryStatus = status
			return os
		},
		"openState.RetryStatus",
	)

	resetAtLens = common.MakeLensWithName(
		func(os openState) time.Time { return os.resetAt },
		func(os openState, tm time.Time) openState {
			os.resetAt = tm
			return os
		},
		"openState.ResetAt",
	)

	openedAtLens = common.MakeLensWithName(
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

	// modifyPrevV is the flavour of [modifyV] that reports the state before and after the
	// update as a Pair, so that callers can tell an actual state change from an update
	// that left the circuit on the same side.
	//
	// Thread Safety: The read-modify-write cycle is a single atomic operation of the
	// IORef, so the two states of the pair are always adjacent versions.
	//
	// Type signature: Reader[IORef[BreakerState], io.Kleisli[Endomorphism[BreakerState], Pair[BreakerState, BreakerState]]]
	modifyPrevV = reader.Sequence(F.Flow2(
		trackPrevious,
		ioref.ModifyWithResult[BreakerState, Pair[BreakerState, BreakerState]],
	))

	initialRetry = retry.DefaultRetryStatus

	// beginCanary marks an open circuit as half-open by setting the canaryRequest flag.
	// Exactly one request (the canary) is then allowed through to probe whether the
	// downstream service has recovered; every other request is rejected until the canary
	// deadline set by [armCanary] has passed.
	//
	// The flag records that a canary is in flight. It is not what blocks the other
	// requests, see [canaryAllowed] and [armCanary] for why.
	//
	// Thread Safety: This is a pure function that returns a new openState; it does not
	// modify its input. Safe for concurrent use.
	//
	// Type signature: Endomorphism[openState]
	beginCanary = canaryRequestLens.Set(true)

	// endCanary clears the canaryRequest flag, which records that no canary is in flight
	// for the circuit any more.
	//
	// This is the counterpart of [beginCanary] and is applied when a canary request
	// completes with a failure: the circuit stays open, but with a fresh reset time.
	//
	// Thread Safety: This is a pure function that returns a new openState; it does not
	// modify its input. Safe for concurrent use.
	//
	// Type signature: Endomorphism[openState]
	endCanary = canaryRequestLens.Set(false)

	// addDelay is the curried, data-last form of [time.Time.Add]: given a duration it
	// returns a Reader that shifts the time taken from its environment by that duration.
	//
	// Type signature: func(time.Duration) Reader[time.Time, time.Time]
	addDelay = F.Curry2(F.Swap(time.Time.Add))

	// resetAtFromStatus computes the point in time at which an open circuit becomes
	// eligible for a canary request: the current time shifted by the delay suggested by
	// the retry policy, or the current time itself when the policy suggests none.
	//
	// Type signature: func(retry.RetryStatus) Reader[time.Time, time.Time]
	resetAtFromStatus = F.Flow2(
		retry.PreviousDelayLens.Get,
		option.Fold(
			lazy.Of(reader.Ask[time.Time]()),
			addDelay,
		),
	)

	// canaryAllowed tests whether an open circuit may let a canary request through at the
	// given time: its reset time must have passed.
	//
	// A canary in flight is not tracked through the canaryRequest flag here but through
	// the reset time itself: [armCanary] pushes the reset time to the canary deadline, so
	// further requests are rejected while the canary is running and exactly one canary is
	// in flight per deadline window. Should that canary never report back, the deadline
	// passes and the next request becomes a new canary instead of the circuit staying
	// open forever.
	//
	// Type signature: func(time.Time) Predicate[openState]
	canaryAllowed = resetTimeExceeded

	// canaryDeadline is the time at which the breaker gives up on a canary that has not
	// reported back and admits another one. It is the same delay that governs the open
	// period, taken from the retry status of the circuit.
	//
	// Type signature: func(time.Time) Reader[openState, time.Time]
	canaryDeadline = reader.Sequence(F.Flow2(
		retryStatusLens.Get,
		resetAtFromStatus,
	))

	// isResetTimeExceeded is the [option.Kleisli] flavour of [canaryAllowed]. It returns
	// Some(openState) when the circuit may transition to half-open at the given time and
	// None when it has to stay fully open.
	//
	// Thread Safety: Pure; safe for concurrent use.
	//
	// Type signature: func(time.Time) option.Kleisli[openState, openState]
	isResetTimeExceeded = F.Flow2(
		canaryAllowed,
		option.FromPredicate[openState],
	)

	// noReport is the metric of a state change that is not worth reporting.
	noReport = readerio.Of[time.Time](F.VOID)
)

// fanout applies two readers to the same environment and pairs up their results.
// It is the applicative product of two readers and allows computing both halves of a
// [transition] from a single state value without naming that value.
//
// Thread Safety: Pure; the returned reader is safe for concurrent use as long as both
// inputs are.
func fanout[S, A, B any](f Reader[S, A], g Reader[S, B]) Reader[S, Pair[A, B]] {
	return F.Pipe1(
		reader.SequenceT2(f, g),
		reader.Map[S](TU.Tupled2(pair.MakePair[A, B])),
	)
}

// trackPrevious turns a state update into the shape [ioref.ModifyWithResult] expects: the
// new state to store, paired with both the previous and the new state as the result.
//
// Thread Safety: Pure; the returned reader is safe for concurrent use.
func trackPrevious(f Endomorphism[BreakerState]) Reader[BreakerState, Pair[BreakerState, Pair[BreakerState, BreakerState]]] {
	return F.Flow2(
		fanout(F.Identity[BreakerState], f),
		fanout(
			pair.Tail[BreakerState, BreakerState],
			F.Identity[Pair[BreakerState, BreakerState]],
		),
	)
}

// resetTimeExceeded builds a predicate that tests whether the given time lies after the
// reset time of an open circuit.
//
// Thread Safety: Pure; safe for concurrent use.
func resetTimeExceeded(ct time.Time) Predicate[openState] {
	return F.Flow2(resetAtLens.Get, ct.After)
}

// armCanary marks an open circuit as half-open and pushes its reset time to the canary
// deadline, i.e. the point in time at which the breaker gives up on a canary that never
// reported back.
//
// The deadline is what keeps exactly one canary in flight: while it has not passed,
// [canaryAllowed] rejects every further request. A canary that completes replaces the
// state anyway (closed on success, reopened with a fresh reset time on failure), so the
// deadline only ever takes effect for a canary that does not complete, for example
// because the protected computation panicked. In that case the next request after the
// deadline becomes a new canary instead of the circuit staying open forever.
//
// Thread Safety: Pure; returns a new openState. Safe for concurrent use.
func armCanary(ct time.Time) Endomorphism[openState] {
	return F.Flow2(
		fanout(beginCanary, canaryDeadline(ct)),
		pair.Merge(resetAtLens.Set),
	)
}

// openCircuitFromStatus builds the [openState] that belongs to a retry status, reading the
// time at which the circuit opens from the reader environment.
//
// The resulting state records when the circuit opened, when it becomes eligible for a canary
// request (see [resetAtFromStatus]), the retry status that produced that delay and a cleared
// canary flag.
//
// Thread Safety: This is a pure function that creates new openState instances.
// Safe for concurrent use.
func openCircuitFromStatus(status retry.RetryStatus) Reader[time.Time, openState] {
	return F.Pipe4(
		reader.Do[time.Time](openState{}),
		reader.LetTo[time.Time](retryStatusLens.Set, status),
		reader.LetTo[time.Time](canaryRequestLens.Set, false),
		reader.ApS(openedAtLens.Set, reader.Ask[time.Time]()),
		reader.ApS(resetAtLens.Set, resetAtFromStatus(status)),
	)
}

// makeOpenCircuitFromPolicy creates a Kleisli arrow that constructs an openState from a retry
// policy. It takes a retry policy and returns a function that maps a retry status to a Reader
// which, given the current time, produces an openState with a calculated reset time.
//
// The retry policy is applied to the incoming status to determine the next retry delay; the
// resetAt time is the current time plus that delay. If the policy does not yield a delay (for
// example because the retry budget is exhausted), resetAt equals the current time, so the next
// request immediately becomes a canary.
//
// Parameters:
//   - policy: The retry policy that determines the backoff strategy (e.g. exponential backoff)
//
// Returns:
//   - A reader.Kleisli that takes a [retry.RetryStatus] and yields a Reader over the current
//     time producing an openState with:
//   - openedAt: the current time
//   - resetAt: the current time plus the delay from the retry policy
//   - retryStatus: the updated retry status from applying the policy
//   - canaryRequest: false (set by [beginCanary] once the reset time is reached)
//
// Thread Safety: This is a pure function that creates new openState instances.
// Safe for concurrent use.
//
// Example:
//
//	policy := retry.ExponentialBackoff(1 * time.Second)
//	makeOpen := makeOpenCircuitFromPolicy(policy)
//	state := makeOpen(retry.DefaultRetryStatus)(time.Now())
//	// state.resetAt is approximately one second from now
func makeOpenCircuitFromPolicy(policy retry.RetryPolicy) reader.Kleisli[time.Time, retry.RetryStatus, openState] {
	return F.Flow2(
		F.Bind1st(retry.ApplyPolicy, policy),
		openCircuitFromStatus,
	)
}

// extendOpenCircuitFromMakeCircuit creates a Reader that extends the open period of a circuit
// breaker after a canary request has failed.
//
// Given the current time it produces an endomorphism on openState that
//  1. takes the retry status of the current open state,
//  2. feeds it through makeCircuit to obtain a state with an increased retry count and a new
//     reset time (typically further in the future because of the backoff policy), and
//  3. clears the canary flag via [endCanary] so that the circuit becomes eligible for another
//     canary once the new reset time has passed.
//
// Step 3 is essential: while the canary flag is set, [canaryAllowed] rejects every request, so
// a circuit that kept the flag after a failed canary would never probe the service again.
//
// Parameters:
//   - makeCircuit: A Kleisli arrow that creates an openState from a retry status and the
//     current time, typically created by [makeOpenCircuitFromPolicy]
//
// Returns:
//   - A Reader[time.Time, Endomorphism[openState]]
//
// Thread Safety: This is a pure function that returns new openState instances.
// Safe for concurrent use.
//
// Usage Context:
//   - Called when a canary request fails in the half-open state
//   - Extends the open period with an increased backoff delay
//   - Prepares the circuit for another canary attempt at the new resetAt time
func extendOpenCircuitFromMakeCircuit(
	makeCircuit reader.Kleisli[time.Time, retry.RetryStatus, openState],
) Reader[time.Time, Endomorphism[openState]] {
	return reader.Sequence(F.Flow3(
		retryStatusLens.Get,
		makeCircuit,
		reader.Map[time.Time](endCanary),
	))
}

// reopenState builds the state update that a failed canary request causes.
//
// An open circuit is extended: it keeps its open state but receives the longer reset time
// that extend computed from the retry policy. A circuit that was closed in the meantime
// is opened from scratch with the fresh open state, because the canary just proved that
// the service is still unhealthy.
//
// Parameters:
//   - extend: extends the open period of a circuit that is still open
//   - fresh: the open state for a circuit that has been closed concurrently
//
// Thread Safety: Pure; the returned endomorphism is safe for concurrent use.
func reopenState(extend Endomorphism[openState], fresh openState) Endomorphism[BreakerState] {
	return F.Flow2(
		either.Fold(extend, reader.Of[ClosedState](fresh)),
		createOpenCircuit,
	)
}

// handleSuccessOnClosed creates a Reader that handles successful requests when the circuit is
// closed. This function is used to update the circuit breaker state after a successful operation
// completes while the circuit is in the closed state.
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
//   - Does not affect the circuit state if it is already open
//   - Part of the normal operation flow when the circuit breaker is functioning properly
func handleSuccessOnClosed(
	addSuccess Reader[time.Time, Endomorphism[ClosedState]],
) Reader[time.Time, Endomorphism[BreakerState]] {
	return F.Flow2(
		addSuccess,
		either.Map[openState],
	)
}

// recordFailure records a failure in the closed state and decides, based on the outcome of the
// threshold check, whether the circuit stays closed or opens.
//
// All three inputs have already been resolved for the same point in time; this function is the
// pure core of [handleFailureOnClosed].
//
// Parameters:
//   - addError: records the failure in the ClosedState
//   - check: yields Some(ClosedState) while the failure threshold is not exceeded, None otherwise
//   - open: the openState to switch to when the threshold is exceeded
//
// Thread Safety: Pure; the returned endomorphism is safe for concurrent use.
func recordFailure(
	addError Endomorphism[ClosedState],
	check option.Kleisli[ClosedState, ClosedState],
	open openState,
) Endomorphism[BreakerState] {
	return either.Chain(F.Flow3(
		addError,
		check,
		option.Fold(
			F.Pipe2(open, createOpenCircuit, lazy.Of),
			createClosedCircuit,
		),
	))
}

// handleFailureOnClosed creates a Reader that handles failed requests when the circuit is closed.
// This function manages the critical logic for determining whether a failure should cause the
// circuit breaker to open (transition from closed to open state).
//
// The three readers are combined applicatively (they all depend on the same current time) and
// their results are fed into [recordFailure], which
//  1. records the failure in the ClosedState,
//  2. checks whether the failure threshold has been exceeded, and
//  3. opens the circuit when it has, or keeps it closed with an updated error count otherwise.
//
// Because the whole thing is wrapped in either.Chain, an already open circuit is left untouched.
//
// Parameters:
//   - addError: A Reader that takes the current time and returns an Endomorphism that updates
//     the ClosedState by recording a failed operation.
//   - checkClosedState: A Reader that takes the current time and returns an option.Kleisli that
//     validates whether the ClosedState is still within acceptable failure thresholds.
//     Returns Some(ClosedState) if the threshold is not exceeded, None if it is.
//   - openCircuit: A Reader that takes the current time and creates a new openState with an
//     appropriate reset time calculated from the retry policy.
//
// Returns:
//   - A Reader[time.Time, Endomorphism[BreakerState]] that either keeps the circuit closed with
//     updated error tracking or opens it with a calculated reset time.
//
// Thread Safety: This is a pure function that creates new state instances. The returned
// endomorphism is safe for concurrent use as it does not mutate its input.
//
// State Transition:
//   - Closed (under threshold) -> Closed (with incremented error count)
//   - Closed (at/over threshold) -> Open (with reset time for recovery attempt)
func handleFailureOnClosed(
	addError Reader[time.Time, Endomorphism[ClosedState]],
	checkClosedState Reader[time.Time, option.Kleisli[ClosedState, ClosedState]],
	openCircuit Reader[time.Time, openState],
) Reader[time.Time, Endomorphism[BreakerState]] {
	return F.Pipe1(
		reader.SequenceT3(addError, checkClosedState, openCircuit),
		reader.Map[time.Time](TU.Tupled3(recordFailure)),
	)
}

// handleErrorOnClosed selects the state update for a request that completed with an error while
// the circuit was closed.
//
// checkError decides whether an error is relevant for the circuit breaker at all: errors it maps
// to None (client errors, validation failures, ...) are treated like a success and do not count
// towards the failure threshold, errors it maps to Some are recorded as failures.
//
// Thread Safety: Pure; safe for concurrent use.
func handleErrorOnClosed[E any](
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

// reportTransition reports the state change a completed request caused, given the state
// before and after its update as a Pair.
//
// Only an actual change of side is reported. A request that leaves the circuit open
// although it was already open reports nothing, which matters under concurrency: when
// several requests are admitted while the circuit is closed and all of them fail, only
// the one that opened the circuit reports it, while the others find it open already.
//
// Parameters:
//   - onOpened: reported when the request closed -> open
//   - onClosed: reported when the request open -> closed
//
// Thread Safety: Pure; the returned operator is safe for concurrent use.
func reportTransition(
	onOpened ReaderIO[time.Time, Void],
	onClosed ReaderIO[time.Time, Void],
) readerio.Operator[time.Time, Pair[BreakerState, BreakerState], Void] {
	// the circuit is open now, report it if it was closed before
	openedEdge := F.Pipe1(
		IsClosed,
		predicate.Fold(
			F.Constant1[BreakerState](noReport),
			F.Constant1[BreakerState](onOpened),
		),
	)

	// the circuit is closed now, report it if it was open before
	closedEdge := F.Pipe1(
		IsOpen,
		predicate.Fold(
			F.Constant1[BreakerState](noReport),
			F.Constant1[BreakerState](onClosed),
		),
	)

	return readerio.Chain(pair.Merge(either.Fold(
		reader.Of[openState](openedEdge),
		reader.Of[ClosedState](closedEdge),
	)))
}

// applyAndReport turns a state update into the effect that performs the update and reports the
// resulting circuit state to the metrics sink.
//
// The returned reader takes a time dependent state update, applies the current time to it,
// commits it to the IORef through modify and hands the resulting BreakerState to report.
//
// Thread Safety: The state update is atomic; the metric is emitted afterwards.
func applyAndReport(
	currentTime IO[time.Time],
	report readerio.Operator[time.Time, Pair[BreakerState, BreakerState], Void],
	modify io.Kleisli[Endomorphism[BreakerState], Pair[BreakerState, BreakerState]],
) Reader[Reader[time.Time, Endomorphism[BreakerState]], IO[Void]] {
	return F.Flow3(
		reader.Map[time.Time](modify),
		report,
		readerio.ReadIO[Void](currentTime),
	)
}

// MakeCircuitBreaker creates a circuit breaker implementation for a higher-kinded type.
//
// This is a generic circuit breaker factory that works with any monad-like type (HKTT).
// It implements the circuit breaker pattern by wrapping operations and managing state transitions
// between closed, open, and half-open states based on failure rates and retry policies.
//
// # State machine
//
//	closed    --- failure threshold exceeded --->  open
//	open      --- resetAt reached ------------->   half-open (one canary request)
//	half-open --- canary succeeds ------------->   closed (failure tracking reset)
//	half-open --- canary fails ---------------->   open (extended resetAt, canary rearmed)
//	half-open --- canary deadline passed ----->   half-open (a new canary replaces the old one)
//
// # Metrics
//
// Every request emits exactly one admission metric when the breaker inspects its state:
//
//   - [Metrics.Accept] - the circuit was closed and the request was let through
//   - [Metrics.Canary] - the circuit was half-open and the request was used as the probe
//   - [Metrics.Reject] - the circuit was open and the request was blocked
//
// In addition, state transitions are reported once the protected computation completed:
//
//   - [Metrics.Open] - the circuit opened, or a failed canary extended the open period
//   - [Metrics.Close] - a canary request succeeded and the circuit closed again
//
// Only an actual change of state is reported. When several requests are admitted while
// the circuit is closed and all of them fail, the one that crosses the threshold reports
// [Metrics.Open]; the others find the circuit open already and report nothing.
//
// Type Parameters:
//   - E: The error type
//   - T: The success value type
//   - HKTT: The higher-kinded type representing the computation (e.g. IO[T], ReaderIO[R, T])
//   - HKTOP: The higher-kinded type for operators (e.g. IO[func(HKTT) HKTT])
//   - HKTHKTT: The nested higher-kinded type (e.g. IO[IO[T]])
//
// Parameters:
//   - left: Constructs an error result in HKTT from an error value
//   - chainFirstIOK: Chains an IO operation that runs after success, preserving the original value
//   - chainFirstLeftIOK: Chains an IO operation that runs after error, preserving the original error
//   - chainFirstEitherIOK: Chains an IO operation that runs after success and after error,
//     preserving the original outcome
//   - fromIO: Lifts an IO operation into HKTOP
//   - flap: Applies a value to a function wrapped in a higher-kinded type
//   - flatten: Flattens nested higher-kinded types (join operation)
//   - currentTime: IO operation that provides the current time
//   - closedState: The initial closed state configuration
//   - makeError: Creates an error from a reset time when the circuit is open
//   - checkError: Decides whether an error counts towards the failure threshold. Errors mapped to
//     None are treated like a success.
//   - policy: Retry policy for determining reset times when the circuit opens
//   - metrics: Metrics sink for circuit breaker events
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

	chainFirstEitherIOK func(io.Kleisli[Either[E, T], Void]) func(HKTT) HKTT,

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

	// the state of a circuit that opens for the first time
	openCircuit := F.Pipe1(
		initialRetry,
		makeOpenCircuit,
	)

	extendOpenCircuit := extendOpenCircuitFromMakeCircuit(makeOpenCircuit)

	// reopenCircuit extends an already open circuit after a failed canary request and opens
	// a closed circuit from scratch, should it have been closed concurrently
	reopenCircuit := F.Pipe1(
		reader.SequenceT2(extendOpenCircuit, openCircuit),
		reader.Map[time.Time](TU.Tupled2(reopenState)),
	)

	// the computation that rejects a request because the circuit is open
	failWithError := F.Flow4(
		resetAtLens.Get,
		makeError,
		left,
		reader.Of[HKTT],
	)

	handleSuccess := handleSuccessOnClosed(addSuccess)
	handleFailure := handleFailureOnClosed(addError, checkClosedState, openCircuit)
	handleError := handleErrorOnClosed(checkError, handleSuccess, handleFailure)

	// while the circuit is closed the only interesting transition is that it opens
	reportClosedPath := reportTransition(metrics.Open, noReport)

	// metrics for the transitions that are decided on the canary path
	reportOpened := io.Chain(metrics.Open)(currentTime)
	reportClosed := io.Chain(metrics.Close)(currentTime)

	onClosed := func(modify io.Kleisli[Endomorphism[BreakerState], Pair[BreakerState, BreakerState]]) Operator {
		return chainFirstEitherIOK(F.Flow2(
			either.Fold(
				handleError,
				reader.Of[T](handleSuccess),
			),
			applyAndReport(currentTime, reportClosedPath, modify),
		))
	}

	onCanary := func(modify io.Kleisli[Endomorphism[BreakerState], BreakerState]) Operator {

		// the canary succeeded, close the circuit and reset the failure tracking
		closeCircuit := F.Pipe3(
			closedCircuit,
			reader.Of[BreakerState],
			modify,
			io.ChainFirst(F.Constant1[BreakerState](reportClosed)),
		)

		// the canary failed, keep the circuit open but with an extended reset time
		extendCircuit := F.Pipe2(
			currentTime,
			io.Chain(F.Flow2(reopenCircuit, modify)),
			io.ChainFirst(F.Constant1[BreakerState](reportOpened)),
		)

		return F.Flow2(
			// the canary request completed with an error
			chainFirstLeftIOK(F.Flow2(
				checkError,
				option.Fold(
					// the error is none of the breaker's business, treat it like a success
					lazy.Of(closeCircuit),
					// the error is relevant, extend the open period
					reader.Of[E](extendCircuit),
				),
			)),
			// the canary request succeeded, close the circuit
			chainFirstIOK(reader.Of[T](closeCircuit)),
		)
	}

	// applyOperator lifts the operator that was selected for a request into HKTT
	applyOperator := F.Flow2(
		F.Flip(flap),
		reader.Map[HKTT](flatten),
	)

	onOpen := func(ref IORef[BreakerState]) Operator {

		modify := modifyV(ref)

		closedOperator := onClosed(modifyPrevV(ref))
		canaryOperator := onCanary(modify)

		// the circuit is closed, let the request through and track its outcome
		admitRequest := func(ct time.Time) Reader[ClosedState, transition[Operator]] {
			return F.Flow2(
				createClosedCircuit,
				F.Bind2nd(
					pair.MakePair[BreakerState, decision[Operator]],
					pair.MakePair(metrics.Accept(ct), closedOperator),
				),
			)
		}

		// the reset time has passed, use this request as the canary
		startCanary := func(ct time.Time) Reader[openState, transition[Operator]] {
			return F.Flow3(
				armCanary(ct),
				createOpenCircuit,
				F.Bind2nd(
					pair.MakePair[BreakerState, decision[Operator]],
					pair.MakePair(metrics.Canary(ct), canaryOperator),
				),
			)
		}

		// the circuit is open, reject the request without calling the service
		blockRequest := func(ct time.Time) Reader[openState, transition[Operator]] {
			return fanout(
				createOpenCircuit,
				F.Flow2(
					failWithError,
					F.Bind1st(pair.MakePair[IO[Void], Operator], metrics.Reject(ct)),
				),
			)
		}

		// decide inspects the breaker state and yields the state to move to together with
		// the decision for the current request
		decide := func(ct time.Time) Reader[BreakerState, transition[Operator]] {
			return either.Fold(
				F.Pipe1(
					canaryAllowed(ct),
					predicate.Fold(blockRequest(ct), startCanary(ct)),
				),
				admitRequest(ct),
			)
		}

		return F.Pipe5(
			currentTime,
			// inspect and update the breaker state atomically
			io.Chain(F.Flow3(
				decide,
				ioref.ModifyWithResult[BreakerState, decision[Operator]],
				identity.Flap[IO[decision[Operator]]](ref),
			)),
			// emit the admission metric outside of the critical section
			io.ChainFirst(pair.Head[IO[Void], Operator]),
			io.Map(pair.Tail[IO[Void], Operator]),
			fromIO,
			applyOperator,
		)
	}

	return fanout(
		F.Identity[Pair[IORef[BreakerState], HKTT]],
		pair.Paired(F.Uncurry2(onOpen)),
	)
}

// MakeSingletonBreaker creates a singleton circuit breaker operator for a higher-kinded type.
//
// This function creates a circuit breaker that maintains its own internal state reference.
// It is called "singleton" because it creates a single, self-contained circuit breaker instance
// with its own IORef for state management. The returned function can be used to wrap
// computations with circuit breaker protection.
//
// Type Parameters:
//   - HKTT: The higher-kinded type representing the computation (e.g. IO[T], ReaderIO[R, T])
//
// Parameters:
//   - cb: The circuit breaker State monad created by [MakeCircuitBreaker]
//   - closedState: The initial closed state configuration for the circuit breaker. Its failure
//     tracking is reset via [ClosedState.Empty] before the breaker starts using it.
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
		F.Pipe4(
			closedState,
			ClosedState.Empty,
			MakeClosedIORef,
			io.Run,
			pair.FromHead[HKTT],
		),
		cb,
		pair.Tail,
	)
}
