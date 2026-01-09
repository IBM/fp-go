package circuitbreaker

import (
	"slices"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
)

type (
	// ClosedState represents the closed state of a circuit breaker.
	// In the closed state, requests are allowed to pass through, but failures are tracked.
	// If a failure condition is met, the circuit breaker transitions to an open state.
	//
	// # Thread Safety
	//
	// All ClosedState implementations MUST be thread-safe. The recommended approach is to
	// make all methods return new copies rather than modifying the receiver, which provides
	// automatic thread safety through immutability.
	//
	// Implementations should ensure that:
	//   - Empty() returns a new instance with cleared state
	//   - AddError() returns a new instance with the error recorded
	//   - AddSuccess() returns a new instance with success recorded
	//   - Check() does not modify the receiver
	//
	// Both provided implementations (closedStateWithErrorCount and closedStateWithHistory)
	// follow this pattern and are safe for concurrent use.
	ClosedState interface {
		// Empty returns a new ClosedState with all tracked failures cleared.
		// This is used when transitioning back to a closed state from an open state.
		//
		// Thread Safety: Returns a new instance; safe for concurrent use.
		Empty() ClosedState

		// AddError records a failure at the given time.
		// Returns an updated ClosedState reflecting the recorded failure.
		//
		// Thread Safety: Returns a new instance; safe for concurrent use.
		// The original ClosedState is not modified.
		AddError(time.Time) ClosedState

		// AddSuccess records a successful request at the given time.
		// Returns an updated ClosedState reflecting the successful request.
		//
		// Thread Safety: Returns a new instance; safe for concurrent use.
		// The original ClosedState is not modified.
		AddSuccess(time.Time) ClosedState

		// Check verifies if the circuit breaker should remain closed at the given time.
		// Returns Some(ClosedState) if the circuit should stay closed,
		// or None if the circuit should open due to exceeding the failure threshold.
		//
		// Thread Safety: Does not modify the receiver; safe for concurrent use.
		Check(time.Time) Option[ClosedState]
	}

	// closedStateWithErrorCount is a counter-based implementation of ClosedState.
	// It tracks the number of consecutive failures and opens the circuit when
	// the failure count exceeds a configured threshold.
	//
	// Thread Safety: This implementation is immutable. All methods return new instances
	// rather than modifying the receiver, making it safe for concurrent use without locks.
	closedStateWithErrorCount struct {
		// checkFailures is a Kleisli arrow that checks if the failure count exceeds the threshold.
		// Returns Some(count) if threshold is exceeded, None otherwise.
		checkFailures option.Kleisli[uint, uint]
		// failureCount tracks the current number of consecutive failures.
		failureCount uint
	}

	// closedStateWithHistory is a time-window-based implementation of ClosedState.
	// It tracks failures within a sliding time window and opens the circuit when
	// the failure count within the window exceeds a configured threshold.
	//
	// Thread Safety: This implementation is immutable. All methods return new instances
	// with new slices rather than modifying the receiver, making it safe for concurrent
	// use without locks. The history slice is never modified in place; addToSlice always
	// creates a new slice.
	closedStateWithHistory struct {
		ordTime Ord[time.Time]
		// maxFailures is the maximum number of failures allowed within the time window.
		checkFailures option.Kleisli[int, int]
		timeWindow    time.Duration
		history       []time.Time
	}
)

var (
	failureCountLens = lens.MakeLensStrictWithName(
		func(s *closedStateWithErrorCount) uint { return s.failureCount },
		func(s *closedStateWithErrorCount, c uint) *closedStateWithErrorCount {
			s.failureCount = c
			return s
		},
		"closeStateWithErrorCount.failureCount",
	)

	historyLens = lens.MakeLensRefWithName(
		func(s *closedStateWithHistory) []time.Time { return s.history },
		func(s *closedStateWithHistory, c []time.Time) *closedStateWithHistory {
			s.history = c
			return s
		},
		"closedStateWithHistory.history",
	)

	resetHistory      = historyLens.Set(A.Empty[time.Time]())
	resetFailureCount = failureCountLens.Set(0)
	incFailureCount   = lens.Modify[*closedStateWithErrorCount](N.Add(uint(1)))(failureCountLens)
)

// Empty returns a new closedStateWithErrorCount with the failure count reset to zero.
//
// Thread Safety: Returns a new instance; the original is not modified.
// Safe for concurrent use.
func (s *closedStateWithErrorCount) Empty() ClosedState {
	return resetFailureCount(s)
}

// AddError increments the failure count and returns a new closedStateWithErrorCount.
// The time parameter is ignored in this counter-based implementation.
//
// Thread Safety: Returns a new instance; the original is not modified.
// Safe for concurrent use.
func (s *closedStateWithErrorCount) AddError(_ time.Time) ClosedState {
	return incFailureCount(s)
}

// AddSuccess resets the failure count to zero and returns a new closedStateWithErrorCount.
// The time parameter is ignored in this counter-based implementation.
//
// Thread Safety: Returns a new instance; the original is not modified.
// Safe for concurrent use.
func (s *closedStateWithErrorCount) AddSuccess(_ time.Time) ClosedState {
	return resetFailureCount(s)
}

// Check verifies if the failure count is below the threshold.
// Returns Some(ClosedState) if below threshold, None if at or above threshold.
// The time parameter is ignored in this counter-based implementation.
//
// Thread Safety: Does not modify the receiver; safe for concurrent use.
func (s *closedStateWithErrorCount) Check(_ time.Time) Option[ClosedState] {
	return F.Pipe3(
		s,
		failureCountLens.Get,
		s.checkFailures,
		option.MapTo[uint](ClosedState(s)),
	)
}

// MakeClosedStateCounter creates a counter-based ClosedState implementation.
// The circuit breaker will open when the number of consecutive failures reaches maxFailures.
//
// Parameters:
//   - maxFailures: The threshold for consecutive failures. The circuit opens when
//     failureCount >= maxFailures (greater than or equal to).
//
// Returns:
//   - A ClosedState that tracks failures using a simple counter.
//
// Example:
//   - If maxFailures is 3, the circuit will open on the 3rd consecutive failure.
//   - Each AddError call increments the counter.
//   - Each AddSuccess call resets the counter to 0 (only consecutive failures count).
//   - Empty resets the counter to 0.
//
// Behavior:
//   - Check returns Some(ClosedState) when failureCount < maxFailures (circuit stays closed)
//   - Check returns None when failureCount >= maxFailures (circuit should open)
//   - AddSuccess resets the failure count, so only consecutive failures trigger circuit opening
//
// Thread Safety: The returned ClosedState is safe for concurrent use. All methods
// return new instances rather than modifying the receiver.
func MakeClosedStateCounter(maxFailures uint) ClosedState {
	return &closedStateWithErrorCount{
		checkFailures: option.FromPredicate(N.LessThan(maxFailures)),
	}
}

// Empty returns a new closedStateWithHistory with an empty failure history.
//
// Thread Safety: Returns a new instance with a new empty slice; the original is not modified.
// Safe for concurrent use.
func (s *closedStateWithHistory) Empty() ClosedState {
	return resetHistory(s)
}

// addToSlice creates a new sorted slice by adding an item to an existing slice.
// This function does not modify the input slice; it creates a new slice with the item added
// and returns it in sorted order.
//
// Parameters:
//   - o: An Ord instance for comparing time.Time values to determine sort order
//   - ar: The existing slice of time.Time values (assumed to be sorted)
//   - item: The new time.Time value to add to the slice
//
// Returns:
//   - A new slice containing all elements from ar plus the new item, sorted in ascending order
//
// Implementation Details:
//   - Creates a new slice with capacity len(ar)+1
//   - Copies all elements from ar to the new slice
//   - Appends the new item
//   - Sorts the entire slice using the provided Ord comparator
//
// Thread Safety: This function is pure and does not modify its inputs. It always returns
// a new slice, making it safe for concurrent use. This is a key component of the immutable
// design of closedStateWithHistory.
//
// Note: This function is used internally by closedStateWithHistory.AddError to maintain
// a sorted history of failure timestamps for efficient binary search operations.
func addToSlice(o ord.Ord[time.Time], ar []time.Time, item time.Time) []time.Time {
	cpy := make([]time.Time, len(ar)+1)
	cpy[copy(cpy, ar)] = item
	slices.SortFunc(cpy, o.Compare)
	return cpy
}

// AddError records a failure at the given time and returns a new closedStateWithHistory.
// The new instance contains the failure in its history, with old failures outside the
// time window automatically pruned.
//
// Thread Safety: Returns a new instance with a new history slice; the original is not modified.
// Safe for concurrent use. The addToSlice function creates a new slice, ensuring immutability.
func (s *closedStateWithHistory) AddError(currentTime time.Time) ClosedState {

	addFailureToHistory := F.Pipe1(
		historyLens,
		lens.Modify[*closedStateWithHistory](func(old []time.Time) []time.Time {
			// oldest valid entry
			idx, _ := slices.BinarySearchFunc(old, currentTime.Add(-s.timeWindow), s.ordTime.Compare)
			return addToSlice(s.ordTime, old[idx:], currentTime)
		}),
	)

	return addFailureToHistory(s)
}

// AddSuccess purges the entire failure history and returns a new closedStateWithHistory.
// The time parameter is ignored; any success clears all tracked failures.
//
// Thread Safety: Returns a new instance with a new empty slice; the original is not modified.
// Safe for concurrent use.
func (s *closedStateWithHistory) AddSuccess(_ time.Time) ClosedState {
	return resetHistory(s)
}

// Check verifies if the number of failures in the history is below the threshold.
// Returns Some(ClosedState) if below threshold, None if at or above threshold.
// The time parameter is ignored; the check is based on the current history size.
//
// Thread Safety: Does not modify the receiver; safe for concurrent use.
func (s *closedStateWithHistory) Check(_ time.Time) Option[ClosedState] {

	return F.Pipe4(
		s,
		historyLens.Get,
		A.Size,
		s.checkFailures,
		option.MapTo[int](ClosedState(s)),
	)
}

// MakeClosedStateHistory creates a time-window-based ClosedState implementation.
// The circuit breaker will open when the number of failures within a sliding time window reaches maxFailures.
//
// Unlike MakeClosedStateCounter which tracks consecutive failures, this implementation tracks
// all failures within a time window. However, any successful request will purge the entire history,
// effectively resetting the failure tracking.
//
// Parameters:
//   - timeWindow: The duration of the sliding time window. Failures older than this are automatically
//     discarded from the history when new failures are added.
//   - maxFailures: The threshold for failures within the time window. The circuit opens when
//     the number of failures in the window reaches this value (failureCount >= maxFailures).
//
// Returns:
//   - A ClosedState that tracks failures using a time-based sliding window.
//
// Example:
//   - If timeWindow is 1 minute and maxFailures is 5, the circuit will open when 5 failures
//     occur within any 1-minute period.
//   - Failures older than 1 minute are automatically removed from the history when AddError is called.
//   - Any successful request immediately purges all tracked failures from the history.
//
// Behavior:
//   - AddError records the failure timestamp and removes failures outside the time window
//     (older than currentTime - timeWindow).
//   - AddSuccess purges the entire failure history (all tracked failures are removed).
//   - Check returns Some(ClosedState) when failureCount < maxFailures (circuit stays closed).
//   - Check returns None when failureCount >= maxFailures (circuit should open).
//   - Empty purges the entire failure history.
//
// Time Window Management:
//   - The history is automatically pruned on each AddError call to remove failures older than
//     currentTime - timeWindow.
//   - The history is kept sorted by time for efficient binary search and pruning.
//
// Important Note:
//   - A successful request resets everything by purging the entire history. This means that
//     unlike a pure sliding window, a single success will clear all tracked failures, even
//     those within the time window. This behavior is similar to MakeClosedStateCounter but
//     with time-based tracking for failures.
//
// Thread Safety: The returned ClosedState is safe for concurrent use. All methods return
// new instances with new slices rather than modifying the receiver. The history slice is
// never modified in place.
//
// Use Cases:
//   - Systems where a successful request indicates recovery and past failures should be forgotten.
//   - Rate limiting with success-based reset: Allow bursts of failures but reset on success.
//   - Hybrid approach: Time-based failure tracking with success-based recovery.
func MakeClosedStateHistory(
	timeWindow time.Duration,
	maxFailures uint) ClosedState {
	return &closedStateWithHistory{
		checkFailures: option.FromPredicate(N.LessThan(int(maxFailures))),
		ordTime:       ord.OrdTime(),
		history:       A.Empty[time.Time](),
		timeWindow:    timeWindow,
	}
}
