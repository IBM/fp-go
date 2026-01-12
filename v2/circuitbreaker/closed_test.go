package circuitbreaker

import (
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	"github.com/stretchr/testify/assert"
)

func TestMakeClosedStateCounter(t *testing.T) {
	t.Run("creates a valid ClosedState", func(t *testing.T) {
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)

		assert.NotNil(t, state, "MakeClosedStateCounter should return a non-nil ClosedState")
	})

	t.Run("initial state passes Check", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		result := state.Check(now)

		assert.True(t, option.IsSome(result), "initial state should pass Check (return Some, circuit stays closed)")
	})

	t.Run("Empty resets failure count", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(2)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Add some errors
		state = state.AddError(now)
		state = state.AddError(now)

		// Reset the state
		state = state.Empty()

		// Should pass check after reset
		result := state.Check(now)
		assert.True(t, option.IsSome(result), "state should pass Check after Empty")
	})

	t.Run("AddSuccess resets failure count", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)

		// Add errors
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)

		// Add success (should reset counter)
		state = state.AddSuccess(vt.Now())
		vt.Advance(1 * time.Second)

		// Add another error (this is now the first consecutive error)
		state = state.AddError(vt.Now())

		// Should still pass check (only 1 consecutive error, threshold is 3)
		result := state.Check(vt.Now())
		assert.True(t, option.IsSome(result), "AddSuccess should reset failure count")
	})

	t.Run("circuit opens when failures reach threshold", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Add errors up to but not including threshold
		state = state.AddError(now)
		state = state.AddError(now)

		// Should still pass before threshold
		result := state.Check(now)
		assert.True(t, option.IsSome(result), "should pass Check before threshold")

		// Add one more error to reach threshold
		state = state.AddError(now)

		// Should fail check at threshold
		result = state.Check(now)
		assert.True(t, option.IsNone(result), "should fail Check when reaching threshold")
	})

	t.Run("circuit opens exactly at maxFailures", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(5)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Add exactly maxFailures - 1 errors
		for i := uint(0); i < maxFailures-1; i++ {
			state = state.AddError(now)
		}

		// Should still pass
		result := state.Check(now)
		assert.True(t, option.IsSome(result), "should pass Check before maxFailures")

		// Add one more to reach maxFailures
		state = state.AddError(now)

		// Should fail now
		result = state.Check(now)
		assert.True(t, option.IsNone(result), "should fail Check at maxFailures")
	})

	t.Run("zero maxFailures means circuit is always open", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(0)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Initial state should already fail (0 >= 0)
		result := state.Check(now)
		assert.True(t, option.IsNone(result), "initial state should fail Check with maxFailures=0")

		// Add one error
		state = state.AddError(now)

		// Should still fail
		result = state.Check(now)
		assert.True(t, option.IsNone(result), "should fail Check after error with maxFailures=0")
	})

	t.Run("AddSuccess resets counter between errors", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)

		// Add errors
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)

		// Add success (resets counter)
		state = state.AddSuccess(vt.Now())
		vt.Advance(1 * time.Second)

		// Add more errors
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())

		// Should still pass (only 2 consecutive errors after reset)
		result := state.Check(vt.Now())
		assert.True(t, option.IsSome(result), "should pass with 2 consecutive errors after reset")

		// Add one more to reach threshold
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())

		// Should fail at threshold
		result = state.Check(vt.Now())
		assert.True(t, option.IsNone(result), "should fail after reaching threshold")
	})

	t.Run("Empty can be called multiple times", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(2)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Add errors
		state = state.AddError(now)
		state = state.AddError(now)
		state = state.AddError(now)

		// Reset multiple times
		state = state.Empty()
		state = state.Empty()
		state = state.Empty()

		// Should still pass
		result := state.Check(now)
		assert.True(t, option.IsSome(result), "state should pass Check after multiple Empty calls")
	})

	t.Run("time parameter is ignored in counter implementation", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)

		// Use different times for each operation
		time1 := vt.Now()
		time2 := time1.Add(1 * time.Hour)

		state = state.AddError(time1)
		state = state.AddError(time2)

		// Check with yet another time
		time3 := time1.Add(2 * time.Hour)
		result := state.Check(time3)

		// Should still pass (2 errors, threshold is 3, not reached yet)
		assert.True(t, option.IsSome(result), "time parameter should not affect counter behavior")

		// Add one more to reach threshold
		state = state.AddError(time1)
		result = state.Check(time1)
		assert.True(t, option.IsNone(result), "should fail after reaching threshold regardless of time")
	})

	t.Run("large maxFailures value", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(1000)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Add many errors but not reaching threshold
		for i := uint(0); i < maxFailures-1; i++ {
			state = state.AddError(now)
		}

		// Should still pass
		result := state.Check(now)
		assert.True(t, option.IsSome(result), "should pass Check with large maxFailures before threshold")

		// Add one more to reach threshold
		state = state.AddError(now)

		// Should fail
		result = state.Check(now)
		assert.True(t, option.IsNone(result), "should fail Check with large maxFailures at threshold")
	})

	t.Run("state is immutable - original unchanged after AddError", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(2)
		originalState := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Create new state by adding error
		newState := originalState.AddError(now)

		// Original should still pass check
		result := originalState.Check(now)
		assert.True(t, option.IsSome(result), "original state should be unchanged")

		// New state should reach threshold (2 errors total, threshold is 2)
		newState = newState.AddError(now)

		result = newState.Check(now)
		assert.True(t, option.IsNone(result), "new state should fail after reaching threshold")

		// Original should still pass
		result = originalState.Check(now)
		assert.True(t, option.IsSome(result), "original state should still be unchanged")
	})

	t.Run("state is immutable - original unchanged after Empty", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(2)
		state := MakeClosedStateCounter(maxFailures)
		now := vt.Now()

		// Add errors to original
		state = state.AddError(now)
		state = state.AddError(now)
		stateWithErrors := state

		// Create new state by calling Empty
		emptyState := stateWithErrors.Empty()

		// Original with errors should reach threshold (2 errors total, threshold is 2)
		result := stateWithErrors.Check(now)
		assert.True(t, option.IsNone(result), "state with errors should fail after reaching threshold")

		// Empty state should pass
		result = emptyState.Check(now)
		assert.True(t, option.IsSome(result), "empty state should pass Check")
	})

	t.Run("AddSuccess prevents circuit from opening", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)

		// Add errors close to threshold
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)

		// Add success before reaching threshold
		state = state.AddSuccess(vt.Now())
		vt.Advance(1 * time.Second)

		// Add more errors
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())

		// Should still pass (only 2 consecutive errors)
		result := state.Check(vt.Now())
		assert.True(t, option.IsSome(result), "circuit should stay closed after success reset")
	})

	t.Run("multiple AddSuccess calls keep counter at zero", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(2)
		state := MakeClosedStateCounter(maxFailures)

		// Add error
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)

		// Multiple successes
		state = state.AddSuccess(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddSuccess(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddSuccess(vt.Now())
		vt.Advance(1 * time.Second)

		// Should still pass
		result := state.Check(vt.Now())
		assert.True(t, option.IsSome(result), "multiple AddSuccess should keep counter at zero")

		// Add errors to reach threshold
		state = state.AddError(vt.Now())
		vt.Advance(1 * time.Second)
		state = state.AddError(vt.Now())

		// Should fail
		result = state.Check(vt.Now())
		assert.True(t, option.IsNone(result), "should fail after reaching threshold")
	})

	t.Run("alternating errors and successes never opens circuit", func(t *testing.T) {
		vt := NewVirtualTimer(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
		maxFailures := uint(3)
		state := MakeClosedStateCounter(maxFailures)

		// Alternate errors and successes
		for i := 0; i < 10; i++ {
			state = state.AddError(vt.Now())
			vt.Advance(500 * time.Millisecond)
			state = state.AddSuccess(vt.Now())
			vt.Advance(500 * time.Millisecond)
		}

		// Should still pass (never had consecutive failures)
		result := state.Check(vt.Now())
		assert.True(t, option.IsSome(result), "alternating errors and successes should never open circuit")
	})
}

func TestAddToSlice(t *testing.T) {
	ordTime := ord.OrdTime()

	t.Run("adds item to empty slice and returns sorted result", func(t *testing.T) {
		input := []time.Time{}
		item := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 1, "result should have 1 element")
		assert.Equal(t, item, result[0], "result should contain the added item")
	})

	t.Run("adds item and maintains sorted order", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		input := []time.Time{
			baseTime,
			baseTime.Add(20 * time.Second),
			baseTime.Add(40 * time.Second),
		}
		item := baseTime.Add(30 * time.Second)

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 4, "result should have 4 elements")
		// Verify sorted order
		assert.Equal(t, baseTime, result[0])
		assert.Equal(t, baseTime.Add(20*time.Second), result[1])
		assert.Equal(t, baseTime.Add(30*time.Second), result[2])
		assert.Equal(t, baseTime.Add(40*time.Second), result[3])
	})

	t.Run("adds item at beginning when it's earliest", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		input := []time.Time{
			baseTime.Add(20 * time.Second),
			baseTime.Add(40 * time.Second),
		}
		item := baseTime

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 3, "result should have 3 elements")
		assert.Equal(t, baseTime, result[0], "earliest item should be first")
		assert.Equal(t, baseTime.Add(20*time.Second), result[1])
		assert.Equal(t, baseTime.Add(40*time.Second), result[2])
	})

	t.Run("adds item at end when it's latest", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		input := []time.Time{
			baseTime,
			baseTime.Add(20 * time.Second),
		}
		item := baseTime.Add(40 * time.Second)

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 3, "result should have 3 elements")
		assert.Equal(t, baseTime, result[0])
		assert.Equal(t, baseTime.Add(20*time.Second), result[1])
		assert.Equal(t, baseTime.Add(40*time.Second), result[2], "latest item should be last")
	})

	t.Run("does not modify original slice", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		input := []time.Time{
			baseTime,
			baseTime.Add(20 * time.Second),
		}
		originalLen := len(input)
		originalFirst := input[0]
		originalLast := input[1]
		item := baseTime.Add(10 * time.Second)

		result := addToSlice(ordTime, input, item)

		// Verify original slice is unchanged
		assert.Len(t, input, originalLen, "original slice length should be unchanged")
		assert.Equal(t, originalFirst, input[0], "original slice first element should be unchanged")
		assert.Equal(t, originalLast, input[1], "original slice last element should be unchanged")

		// Verify result is different and has correct length
		assert.Len(t, result, 3, "result should have new length")
		// Verify the result contains the new item in sorted order
		assert.Equal(t, baseTime, result[0])
		assert.Equal(t, baseTime.Add(10*time.Second), result[1])
		assert.Equal(t, baseTime.Add(20*time.Second), result[2])
	})

	t.Run("handles duplicate timestamps", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		input := []time.Time{
			baseTime,
			baseTime.Add(20 * time.Second),
		}
		item := baseTime // duplicate of first element

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 3, "result should have 3 elements including duplicate")
		// Both instances of baseTime should be present
		count := 0
		for _, t := range result {
			if t.Equal(baseTime) {
				count++
			}
		}
		assert.Equal(t, 2, count, "should have 2 instances of the duplicate timestamp")
	})

	t.Run("maintains sort order with unsorted input", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		// Input is intentionally unsorted
		input := []time.Time{
			baseTime.Add(40 * time.Second),
			baseTime,
			baseTime.Add(20 * time.Second),
		}
		item := baseTime.Add(30 * time.Second)

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 4, "result should have 4 elements")
		// Verify result is sorted regardless of input order
		for i := 0; i < len(result)-1; i++ {
			assert.True(t, result[i].Before(result[i+1]) || result[i].Equal(result[i+1]),
				"result should be sorted: element %d (%v) should be <= element %d (%v)",
				i, result[i], i+1, result[i+1])
		}
	})

	t.Run("works with nanosecond precision", func(t *testing.T) {
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		input := []time.Time{
			baseTime,
			baseTime.Add(2 * time.Nanosecond),
		}
		item := baseTime.Add(1 * time.Nanosecond)

		result := addToSlice(ordTime, input, item)

		assert.Len(t, result, 3, "result should have 3 elements")
		assert.Equal(t, baseTime, result[0])
		assert.Equal(t, baseTime.Add(1*time.Nanosecond), result[1])
		assert.Equal(t, baseTime.Add(2*time.Nanosecond), result[2])
	})
}

func TestMakeClosedStateHistory(t *testing.T) {
	t.Run("creates a valid ClosedState", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)

		assert.NotNil(t, state, "MakeClosedStateHistory should return a non-nil ClosedState")
	})

	t.Run("initial state passes Check", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		result := state.Check(now)

		assert.True(t, option.IsSome(result), "initial state should pass Check (return Some, circuit stays closed)")
	})

	t.Run("Empty purges failure history", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(2)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add some errors
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))

		// Reset the state
		state = state.Empty()

		// Should pass check after reset
		result := state.Check(baseTime.Add(20 * time.Second))
		assert.True(t, option.IsSome(result), "state should pass Check after Empty")
	})

	t.Run("AddSuccess purges entire failure history", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))

		// Add success (should purge all history)
		state = state.AddSuccess(baseTime.Add(20 * time.Second))

		// Add another error (this is now the first error in history)
		state = state.AddError(baseTime.Add(30 * time.Second))

		// Should still pass check (only 1 error in history, threshold is 3)
		result := state.Check(baseTime.Add(30 * time.Second))
		assert.True(t, option.IsSome(result), "AddSuccess should purge entire failure history")
	})

	t.Run("circuit opens when failures reach threshold within time window", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors within time window but not reaching threshold
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))

		// Should still pass before threshold
		result := state.Check(baseTime.Add(20 * time.Second))
		assert.True(t, option.IsSome(result), "should pass Check before threshold")

		// Add one more error to reach threshold
		state = state.AddError(baseTime.Add(30 * time.Second))

		// Should fail check at threshold
		result = state.Check(baseTime.Add(30 * time.Second))
		assert.True(t, option.IsNone(result), "should fail Check when reaching threshold")
	})

	t.Run("old failures outside time window are automatically removed", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors that will be outside the time window
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))

		// Add error after time window has passed (this should remove old errors)
		state = state.AddError(baseTime.Add(2 * time.Minute))

		// Should pass check (only 1 error in window, old ones removed)
		result := state.Check(baseTime.Add(2 * time.Minute))
		assert.True(t, option.IsSome(result), "old failures should be removed from history")
	})

	t.Run("failures within time window are retained", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors within time window
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(30 * time.Second))
		state = state.AddError(baseTime.Add(50 * time.Second))

		// All errors are within 1 minute window, should fail check
		result := state.Check(baseTime.Add(50 * time.Second))
		assert.True(t, option.IsNone(result), "failures within time window should be retained")
	})

	t.Run("sliding window behavior - errors slide out over time", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add 3 errors to reach threshold
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))
		state = state.AddError(baseTime.Add(20 * time.Second))

		// Circuit should be open
		result := state.Check(baseTime.Add(20 * time.Second))
		assert.True(t, option.IsNone(result), "circuit should be open with 3 failures")

		// Add error after first failure has expired (> 1 minute from first error)
		// This should remove the first error, leaving only 3 in window
		state = state.AddError(baseTime.Add(70 * time.Second))

		// Should still fail check (3 errors in window after pruning)
		result = state.Check(baseTime.Add(70 * time.Second))
		assert.True(t, option.IsNone(result), "circuit should remain open with 3 failures in window")
	})

	t.Run("zero maxFailures means circuit is always open", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(0)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Initial state should already fail (0 >= 0)
		result := state.Check(baseTime)
		assert.True(t, option.IsNone(result), "initial state should fail Check with maxFailures=0")

		// Add one error
		state = state.AddError(baseTime)

		// Should still fail
		result = state.Check(baseTime)
		assert.True(t, option.IsNone(result), "should fail Check after error with maxFailures=0")
	})

	t.Run("success purges history even with failures in time window", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors within time window
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))

		// Add success (purges all history)
		state = state.AddSuccess(baseTime.Add(20 * time.Second))

		// Add more errors
		state = state.AddError(baseTime.Add(30 * time.Second))
		state = state.AddError(baseTime.Add(40 * time.Second))

		// Should still pass (only 2 errors after purge)
		result := state.Check(baseTime.Add(40 * time.Second))
		assert.True(t, option.IsSome(result), "success should purge all history")

		// Add one more to reach threshold
		state = state.AddError(baseTime.Add(50 * time.Second))

		// Should fail at threshold
		result = state.Check(baseTime.Add(50 * time.Second))
		assert.True(t, option.IsNone(result), "should fail after reaching threshold")
	})

	t.Run("multiple successes keep history empty", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(2)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add error
		state = state.AddError(baseTime)

		// Multiple successes
		state = state.AddSuccess(baseTime.Add(10 * time.Second))
		state = state.AddSuccess(baseTime.Add(20 * time.Second))
		state = state.AddSuccess(baseTime.Add(30 * time.Second))

		// Should still pass
		result := state.Check(baseTime.Add(30 * time.Second))
		assert.True(t, option.IsSome(result), "multiple AddSuccess should keep history empty")

		// Add errors to reach threshold
		state = state.AddError(baseTime.Add(40 * time.Second))
		state = state.AddError(baseTime.Add(50 * time.Second))

		// Should fail
		result = state.Check(baseTime.Add(50 * time.Second))
		assert.True(t, option.IsNone(result), "should fail after reaching threshold")
	})

	t.Run("state is immutable - original unchanged after AddError", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(2)
		originalState := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Create new state by adding error
		newState := originalState.AddError(baseTime)

		// Original should still pass check
		result := originalState.Check(baseTime)
		assert.True(t, option.IsSome(result), "original state should be unchanged")

		// New state should reach threshold after another error
		newState = newState.AddError(baseTime.Add(10 * time.Second))

		result = newState.Check(baseTime.Add(10 * time.Second))
		assert.True(t, option.IsNone(result), "new state should fail after reaching threshold")

		// Original should still pass
		result = originalState.Check(baseTime)
		assert.True(t, option.IsSome(result), "original state should still be unchanged")
	})

	t.Run("state is immutable - original unchanged after Empty", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(2)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors to original
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))
		stateWithErrors := state

		// Create new state by calling Empty
		emptyState := stateWithErrors.Empty()

		// Original with errors should fail check
		result := stateWithErrors.Check(baseTime.Add(10 * time.Second))
		assert.True(t, option.IsNone(result), "state with errors should fail after reaching threshold")

		// Empty state should pass
		result = emptyState.Check(baseTime.Add(10 * time.Second))
		assert.True(t, option.IsSome(result), "empty state should pass Check")
	})

	t.Run("exact time window boundary behavior", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add error at baseTime
		state = state.AddError(baseTime)

		// Add error exactly at time window boundary
		state = state.AddError(baseTime.Add(1 * time.Minute))

		// The first error should be removed (it's now outside the window)
		// Only 1 error should remain
		result := state.Check(baseTime.Add(1 * time.Minute))
		assert.True(t, option.IsSome(result), "error at exact window boundary should remove older errors")
	})

	t.Run("multiple errors at same timestamp", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add multiple errors at same time
		state = state.AddError(baseTime)
		state = state.AddError(baseTime)
		state = state.AddError(baseTime)

		// Should fail check (3 errors at same time)
		result := state.Check(baseTime)
		assert.True(t, option.IsNone(result), "multiple errors at same timestamp should count separately")
	})

	t.Run("errors added out of chronological order are sorted", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(4)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors out of order
		state = state.AddError(baseTime.Add(30 * time.Second))
		state = state.AddError(baseTime.Add(5 * time.Second))
		state = state.AddError(baseTime.Add(50 * time.Second))

		// Add error that should trigger pruning
		state = state.AddError(baseTime.Add(70 * time.Second))

		// The error at 5s should be removed (> 1 minute from 70s: 70-5=65 > 60)
		// Should have 3 errors remaining (30s, 50s, 70s)
		result := state.Check(baseTime.Add(70 * time.Second))
		assert.True(t, option.IsSome(result), "errors should be sorted and pruned correctly")
	})

	t.Run("large time window with many failures", func(t *testing.T) {
		timeWindow := 24 * time.Hour
		maxFailures := uint(100)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add many failures within the window
		for i := 0; i < 99; i++ {
			state = state.AddError(baseTime.Add(time.Duration(i) * time.Minute))
		}

		// Should still pass (99 < 100)
		result := state.Check(baseTime.Add(99 * time.Minute))
		assert.True(t, option.IsSome(result), "should pass with 99 failures when threshold is 100")

		// Add one more to reach threshold
		state = state.AddError(baseTime.Add(100 * time.Minute))

		// Should fail
		result = state.Check(baseTime.Add(100 * time.Minute))
		assert.True(t, option.IsNone(result), "should fail at threshold with large window")
	})

	t.Run("very short time window", func(t *testing.T) {
		timeWindow := 100 * time.Millisecond
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors within short window
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(50 * time.Millisecond))
		state = state.AddError(baseTime.Add(90 * time.Millisecond))

		// Should fail (3 errors within 100ms)
		result := state.Check(baseTime.Add(90 * time.Millisecond))
		assert.True(t, option.IsNone(result), "should fail with errors in short time window")

		// Add error after window expires
		state = state.AddError(baseTime.Add(200 * time.Millisecond))

		// Should pass (old errors removed, only 1 in window)
		result = state.Check(baseTime.Add(200 * time.Millisecond))
		assert.True(t, option.IsSome(result), "should pass after short window expires")
	})

	t.Run("success prevents circuit from opening", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(3)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors close to threshold
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))

		// Add success before reaching threshold
		state = state.AddSuccess(baseTime.Add(20 * time.Second))

		// Add more errors
		state = state.AddError(baseTime.Add(30 * time.Second))
		state = state.AddError(baseTime.Add(40 * time.Second))

		// Should still pass (only 2 errors after success purge)
		result := state.Check(baseTime.Add(40 * time.Second))
		assert.True(t, option.IsSome(result), "circuit should stay closed after success purge")
	})

	t.Run("Empty can be called multiple times", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(2)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add errors
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(10 * time.Second))
		state = state.AddError(baseTime.Add(20 * time.Second))

		// Reset multiple times
		state = state.Empty()
		state = state.Empty()
		state = state.Empty()

		// Should still pass
		result := state.Check(baseTime.Add(30 * time.Second))
		assert.True(t, option.IsSome(result), "state should pass Check after multiple Empty calls")
	})

	t.Run("gradual failure accumulation within window", func(t *testing.T) {
		timeWindow := 1 * time.Minute
		maxFailures := uint(5)
		state := MakeClosedStateHistory(timeWindow, maxFailures)
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		// Add failures gradually
		state = state.AddError(baseTime)
		state = state.AddError(baseTime.Add(15 * time.Second))
		state = state.AddError(baseTime.Add(30 * time.Second))
		state = state.AddError(baseTime.Add(45 * time.Second))

		// Should still pass (4 < 5)
		result := state.Check(baseTime.Add(45 * time.Second))
		assert.True(t, option.IsSome(result), "should pass before threshold")

		// Add one more within window
		state = state.AddError(baseTime.Add(55 * time.Second))

		// Should fail (5 >= 5)
		result = state.Check(baseTime.Add(55 * time.Second))
		assert.True(t, option.IsNone(result), "should fail at threshold")
	})
}
