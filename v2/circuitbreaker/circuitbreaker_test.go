package circuitbreaker

import (
	"sync"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

type testMetrics struct {
	accepts int
	rejects int
	opens   int
	closes  int
	canary  int

	mu sync.Mutex
}

func (m *testMetrics) Accept(_ time.Time) IO[Void] {
	return func() Void {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.accepts++
		return function.VOID
	}
}

func (m *testMetrics) Open(_ time.Time) IO[Void] {
	return func() Void {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.opens++
		return function.VOID
	}
}

func (m *testMetrics) Close(_ time.Time) IO[Void] {
	return func() Void {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.closes++
		return function.VOID
	}
}

func (m *testMetrics) Reject(_ time.Time) IO[Void] {
	return func() Void {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.rejects++
		return function.VOID
	}
}

func (m *testMetrics) Canary(_ time.Time) IO[Void] {
	return func() Void {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.canary++
		return function.VOID
	}
}

// VirtualTimer provides a controllable time source for testing
type VirtualTimer struct {
	mu      sync.Mutex
	current time.Time
}

func NewMockMetrics() Metrics {
	return &testMetrics{}
}

// NewVirtualTimer creates a new virtual timer starting at the given time
func NewVirtualTimer(start time.Time) *VirtualTimer {
	return &VirtualTimer{current: start}
}

// Now returns the current virtual time
func (vt *VirtualTimer) Now() time.Time {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	return vt.current
}

// Advance moves the virtual time forward by the given duration
func (vt *VirtualTimer) Advance(d time.Duration) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.current = vt.current.Add(d)
}

// Set sets the virtual time to a specific value
func (vt *VirtualTimer) Set(t time.Time) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.current = t
}

// TestModifyV tests the modifyV variable
func TestModifyV(t *testing.T) {
	t.Run("modifyV creates a Reader that modifies IORef", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))

		// Create initial state
		initialState := createClosedCircuit(MakeClosedStateCounter(3))
		ref := io.Run(ioref.MakeIORef(initialState))

		// Create an endomorphism that opens the circuit
		now := vt.Now()
		openState := openState{
			openedAt:      now,
			resetAt:       now.Add(1 * time.Minute),
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}
		endomorphism := func(bs BreakerState) BreakerState {
			return createOpenCircuit(openState)
		}

		// Apply modifyV
		modifyOp := modifyV(ref)
		result := io.Run(modifyOp(endomorphism))

		// Verify the state was modified
		assert.True(t, IsOpen(result), "state should be open after modification")
	})

	t.Run("modifyV returns the new state", func(t *testing.T) {
		initialState := createClosedCircuit(MakeClosedStateCounter(3))
		ref := io.Run(ioref.MakeIORef(initialState))

		// Create a simple endomorphism
		endomorphism := F.Identity[BreakerState]

		modifyOp := modifyV(ref)
		result := io.Run(modifyOp(endomorphism))

		assert.True(t, IsClosed(result), "state should remain closed")
	})
}

// TestTestCircuit tests the testCircuit variable
func TestTestCircuit(t *testing.T) {
	t.Run("testCircuit sets canaryRequest to true", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		now := vt.Now()

		openState := openState{
			openedAt:      now,
			resetAt:       now.Add(1 * time.Minute),
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		result := testCircuit(openState)

		assert.True(t, result.canaryRequest, "canaryRequest should be set to true")
		assert.Equal(t, openState.openedAt, result.openedAt, "openedAt should be unchanged")
		assert.Equal(t, openState.resetAt, result.resetAt, "resetAt should be unchanged")
	})

	t.Run("testCircuit is idempotent", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		now := vt.Now()

		openState := openState{
			openedAt:      now,
			resetAt:       now.Add(1 * time.Minute),
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: true, // already true
		}

		result := testCircuit(openState)

		assert.True(t, result.canaryRequest, "canaryRequest should remain true")
	})

	t.Run("testCircuit preserves other fields", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		now := vt.Now()
		resetTime := now.Add(2 * time.Minute)
		retryStatus := retry.RetryStatus{
			IterNumber:    5,
			PreviousDelay: option.Some(30 * time.Second),
		}

		openState := openState{
			openedAt:      now,
			resetAt:       resetTime,
			retryStatus:   retryStatus,
			canaryRequest: false,
		}

		result := testCircuit(openState)

		assert.Equal(t, now, result.openedAt, "openedAt should be preserved")
		assert.Equal(t, resetTime, result.resetAt, "resetAt should be preserved")
		assert.Equal(t, retryStatus.IterNumber, result.retryStatus.IterNumber, "retryStatus should be preserved")
		assert.True(t, result.canaryRequest, "canaryRequest should be set to true")
	})
}

// TestMakeOpenCircuitFromPolicy tests the makeOpenCircuitFromPolicy function
func TestMakeOpenCircuitFromPolicy(t *testing.T) {
	t.Run("creates openState with calculated reset time", func(t *testing.T) {
		policy := retry.LimitRetries(5)
		makeOpen := makeOpenCircuitFromPolicy(policy)

		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		result := makeOpen(retry.DefaultRetryStatus)(currentTime)

		assert.Equal(t, currentTime, result.openedAt, "openedAt should be current time")
		assert.False(t, result.canaryRequest, "canaryRequest should be false initially")
		assert.NotNil(t, result.retryStatus, "retryStatus should be set")
	})

	t.Run("applies retry policy to calculate delay", func(t *testing.T) {
		// Use exponential backoff policy with limit and cap
		policy := retry.Monoid.Concat(
			retry.LimitRetries(10),
			retry.CapDelay(10*time.Second, retry.ExponentialBackoff(1*time.Second)),
		)
		makeOpen := makeOpenCircuitFromPolicy(policy)

		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)

		// First retry (iter 0)
		result1 := makeOpen(retry.DefaultRetryStatus)(currentTime)

		// The first delay should be approximately 1 second
		expectedResetTime1 := currentTime.Add(1 * time.Second)
		assert.WithinDuration(t, expectedResetTime1, result1.resetAt, 100*time.Millisecond,
			"first reset time should be ~1 second from now")

		// Second retry (iter 1) - should double
		result2 := makeOpen(result1.retryStatus)(currentTime)
		expectedResetTime2 := currentTime.Add(2 * time.Second)
		assert.WithinDuration(t, expectedResetTime2, result2.resetAt, 100*time.Millisecond,
			"second reset time should be ~2 seconds from now")
	})

	t.Run("handles first failure with no previous delay", func(t *testing.T) {
		policy := retry.LimitRetries(3)
		makeOpen := makeOpenCircuitFromPolicy(policy)

		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		result := makeOpen(retry.DefaultRetryStatus)(currentTime)

		// With no previous delay, resetAt should be current time
		assert.Equal(t, currentTime, result.resetAt, "resetAt should be current time when no previous delay")
	})

	t.Run("increments retry iteration number", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		policy := retry.LimitRetries(10)
		makeOpen := makeOpenCircuitFromPolicy(policy)

		currentTime := vt.Now()
		initialStatus := retry.DefaultRetryStatus

		result := makeOpen(initialStatus)(currentTime)

		assert.Greater(t, result.retryStatus.IterNumber, initialStatus.IterNumber,
			"retry iteration should be incremented")
	})

	t.Run("curried function can be partially applied", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		policy := retry.LimitRetries(5)
		makeOpen := makeOpenCircuitFromPolicy(policy)

		// Partially apply with retry status
		makeOpenWithStatus := makeOpen(retry.DefaultRetryStatus)

		currentTime := vt.Now()
		result := makeOpenWithStatus(currentTime)

		assert.NotNil(t, result, "partially applied function should work")
		assert.Equal(t, currentTime, result.openedAt)
	})
}

// TestExtendOpenCircuitFromMakeCircuit tests the extendOpenCircuitFromMakeCircuit function
func TestExtendOpenCircuitFromMakeCircuit(t *testing.T) {
	t.Run("extends open circuit with new retry status", func(t *testing.T) {
		policy := retry.Monoid.Concat(
			retry.LimitRetries(10),
			retry.ExponentialBackoff(1*time.Second),
		)
		makeCircuit := makeOpenCircuitFromPolicy(policy)
		extendCircuit := extendOpenCircuitFromMakeCircuit(makeCircuit)

		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)

		// Create initial open state
		initialOpen := openState{
			openedAt:      currentTime.Add(-1 * time.Minute),
			resetAt:       currentTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		// Extend the circuit
		extendOp := extendCircuit(currentTime)
		result := extendOp(initialOpen)

		assert.True(t, result.canaryRequest, "canaryRequest should be set to true")
		assert.Greater(t, result.retryStatus.IterNumber, initialOpen.retryStatus.IterNumber,
			"retry iteration should be incremented")
		assert.True(t, result.resetAt.After(currentTime), "resetAt should be in the future")
	})

	t.Run("sets canaryRequest to true for next test", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		policy := retry.LimitRetries(5)
		makeCircuit := makeOpenCircuitFromPolicy(policy)
		extendCircuit := extendOpenCircuitFromMakeCircuit(makeCircuit)

		currentTime := vt.Now()
		initialOpen := openState{
			openedAt:      currentTime.Add(-30 * time.Second),
			resetAt:       currentTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		result := extendCircuit(currentTime)(initialOpen)

		assert.True(t, result.canaryRequest, "canaryRequest must be true after extension")
	})

	t.Run("applies exponential backoff on successive extensions", func(t *testing.T) {
		policy := retry.Monoid.Concat(
			retry.LimitRetries(10),
			retry.ExponentialBackoff(1*time.Second),
		)
		makeCircuit := makeOpenCircuitFromPolicy(policy)
		extendCircuit := extendOpenCircuitFromMakeCircuit(makeCircuit)

		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)

		// First extension
		state1 := openState{
			openedAt:      currentTime,
			resetAt:       currentTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}
		result1 := extendCircuit(currentTime)(state1)
		delay1 := result1.resetAt.Sub(currentTime)

		// Second extension (should have longer delay)
		result2 := extendCircuit(currentTime)(result1)
		delay2 := result2.resetAt.Sub(currentTime)

		assert.Greater(t, delay2, delay1, "second extension should have longer delay due to exponential backoff")
	})
}

// TestIsResetTimeExceeded tests the isResetTimeExceeded function
func TestIsResetTimeExceeded(t *testing.T) {
	t.Run("returns Some when reset time is exceeded and no canary active", func(t *testing.T) {
		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		resetTime := currentTime.Add(-1 * time.Second) // in the past

		openState := openState{
			openedAt:      currentTime.Add(-1 * time.Minute),
			resetAt:       resetTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		result := isResetTimeExceeded(currentTime)(openState)

		assert.True(t, option.IsSome(result), "should return Some when reset time exceeded")
	})

	t.Run("returns None when reset time not yet exceeded", func(t *testing.T) {
		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		resetTime := currentTime.Add(1 * time.Minute) // in the future

		openState := openState{
			openedAt:      currentTime.Add(-30 * time.Second),
			resetAt:       resetTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		result := isResetTimeExceeded(currentTime)(openState)

		assert.True(t, option.IsNone(result), "should return None when reset time not exceeded")
	})

	t.Run("returns None when canary request is already active", func(t *testing.T) {
		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		resetTime := currentTime.Add(-1 * time.Second) // in the past

		openState := openState{
			openedAt:      currentTime.Add(-1 * time.Minute),
			resetAt:       resetTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: true, // canary already active
		}

		result := isResetTimeExceeded(currentTime)(openState)

		assert.True(t, option.IsNone(result), "should return None when canary is already active")
	})

	t.Run("returns Some at exact reset time boundary", func(t *testing.T) {
		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		resetTime := currentTime.Add(-1 * time.Nanosecond) // just passed

		openState := openState{
			openedAt:      currentTime.Add(-1 * time.Minute),
			resetAt:       resetTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		result := isResetTimeExceeded(currentTime)(openState)

		assert.True(t, option.IsSome(result), "should return Some when current time is after reset time")
	})

	t.Run("returns None when current time equals reset time", func(t *testing.T) {
		currentTime := time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC)
		resetTime := currentTime // exactly equal

		openState := openState{
			openedAt:      currentTime.Add(-1 * time.Minute),
			resetAt:       resetTime,
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}

		result := isResetTimeExceeded(currentTime)(openState)

		assert.True(t, option.IsNone(result), "should return None when times are equal (not After)")
	})
}

// TestHandleSuccessOnClosed tests the handleSuccessOnClosed function
func TestHandleSuccessOnClosed(t *testing.T) {
	t.Run("updates closed state with success when circuit is closed", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		// Create a simple addSuccess reader that increments a counter
		addSuccess := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddSuccess(ct)
			}
		}

		// Create initial closed state
		initialClosed := MakeClosedStateCounter(3)
		initialState := createClosedCircuit(initialClosed)

		// Apply handleSuccessOnClosed
		handler := handleSuccessOnClosed(addSuccess)
		endomorphism := handler(currentTime)
		result := endomorphism(initialState)

		// Verify the state is still closed
		assert.True(t, IsClosed(result), "state should remain closed after success")

		// Verify the closed state was updated
		closedState := either.Fold(
			func(openState) ClosedState { return initialClosed },
			F.Identity[ClosedState],
		)(result)
		// The success should have been recorded (implementation-specific verification)
		assert.NotNil(t, closedState, "closed state should be present")
	})

	t.Run("does not affect open state", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		addSuccess := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddSuccess(ct)
			}
		}

		// Create initial open state
		initialOpen := openState{
			openedAt:      currentTime.Add(-1 * time.Minute),
			resetAt:       currentTime.Add(1 * time.Minute),
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: false,
		}
		initialState := createOpenCircuit(initialOpen)

		// Apply handleSuccessOnClosed
		handler := handleSuccessOnClosed(addSuccess)
		endomorphism := handler(currentTime)
		result := endomorphism(initialState)

		// Verify the state remains open and unchanged
		assert.True(t, IsOpen(result), "state should remain open")

		// Extract and verify the open state is unchanged
		openResult := either.Fold(
			func(os openState) openState { return os },
			func(ClosedState) openState { return initialOpen },
		)(result)
		assert.Equal(t, initialOpen.openedAt, openResult.openedAt, "openedAt should be unchanged")
		assert.Equal(t, initialOpen.resetAt, openResult.resetAt, "resetAt should be unchanged")
		assert.Equal(t, initialOpen.canaryRequest, openResult.canaryRequest, "canaryRequest should be unchanged")
	})

	t.Run("preserves time parameter through reader", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		time1 := vt.Now()
		vt.Advance(1 * time.Hour)
		time2 := vt.Now()

		var capturedTime time.Time
		addSuccess := func(ct time.Time) Endomorphism[ClosedState] {
			capturedTime = ct
			return F.Identity[ClosedState]
		}

		initialClosed := MakeClosedStateCounter(3)
		initialState := createClosedCircuit(initialClosed)

		handler := handleSuccessOnClosed(addSuccess)

		// Apply with time1
		endomorphism1 := handler(time1)
		endomorphism1(initialState)
		assert.Equal(t, time1, capturedTime, "should pass time1 to addSuccess")

		// Apply with time2
		endomorphism2 := handler(time2)
		endomorphism2(initialState)
		assert.Equal(t, time2, capturedTime, "should pass time2 to addSuccess")
	})

	t.Run("composes correctly with multiple successes", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		addSuccess := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddSuccess(ct)
			}
		}

		initialClosed := MakeClosedStateCounter(3)
		initialState := createClosedCircuit(initialClosed)

		handler := handleSuccessOnClosed(addSuccess)
		endomorphism := handler(currentTime)

		// Apply multiple times
		result1 := endomorphism(initialState)
		result2 := endomorphism(result1)
		result3 := endomorphism(result2)

		// All should remain closed
		assert.True(t, IsClosed(result1), "state should remain closed after first success")
		assert.True(t, IsClosed(result2), "state should remain closed after second success")
		assert.True(t, IsClosed(result3), "state should remain closed after third success")
	})
}

// TestHandleFailureOnClosed tests the handleFailureOnClosed function
func TestHandleFailureOnClosed(t *testing.T) {
	t.Run("keeps circuit closed when threshold not exceeded", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		// Create a closed state that allows 3 errors
		initialClosed := MakeClosedStateCounter(3)

		// addError increments error count
		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		// checkClosedState returns Some if under threshold
		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		// openCircuit creates an open state (shouldn't be called in this test)
		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(initialClosed)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// First error - should stay closed
		result1 := endomorphism(initialState)
		assert.True(t, IsClosed(result1), "circuit should remain closed after first error")

		// Second error - should stay closed
		result2 := endomorphism(result1)
		assert.True(t, IsClosed(result2), "circuit should remain closed after second error")
	})

	t.Run("opens circuit when threshold exceeded", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		// Create a closed state that allows only 2 errors (opens at 2nd error)
		initialClosed := MakeClosedStateCounter(2)

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(initialClosed)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// First error - should stay closed (count=1, threshold=2)
		result1 := endomorphism(initialState)
		assert.True(t, IsClosed(result1), "circuit should remain closed after first error")

		// Second error - should open (count=2, threshold=2)
		result2 := endomorphism(result1)
		assert.True(t, IsOpen(result2), "circuit should open when threshold reached")
	})

	t.Run("creates open state with correct reset time", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()
		expectedResetTime := currentTime.Add(5 * time.Minute)

		initialClosed := MakeClosedStateCounter(1) // Opens at 1st error

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       expectedResetTime,
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(initialClosed)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// First error - should open immediately (threshold=1)
		result1 := endomorphism(initialState)
		assert.True(t, IsOpen(result1), "circuit should open after first error")

		// Verify the open state has correct reset time
		resultOpen := either.Fold(
			func(os openState) openState { return os },
			func(ClosedState) openState { return openState{} },
		)(result1)
		assert.Equal(t, expectedResetTime, resultOpen.resetAt, "reset time should match expected")
		assert.Equal(t, currentTime, resultOpen.openedAt, "opened time should be current time")
	})

	t.Run("edge case: zero error threshold", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		// Create a closed state that allows 0 errors (opens immediately)
		initialClosed := MakeClosedStateCounter(0)

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(initialClosed)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// First error should immediately open the circuit
		result := endomorphism(initialState)
		assert.True(t, IsOpen(result), "circuit should open immediately with zero threshold")
	})

	t.Run("edge case: very high error threshold", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		// Create a closed state that allows 1000 errors
		initialClosed := MakeClosedStateCounter(1000)

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(initialClosed)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// Apply many errors
		result := initialState
		for i := 0; i < 100; i++ {
			result = endomorphism(result)
		}

		// Should still be closed after 100 errors
		assert.True(t, IsClosed(result), "circuit should remain closed with high threshold")
	})

	t.Run("preserves time parameter through reader chain", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		time1 := vt.Now()
		vt.Advance(2 * time.Hour)
		time2 := vt.Now()

		var capturedAddErrorTime, capturedCheckTime, capturedOpenTime time.Time

		initialClosed := MakeClosedStateCounter(2) // Need 2 errors to open

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			capturedAddErrorTime = ct
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			capturedCheckTime = ct
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			capturedOpenTime = ct
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(initialClosed)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)

		// Apply with time1 - first error, stays closed
		endomorphism1 := handler(time1)
		result1 := endomorphism1(initialState)
		assert.Equal(t, time1, capturedAddErrorTime, "addError should receive time1")
		assert.Equal(t, time1, capturedCheckTime, "checkClosedState should receive time1")

		// Apply with time2 - second error, should trigger open
		endomorphism2 := handler(time2)
		endomorphism2(result1)
		assert.Equal(t, time2, capturedAddErrorTime, "addError should receive time2")
		assert.Equal(t, time2, capturedCheckTime, "checkClosedState should receive time2")
		assert.Equal(t, time2, capturedOpenTime, "openCircuit should receive time2")
	})

	t.Run("handles transition from closed to open correctly", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		initialClosed := MakeClosedStateCounter(2) // Opens at 2nd error

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// Start with closed state
		state := createClosedCircuit(initialClosed)
		assert.True(t, IsClosed(state), "initial state should be closed")

		// First error - should stay closed (count=1, threshold=2)
		state = endomorphism(state)
		assert.True(t, IsClosed(state), "should remain closed after first error")

		// Second error - should open (count=2, threshold=2)
		state = endomorphism(state)
		assert.True(t, IsOpen(state), "should open after second error")

		// Verify it's truly open with correct properties
		resultOpen := either.Fold(
			func(os openState) openState { return os },
			func(ClosedState) openState { return openState{} },
		)(state)
		assert.False(t, resultOpen.canaryRequest, "canaryRequest should be false initially")
		assert.Equal(t, currentTime, resultOpen.openedAt, "openedAt should be current time")
	})

	t.Run("does not affect already open state", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now()

		addError := func(ct time.Time) Endomorphism[ClosedState] {
			return func(cs ClosedState) ClosedState {
				return cs.AddError(ct)
			}
		}

		checkClosedState := func(ct time.Time) option.Kleisli[ClosedState, ClosedState] {
			return func(cs ClosedState) Option[ClosedState] {
				return cs.Check(ct)
			}
		}

		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		// Start with an already open state
		existingOpen := openState{
			openedAt:      currentTime.Add(-5 * time.Minute),
			resetAt:       currentTime.Add(5 * time.Minute),
			retryStatus:   retry.DefaultRetryStatus,
			canaryRequest: true,
		}
		initialState := createOpenCircuit(existingOpen)

		handler := handleFailureOnClosed(addError, checkClosedState, openCircuit)
		endomorphism := handler(currentTime)

		// Apply to open state - should not change it
		result := endomorphism(initialState)

		assert.True(t, IsOpen(result), "state should remain open")

		// The open state should be unchanged since handleFailureOnClosed
		// only operates on the Right (closed) side of the Either
		openResult := either.Fold(
			func(os openState) openState { return os },
			func(ClosedState) openState { return openState{} },
		)(result)
		assert.Equal(t, existingOpen.openedAt, openResult.openedAt, "openedAt should be unchanged")
		assert.Equal(t, existingOpen.resetAt, openResult.resetAt, "resetAt should be unchanged")
		assert.Equal(t, existingOpen.canaryRequest, openResult.canaryRequest, "canaryRequest should be unchanged")
	})
}
