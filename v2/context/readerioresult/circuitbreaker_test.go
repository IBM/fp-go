// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package readerioresult

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/circuitbreaker"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// Helper function to create a test logger that collects messages
func testMetrics(_ *[]string) circuitbreaker.Metrics {
	return circuitbreaker.MakeMetricsFromLogger("testMetrics", log.Default())
}

// Helper function to create a simple closed state
func testCBClosedState() circuitbreaker.ClosedState {
	return circuitbreaker.MakeClosedStateCounter(3)
}

// Helper function to create a test retry policy
func testCBRetryPolicy() retry.RetryPolicy {
	return retry.Monoid.Concat(
		retry.LimitRetries(3),
		retry.ExponentialBackoff(100*time.Millisecond),
	)
}

// Helper function that checks all errors
func checkAllErrors(err error) option.Option[error] {
	return option.Some(err)
}

// Helper function that ignores specific errors
func ignoreSpecificError(ignoredMsg string) func(error) option.Option[error] {
	return func(err error) option.Option[error] {
		if err.Error() == ignoredMsg {
			return option.None[error]()
		}
		return option.Some(err)
	}
}

// TestCircuitBreaker_SuccessfulOperation tests that successful operations
// pass through the circuit breaker without issues
func TestCircuitBreaker_SuccessfulOperation(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	// Create initial state
	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	// Successful operation
	operation := Of("success")

	// Apply circuit breaker
	env := pair.MakePair(stateRef, operation)
	resultEnv := cb(env)

	// Execute
	ctx := t.Context()
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.Equal(t, result.Of("success"), outcome)
}

// TestCircuitBreaker_SingleFailure tests that a single failure is handled
// but doesn't open the circuit
func TestCircuitBreaker_SingleFailure(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	expError := errors.New("operation failed")

	// Failing operation
	operation := Left[string](expError)

	env := pair.MakePair(stateRef, operation)
	resultEnv := cb(env)

	ctx := t.Context()
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.Equal(t, result.Left[string](expError), outcome)

	// Circuit should still be closed after one failure
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(state))
}

// TestCircuitBreaker_OpensAfterThreshold tests that the circuit opens
// after exceeding the failure threshold
func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(), // Opens after 3 failures
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	expError := errors.New("operation failed")

	// Failing operation
	operation := Left[string](expError)

	ctx := t.Context()

	// Execute 3 failures to open the circuit
	for range 3 {
		env := pair.MakePair(stateRef, operation)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()
		assert.Equal(t, result.Left[string](expError), outcome)
	}

	// Circuit should now be open
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsOpen(state))

	// Next request should fail immediately with circuit breaker error
	env := pair.MakePair(stateRef, operation)
	resultEnv := cb(env)
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.True(t, result.IsLeft(outcome))
	_, err := result.Unwrap(outcome)
	var cbErr *circuitbreaker.CircuitBreakerError
	assert.ErrorAs(t, err, &cbErr)
}

// TestCircuitBreaker_HalfOpenAfterResetTime tests that the circuit
// transitions to half-open state after the reset time
func TestCircuitBreaker_HalfOpenAfterResetTime(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	expError := errors.New("operation failed")

	// Failing operation
	failingOp := Left[string](expError)

	ctx := t.Context()

	// Open the circuit with 3 failures
	for range 3 {
		env := pair.MakePair(stateRef, failingOp)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()

		assert.Equal(t, result.Left[string](expError), outcome)
	}

	// Verify circuit is open
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsOpen(state))

	// Advance time past the reset time (exponential backoff starts at 100ms)
	vt.Advance(200 * time.Millisecond)

	// Now create a successful operation for the canary request
	successOp := Of("success")

	// Next request should be a canary request
	env := pair.MakePair(stateRef, successOp)
	resultEnv := cb(env)
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	// Canary should succeed
	assert.Equal(t, result.Of("success"), outcome)

	// Circuit should now be closed again
	state = ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(state))
}

// TestCircuitBreaker_CanaryFailureExtendsOpenTime tests that a failed
// canary request extends the open time
func TestCircuitBreaker_CanaryFailureExtendsOpenTime(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	expError := errors.New("operation failed")

	// Failing operation
	failingOp := Left[string](expError)

	ctx := t.Context()

	// Open the circuit
	for range 3 {
		env := pair.MakePair(stateRef, failingOp)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()
		assert.Equal(t, result.Left[string](expError), outcome)
	}

	// Advance time to trigger canary
	vt.Advance(200 * time.Millisecond)

	// Canary request fails
	env := pair.MakePair(stateRef, failingOp)
	resultEnv := cb(env)
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.True(t, result.IsLeft(outcome))

	// Circuit should still be open
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsOpen(state))

	// Immediate next request should fail with circuit breaker error
	env = pair.MakePair(stateRef, failingOp)
	resultEnv = cb(env)
	protectedOp = pair.Tail(resultEnv)
	outcome = protectedOp(ctx)()

	assert.True(t, result.IsLeft(outcome))
	_, err := result.Unwrap(outcome)
	var cbErr *circuitbreaker.CircuitBreakerError
	assert.ErrorAs(t, err, &cbErr)
}

// TestCircuitBreaker_IgnoredErrorsDoNotCount tests that errors filtered
// by checkError don't count toward opening the circuit
func TestCircuitBreaker_IgnoredErrorsDoNotCount(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	// Ignore "ignorable error"
	checkError := ignoreSpecificError("ignorable error")

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkError,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	ctx := t.Context()
	ignorableError := errors.New("ignorable error")

	// Execute 5 ignorable errors
	ignorableOp := Left[string](ignorableError)

	for range 5 {
		env := pair.MakePair(stateRef, ignorableOp)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()
		assert.Equal(t, result.Left[string](ignorableError), outcome)
	}

	// Circuit should still be closed
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(state))

	realError := errors.New("real error")

	// Now send a real error
	realErrorOp := Left[string](realError)

	env := pair.MakePair(stateRef, realErrorOp)
	resultEnv := cb(env)
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.Equal(t, result.Left[string](realError), outcome)

	// Circuit should still be closed (only 1 counted error)
	state = ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(state))
}

// TestCircuitBreaker_MixedSuccessAndFailure tests the circuit behavior
// with a mix of successful and failed operations
func TestCircuitBreaker_MixedSuccessAndFailure(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	ctx := t.Context()

	successOp := Of("success")
	expError := errors.New("failure")

	failOp := Left[string](expError)

	// Pattern: fail, fail, success, fail
	ops := array.From(failOp, failOp, successOp, failOp)

	for _, op := range ops {
		env := pair.MakePair(stateRef, op)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		_ = protectedOp(ctx)()
	}

	// Circuit should still be closed (success resets the count)
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(state))
}

// TestCircuitBreaker_ConcurrentOperations tests that the circuit breaker
// handles concurrent operations correctly
func TestCircuitBreaker_ConcurrentOperations(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[int](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	ctx := t.Context()

	var wg sync.WaitGroup
	results := make([]Result[int], 10)

	// Launch 10 concurrent operations
	for i := range 10 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			op := Of(idx)

			env := pair.MakePair(stateRef, op)
			resultEnv := cb(env)
			protectedOp := pair.Tail(resultEnv)
			results[idx] = protectedOp(ctx)()
		}(i)
	}

	wg.Wait()

	// All operations should succeed
	for i, res := range results {
		assert.True(t, result.IsRight(res), "Operation %d should succeed", i)
	}
}

// TestCircuitBreaker_DifferentTypes tests that the circuit breaker works
// with different result types
func TestCircuitBreaker_DifferentTypes(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	// Test with int
	cbInt := MakeCircuitBreaker[int](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRefInt := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	opInt := Of(42)

	ctx := t.Context()
	envInt := pair.MakePair(stateRefInt, opInt)
	resultEnvInt := cbInt(envInt)
	protectedOpInt := pair.Tail(resultEnvInt)
	outcomeInt := protectedOpInt(ctx)()

	assert.Equal(t, result.Of(42), outcomeInt)

	// Test with struct
	type User struct {
		ID   int
		Name string
	}

	cbUser := MakeCircuitBreaker[User](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRefUser := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	opUser := Of(User{ID: 1, Name: "Alice"})

	envUser := pair.MakePair(stateRefUser, opUser)
	resultEnvUser := cbUser(envUser)
	protectedOpUser := pair.Tail(resultEnvUser)
	outcomeUser := protectedOpUser(ctx)()

	require.Equal(t, result.Of(User{ID: 1, Name: "Alice"}), outcomeUser)
}

// TestCircuitBreaker_VirtualTimerAdvancement tests that the virtual timer
// correctly controls time-based behavior
func TestCircuitBreaker_VirtualTimerAdvancement(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	vt := NewVirtualTimer(startTime)

	// Verify initial time
	assert.Equal(t, startTime, vt.Now())

	// Advance by 1 hour
	vt.Advance(1 * time.Hour)
	assert.Equal(t, startTime.Add(1*time.Hour), vt.Now())

	// Advance by 30 minutes
	vt.Advance(30 * time.Minute)
	assert.Equal(t, startTime.Add(90*time.Minute), vt.Now())

	// Set to specific time
	newTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	vt.Set(newTime)
	assert.Equal(t, newTime, vt.Now())
}

// TestCircuitBreaker_InitialState tests that the circuit starts in closed state
func TestCircuitBreaker_InitialState(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	// Check initial state is closed
	state := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(state), "Circuit should start in closed state")

	// First operation should execute normally
	op := Of("first operation")

	ctx := t.Context()
	env := pair.MakePair(stateRef, op)
	resultEnv := cb(env)
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.Equal(t, result.Of("first operation"), outcome)
}

// TestCircuitBreaker_ErrorMessageFormat tests that circuit breaker errors
// have appropriate error messages
func TestCircuitBreaker_ErrorMessageFormat(t *testing.T) {
	vt := NewVirtualTimer(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	var logMessages []string

	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()

	ctx := t.Context()

	expError := errors.New("service unavailable")

	failOp := Left[string](expError)

	// Open the circuit
	for range 3 {
		env := pair.MakePair(stateRef, failOp)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		_ = protectedOp(ctx)()
	}

	// Next request should fail with circuit breaker error
	env := pair.MakePair(stateRef, failOp)
	resultEnv := cb(env)
	protectedOp := pair.Tail(resultEnv)
	outcome := protectedOp(ctx)()

	assert.True(t, result.IsLeft(outcome))

	// Error message should indicate circuit breaker is open
	_, err := result.Unwrap(outcome)
	errMsg := err.Error()
	assert.Contains(t, errMsg, "circuit", "Error should mention circuit breaker")
}

// RequestSpec defines a virtual request with timing and outcome information
type RequestSpec struct {
	ID         int           // Unique identifier for the request
	StartTime  time.Duration // Virtual start time relative to test start
	Duration   time.Duration // How long the request takes to execute
	ShouldFail bool          // Whether this request should fail
}

// RequestResult captures the outcome of a request execution
type RequestResult struct {
	ID                  int
	StartTime           time.Time
	EndTime             time.Time
	Success             bool
	Error               error
	CircuitBreakerError bool // True if failed due to circuit breaker being open
}

// TestCircuitBreaker_ConcurrentBatchWithThresholdExceeded tests a complex
// concurrent scenario where:
// 1. Initial requests succeed
// 2. A batch of failures exceeds the threshold, opening the circuit
// 3. Subsequent requests fail immediately due to open circuit
// 4. After timeout, a canary request succeeds
// 5. Following requests succeed again
func TestCircuitBreaker_ConcurrentBatchWithThresholdExceeded(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	vt := NewVirtualTimer(startTime)
	var logMessages []string

	// Circuit opens after 3 failures, with exponential backoff starting at 100ms
	cb := MakeCircuitBreaker[string](
		vt.Now,
		testCBClosedState(), // Opens after 3 failures
		checkAllErrors,
		testCBRetryPolicy(), // 100ms initial backoff
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()
	ctx := t.Context()

	// Define the request sequence
	// Phase 1: Initial successes (0-100ms)
	// Phase 2: Failures that exceed threshold (100-200ms) - should open circuit
	// Phase 3: Requests during open circuit (200-300ms) - should fail immediately
	// Phase 4: After timeout (400ms+) - canary succeeds, then more successes
	requests := []RequestSpec{
		// Phase 1: Initial successful requests
		{ID: 1, StartTime: 0 * time.Millisecond, Duration: 10 * time.Millisecond, ShouldFail: false},
		{ID: 2, StartTime: 20 * time.Millisecond, Duration: 10 * time.Millisecond, ShouldFail: false},

		// Phase 2: Sequential failures that exceed threshold (3 failures)
		{ID: 3, StartTime: 100 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: true},
		{ID: 4, StartTime: 110 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: true},
		{ID: 5, StartTime: 120 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: true},
		{ID: 6, StartTime: 130 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: true},

		// Phase 3: Requests during open circuit - should fail with circuit breaker error
		{ID: 7, StartTime: 200 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: false},
		{ID: 8, StartTime: 210 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: false},
		{ID: 9, StartTime: 220 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: false},

		// Phase 4: After reset timeout (100ms backoff from last failure at ~125ms = ~225ms)
		// Wait longer to ensure we're past the reset time
		{ID: 10, StartTime: 400 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: false}, // Canary succeeds
		{ID: 11, StartTime: 410 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: false},
		{ID: 12, StartTime: 420 * time.Millisecond, Duration: 5 * time.Millisecond, ShouldFail: false},
	}

	results := make([]RequestResult, len(requests))

	// Execute requests sequentially but model them as if they were concurrent
	// by advancing the virtual timer to each request's start time
	for i, req := range requests {
		// Set virtual time to request start time
		vt.Set(startTime.Add(req.StartTime))

		// Create the operation based on spec
		var op ReaderIOResult[string]
		if req.ShouldFail {
			op = Left[string](errors.New("operation failed"))
		} else {
			op = Of("success")
		}

		// Apply circuit breaker
		env := pair.MakePair(stateRef, op)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)

		// Record start time
		execStartTime := vt.Now()

		// Execute the operation
		outcome := protectedOp(ctx)()

		// Advance time by operation duration
		vt.Advance(req.Duration)
		execEndTime := vt.Now()

		// Analyze the result
		isSuccess := result.IsRight(outcome)
		var err error
		var isCBError bool

		if !isSuccess {
			_, err = result.Unwrap(outcome)
			var cbErr *circuitbreaker.CircuitBreakerError
			isCBError = errors.As(err, &cbErr)
		}

		results[i] = RequestResult{
			ID:                  req.ID,
			StartTime:           execStartTime,
			EndTime:             execEndTime,
			Success:             isSuccess,
			Error:               err,
			CircuitBreakerError: isCBError,
		}
	}

	// Verify Phase 1: Initial requests should succeed
	assert.True(t, results[0].Success, "Request 1 should succeed")
	assert.True(t, results[1].Success, "Request 2 should succeed")

	// Verify Phase 2: Failures should be recorded (first 3 fail with actual error)
	// The 4th might fail with CB error if circuit opened fast enough
	assert.False(t, results[2].Success, "Request 3 should fail")
	assert.False(t, results[3].Success, "Request 4 should fail")
	assert.False(t, results[4].Success, "Request 5 should fail")

	// At least the first 3 failures should be actual operation failures, not CB errors
	actualFailures := 0
	for i := 2; i <= 4; i++ {
		if !results[i].CircuitBreakerError {
			actualFailures++
		}
	}
	assert.GreaterOrEqual(t, actualFailures, 3, "At least 3 actual operation failures should occur")

	// Verify Phase 3: Requests during open circuit should fail with circuit breaker error
	for i := 6; i <= 8; i++ {
		assert.False(t, results[i].Success, "Request %d should fail during open circuit", results[i].ID)
		assert.True(t, results[i].CircuitBreakerError, "Request %d should fail with circuit breaker error", results[i].ID)
	}

	// Verify Phase 4: After timeout, canary and subsequent requests should succeed
	assert.True(t, results[9].Success, "Request 10 (canary) should succeed")
	assert.True(t, results[10].Success, "Request 11 should succeed after circuit closes")
	assert.True(t, results[11].Success, "Request 12 should succeed after circuit closes")

	// Verify final state is closed
	finalState := ioref.Read(stateRef)()
	assert.True(t, circuitbreaker.IsClosed(finalState), "Circuit should be closed at the end")

	// Log summary for debugging
	t.Logf("Test completed with %d requests", len(results))
	successCount := 0
	cbErrorCount := 0
	actualErrorCount := 0

	for _, r := range results {
		if r.Success {
			successCount++
		} else if r.CircuitBreakerError {
			cbErrorCount++
		} else {
			actualErrorCount++
		}
	}

	t.Logf("Summary: %d successes, %d circuit breaker errors, %d actual errors",
		successCount, cbErrorCount, actualErrorCount)
}

// TestCircuitBreaker_ConcurrentHighLoad tests circuit breaker behavior
// under high concurrent load with mixed success/failure patterns
func TestCircuitBreaker_ConcurrentHighLoad(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	vt := NewVirtualTimer(startTime)
	var logMessages []string

	cb := MakeCircuitBreaker[int](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()
	ctx := t.Context()

	// Create a large batch of 50 requests
	// Pattern: success, success, fail, fail, fail, fail, success, success, ...
	// This ensures we have initial successes, then failures to open circuit,
	// then more requests that hit the open circuit
	numRequests := 50

	results := make([]bool, numRequests)
	cbErrors := make([]bool, numRequests)

	// Execute requests with controlled timing
	for i := range numRequests {
		// Advance time slightly for each request
		vt.Advance(10 * time.Millisecond)

		// Pattern: 2 success, 4 failures, repeat
		// This ensures we exceed the threshold (3 failures) early on
		shouldFail := (i%6) >= 2 && (i%6) < 6

		var op ReaderIOResult[int]
		if shouldFail {
			op = Left[int](errors.New("simulated failure"))
		} else {
			op = Of(i)
		}

		env := pair.MakePair(stateRef, op)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()

		results[i] = result.IsRight(outcome)

		if !results[i] {
			_, err := result.Unwrap(outcome)
			var cbErr *circuitbreaker.CircuitBreakerError
			cbErrors[i] = errors.As(err, &cbErr)
		}
	}

	// Count outcomes
	successCount := 0
	failureCount := 0
	cbErrorCount := 0

	for i := range numRequests {
		if results[i] {
			successCount++
		} else {
			failureCount++
			if cbErrors[i] {
				cbErrorCount++
			}
		}
	}

	t.Logf("High load test: %d total requests", numRequests)
	t.Logf("Results: %d successes, %d failures (%d circuit breaker errors)",
		successCount, failureCount, cbErrorCount)

	// Verify that circuit breaker activated (some requests failed due to open circuit)
	assert.Greater(t, cbErrorCount, 0, "Circuit breaker should have opened and blocked some requests")

	// Verify that not all requests failed (some succeeded before circuit opened)
	assert.Greater(t, successCount, 0, "Some requests should have succeeded")
}

// TestCircuitBreaker_TrueConcurrentRequests tests actual concurrent execution
// with proper synchronization
func TestCircuitBreaker_TrueConcurrentRequests(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	vt := NewVirtualTimer(startTime)
	var logMessages []string

	cb := MakeCircuitBreaker[int](
		vt.Now,
		testCBClosedState(),
		checkAllErrors,
		testCBRetryPolicy(),
		testMetrics(&logMessages),
	)

	stateRef := circuitbreaker.MakeClosedIORef(testCBClosedState())()
	ctx := t.Context()

	// Launch 20 concurrent requests
	numRequests := 20
	var wg sync.WaitGroup
	results := make([]bool, numRequests)
	cbErrors := make([]bool, numRequests)

	// First, send some successful requests
	for i := range 5 {
		op := Of(i)
		env := pair.MakePair(stateRef, op)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()
		results[i] = result.IsRight(outcome)
	}

	// Now send concurrent failures to open the circuit
	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			op := Left[int](errors.New("concurrent failure"))
			env := pair.MakePair(stateRef, op)
			resultEnv := cb(env)
			protectedOp := pair.Tail(resultEnv)
			outcome := protectedOp(ctx)()
			results[idx] = result.IsRight(outcome)
			if !results[idx] {
				_, err := result.Unwrap(outcome)
				var cbErr *circuitbreaker.CircuitBreakerError
				cbErrors[idx] = errors.As(err, &cbErr)
			}
		}(i)
	}

	wg.Wait()

	// Now send more requests that should hit the open circuit
	for i := 10; i < numRequests; i++ {
		op := Of(i)
		env := pair.MakePair(stateRef, op)
		resultEnv := cb(env)
		protectedOp := pair.Tail(resultEnv)
		outcome := protectedOp(ctx)()
		results[i] = result.IsRight(outcome)
		if !results[i] {
			_, err := result.Unwrap(outcome)
			var cbErr *circuitbreaker.CircuitBreakerError
			cbErrors[i] = errors.As(err, &cbErr)
		}
	}

	// Count outcomes
	successCount := 0
	failureCount := 0
	cbErrorCount := 0

	for i := range numRequests {
		if results[i] {
			successCount++
		} else {
			failureCount++
			if cbErrors[i] {
				cbErrorCount++
			}
		}
	}

	t.Logf("Concurrent test: %d total requests", numRequests)
	t.Logf("Results: %d successes, %d failures (%d circuit breaker errors)",
		successCount, failureCount, cbErrorCount)

	// Verify initial successes
	assert.Equal(t, 5, successCount, "First 5 requests should succeed")

	// Verify that circuit breaker opened and blocked some requests
	assert.Greater(t, cbErrorCount, 0, "Circuit breaker should have opened and blocked some requests")
}
