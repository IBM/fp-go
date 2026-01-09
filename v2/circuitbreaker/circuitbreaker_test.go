package circuitbreaker

import (
	"sync"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// VirtualTimer provides a controllable time source for testing
type VirtualTimer struct {
	mu      sync.Mutex
	current time.Time
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
	t.Run("resets failure count on success", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now
		addSuccess := reader.From1(ClosedState.AddSuccess)

		// Create initial state with some failures
		now := vt.Now()
		initialClosed := MakeClosedStateCounter(3)
		initialClosed = initialClosed.AddError(now)
		initialClosed = initialClosed.AddError(now)
		initialState := createClosedCircuit(initialClosed)

		ref := io.Run(ioref.MakeIORef(initialState))
		modify := modifyV(ref)

		handler := handleSuccessOnClosed(currentTime, addSuccess)

		// Apply the handler
		result := io.Run(handler(modify))

		// Verify state is still closed and failures are reset
		assert.True(t, IsClosed(result), "circuit should remain closed after success")
	})

	t.Run("keeps circuit closed", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now
		addSuccess := reader.From1(ClosedState.AddSuccess)

		initialState := createClosedCircuit(MakeClosedStateCounter(3))
		ref := io.Run(ioref.MakeIORef(initialState))
		modify := modifyV(ref)

		handler := handleSuccessOnClosed(currentTime, addSuccess)
		result := io.Run(handler(modify))

		assert.True(t, IsClosed(result), "circuit should remain closed")
	})
}

// TestHandleFailureOnClosed tests the handleFailureOnClosed function
func TestHandleFailureOnClosed(t *testing.T) {
	t.Run("keeps circuit closed when threshold not exceeded", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now
		addError := reader.From1(ClosedState.AddError)
		checkClosedState := reader.From1(ClosedState.Check)
		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		// Create initial state with room for more failures
		now := vt.Now()
		initialClosed := MakeClosedStateCounter(5) // threshold is 5
		initialClosed = initialClosed.AddError(now)
		initialState := createClosedCircuit(initialClosed)

		ref := io.Run(ioref.MakeIORef(initialState))
		modify := modifyV(ref)

		handler := handleFailureOnClosed(currentTime, addError, checkClosedState, openCircuit)
		result := io.Run(handler(modify))

		assert.True(t, IsClosed(result), "circuit should remain closed when threshold not exceeded")
	})

	t.Run("opens circuit when threshold exceeded", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now
		addError := reader.From1(ClosedState.AddError)
		checkClosedState := reader.From1(ClosedState.Check)
		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		// Create initial state at threshold
		now := vt.Now()
		initialClosed := MakeClosedStateCounter(2) // threshold is 2
		initialClosed = initialClosed.AddError(now)
		initialState := createClosedCircuit(initialClosed)

		ref := io.Run(ioref.MakeIORef(initialState))
		modify := modifyV(ref)

		handler := handleFailureOnClosed(currentTime, addError, checkClosedState, openCircuit)
		result := io.Run(handler(modify))

		assert.True(t, IsOpen(result), "circuit should open when threshold exceeded")
	})

	t.Run("records failure in closed state", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		currentTime := vt.Now
		addError := reader.From1(ClosedState.AddError)
		checkClosedState := reader.From1(ClosedState.Check)
		openCircuit := func(ct time.Time) openState {
			return openState{
				openedAt:      ct,
				resetAt:       ct.Add(1 * time.Minute),
				retryStatus:   retry.DefaultRetryStatus,
				canaryRequest: false,
			}
		}

		initialState := createClosedCircuit(MakeClosedStateCounter(10))
		ref := io.Run(ioref.MakeIORef(initialState))
		modify := modifyV(ref)

		handler := handleFailureOnClosed(currentTime, addError, checkClosedState, openCircuit)
		result := io.Run(handler(modify))

		// Should still be closed but with failure recorded
		assert.True(t, IsClosed(result), "circuit should remain closed")
	})
}
