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

package circuitbreaker

import (
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// testTime is the fixed point in time all tests in this file are anchored at
var testTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

// openStateOf extracts the openState of an open circuit, failing the test when the
// circuit is closed
func openStateOf(t *testing.T, state BreakerState) openState {
	t.Helper()
	assert.True(t, IsOpen(state), "expected an open circuit")
	return either.Fold(F.Identity[openState], func(ClosedState) openState { return openState{} })(state)
}

// closedStateOf extracts the ClosedState of a closed circuit, failing the test when the
// circuit is open
func closedStateOf(t *testing.T, state BreakerState) ClosedState {
	t.Helper()
	assert.True(t, IsClosed(state), "expected a closed circuit")
	return either.GetOrElse(func(openState) ClosedState { return nil })(state)
}

// TestFanout tests the fanout combinator that applies two readers to the same
// environment and pairs up their results
func TestFanout(t *testing.T) {
	t.Run("pairs the results of both readers", func(t *testing.T) {
		length := func(s string) int { return len(s) }
		exclaim := func(s string) string { return s + "!" }

		result := fanout(length, exclaim)("abc")

		assert.Equal(t, 3, pair.Head(result))
		assert.Equal(t, "abc!", pair.Tail(result))
	})

	t.Run("passes the same environment to both readers", func(t *testing.T) {
		var seen []string
		record := func(tag string) Reader[string, string] {
			return func(env string) string {
				seen = append(seen, tag+":"+env)
				return env
			}
		}

		fanout(record("left"), record("right"))("env")

		assert.Equal(t, []string{"left:env", "right:env"}, seen)
	})

	t.Run("works with the state types of this package", func(t *testing.T) {
		state := openState{openedAt: testTime, resetAt: testTime.Add(time.Minute)}

		result := fanout(openedAtLens.Get, resetAtLens.Get)(state)

		assert.Equal(t, testTime, pair.Head(result))
		assert.Equal(t, testTime.Add(time.Minute), pair.Tail(result))
	})
}

// TestResetTimeExceeded tests the predicate that decides whether the reset time of an
// open circuit has passed
func TestResetTimeExceeded(t *testing.T) {
	stateWithResetAt := func(resetAt time.Time) openState {
		return openState{
			openedAt:    testTime.Add(-time.Minute),
			resetAt:     resetAt,
			retryStatus: retry.DefaultRetryStatus,
		}
	}

	t.Run("holds when the reset time is in the past", func(t *testing.T) {
		assert.True(t, resetTimeExceeded(testTime)(stateWithResetAt(testTime.Add(-time.Second))))
	})

	t.Run("fails when the reset time is in the future", func(t *testing.T) {
		assert.False(t, resetTimeExceeded(testTime)(stateWithResetAt(testTime.Add(time.Second))))
	})

	t.Run("fails when the reset time is exactly the current time", func(t *testing.T) {
		assert.False(t, resetTimeExceeded(testTime)(stateWithResetAt(testTime)),
			"the boundary is exclusive because the predicate is based on time.Time.After")
	})

	t.Run("ignores a pending canary", func(t *testing.T) {
		state := stateWithResetAt(testTime.Add(-time.Second))
		state.canaryRequest = true

		assert.True(t, resetTimeExceeded(testTime)(state),
			"the canary flag is the concern of canaryAllowed, not of resetTimeExceeded")
	})
}

// TestOpenCircuitFromStatus tests building the openState that belongs to a retry status
func TestOpenCircuitFromStatus(t *testing.T) {
	t.Run("delays the reset time by the previous delay of the status", func(t *testing.T) {
		status := retry.PreviousDelayLens.Set(option.Of(2 * time.Second))(retry.DefaultRetryStatus)

		result := openCircuitFromStatus(status)(testTime)

		assert.Equal(t, testTime, result.openedAt)
		assert.Equal(t, testTime.Add(2*time.Second), result.resetAt)
	})

	t.Run("resets immediately when the status carries no delay", func(t *testing.T) {
		status := retry.PreviousDelayLens.Set(option.None[time.Duration]())(retry.DefaultRetryStatus)

		result := openCircuitFromStatus(status)(testTime)

		assert.Equal(t, testTime, result.resetAt)
	})

	t.Run("keeps the retry status and clears the canary flag", func(t *testing.T) {
		status := retry.RetryStatus{
			IterNumber:      3,
			CumulativeDelay: 7 * time.Second,
			PreviousDelay:   option.Of(time.Second),
		}

		result := openCircuitFromStatus(status)(testTime)

		assert.Equal(t, status, result.retryStatus)
		assert.False(t, result.canaryRequest, "a freshly opened circuit has no canary in flight")
	})

	t.Run("reads the opening time from the environment", func(t *testing.T) {
		later := testTime.Add(time.Hour)

		result := openCircuitFromStatus(retry.DefaultRetryStatus)(later)

		assert.Equal(t, later, result.openedAt)
	})
}

// TestRecordFailure tests the pure core of handleFailureOnClosed that decides whether a
// recorded failure keeps the circuit closed or opens it
func TestRecordFailure(t *testing.T) {
	open := openState{
		openedAt:    testTime,
		resetAt:     testTime.Add(time.Minute),
		retryStatus: retry.DefaultRetryStatus,
	}

	addError := func(cs ClosedState) ClosedState { return cs.AddError(testTime) }
	staysClosed := option.Of[ClosedState]
	opens := func(ClosedState) Option[ClosedState] { return option.None[ClosedState]() }

	t.Run("keeps the circuit closed while the check succeeds", func(t *testing.T) {
		initial := createClosedCircuit(MakeClosedStateCounter(3))

		result := recordFailure(addError, staysClosed, open)(initial)

		assert.True(t, IsClosed(result))
	})

	t.Run("records the failure in the resulting closed state", func(t *testing.T) {
		// a counter that opens on the second failure: after recordFailure applied the
		// first one, a single further failure must reach the threshold
		initial := createClosedCircuit(MakeClosedStateCounter(2))

		recorded := closedStateOf(t, recordFailure(addError, staysClosed, open)(initial))

		assert.True(t, option.IsNone(recorded.AddError(testTime).Check(testTime)),
			"the failure recorded by addError must be part of the resulting state")
	})

	t.Run("opens the circuit with the given open state when the check fails", func(t *testing.T) {
		initial := createClosedCircuit(MakeClosedStateCounter(1))

		result := recordFailure(addError, opens, open)(initial)

		assert.Equal(t, open, openStateOf(t, result))
	})

	t.Run("leaves an already open circuit untouched", func(t *testing.T) {
		other := openState{openedAt: testTime.Add(-time.Hour), resetAt: testTime}
		called := false
		trackingAddError := func(cs ClosedState) ClosedState {
			called = true
			return addError(cs)
		}

		result := recordFailure(trackingAddError, opens, open)(createOpenCircuit(other))

		assert.Equal(t, other, openStateOf(t, result), "an open circuit keeps its own open state")
		assert.False(t, called, "no failure is recorded while the circuit is open")
	})
}

// TestHandleErrorOnClosed tests that checkError decides whether an error counts as a
// failure for the circuit breaker or is treated like a success
func TestHandleErrorOnClosed(t *testing.T) {
	// tag returns a handler that records that it was selected, together with the time
	// it was resolved for
	tag := func(name string, seen *[]string) Reader[time.Time, Endomorphism[BreakerState]] {
		return func(ct time.Time) Endomorphism[BreakerState] {
			*seen = append(*seen, name+"@"+ct.Format(time.RFC3339))
			return F.Identity[BreakerState]
		}
	}

	countsAsFailure := option.Of[error]
	ignored := func(error) Option[error] { return option.None[error]() }

	t.Run("routes a relevant error to the failure handler", func(t *testing.T) {
		var seen []string

		handleErrorOnClosed(countsAsFailure, tag("success", &seen), tag("failure", &seen))(assert.AnError)(testTime)

		assert.Equal(t, []string{"failure@" + testTime.Format(time.RFC3339)}, seen)
	})

	t.Run("routes an irrelevant error to the success handler", func(t *testing.T) {
		var seen []string

		handleErrorOnClosed(ignored, tag("success", &seen), tag("failure", &seen))(assert.AnError)(testTime)

		assert.Equal(t, []string{"success@" + testTime.Format(time.RFC3339)}, seen)
	})

	t.Run("hands the error to checkError", func(t *testing.T) {
		var seen []string
		var inspected []error

		checkError := func(e error) Option[error] {
			inspected = append(inspected, e)
			return option.None[error]()
		}

		handleErrorOnClosed(checkError, tag("success", &seen), tag("failure", &seen))(assert.AnError)(testTime)

		assert.Equal(t, []error{assert.AnError}, inspected)
	})

	t.Run("applies the endomorphism of the selected branch", func(t *testing.T) {
		opened := createOpenCircuit(openState{openedAt: testTime, resetAt: testTime.Add(time.Minute)})
		toOpen := reader.Of[time.Time](F.Constant1[BreakerState](opened))
		unchanged := reader.Of[time.Time](F.Identity[BreakerState])

		handler := handleErrorOnClosed(countsAsFailure, unchanged, toOpen)
		result := handler(assert.AnError)(testTime)(createClosedCircuit(MakeClosedStateCounter(3)))

		assert.True(t, IsOpen(result), "the failure branch replaced the state")
	})
}

// TestReportTransition tests that only an actual change of the circuit state is
// reported through the matching metric
func TestReportTransition(t *testing.T) {
	// recorder builds a ReaderIO that notes the time it was invoked with
	recorder := func(seen *[]time.Time) ReaderIO[time.Time, Void] {
		return func(ct time.Time) IO[Void] {
			return func() Void {
				*seen = append(*seen, ct)
				return function.VOID
			}
		}
	}

	openBreaker := createOpenCircuit(openState{openedAt: testTime, resetAt: testTime.Add(time.Minute)})
	otherOpenBreaker := createOpenCircuit(openState{openedAt: testTime, resetAt: testTime.Add(time.Hour)})
	closedBreaker := createClosedCircuit(MakeClosedStateCounter(3))

	// report runs the transition from previous to current and returns what was recorded
	report := func(previous, current BreakerState) (opened, closed []time.Time) {
		change := readerio.Of[time.Time](pair.MakePair(previous, current))
		reportTransition(recorder(&opened), recorder(&closed))(change)(testTime)()
		return
	}

	t.Run("reports closed to open through onOpened", func(t *testing.T) {
		opened, closed := report(closedBreaker, openBreaker)

		assert.Equal(t, []time.Time{testTime}, opened)
		assert.Empty(t, closed)
	})

	t.Run("reports open to closed through onClosed", func(t *testing.T) {
		opened, closed := report(openBreaker, closedBreaker)

		assert.Equal(t, []time.Time{testTime}, closed)
		assert.Empty(t, opened)
	})

	t.Run("reports nothing when the circuit stayed open", func(t *testing.T) {
		opened, closed := report(openBreaker, otherOpenBreaker)

		assert.Empty(t, opened, "a circuit that was open already must not be reported as opening")
		assert.Empty(t, closed)
	})

	t.Run("reports nothing when the circuit stayed closed", func(t *testing.T) {
		opened, closed := report(closedBreaker, createClosedCircuit(MakeClosedStateCounter(5)))

		assert.Empty(t, opened)
		assert.Empty(t, closed)
	})

	t.Run("does not report before the effect is executed", func(t *testing.T) {
		var opened, closed []time.Time
		change := readerio.Of[time.Time](pair.MakePair(closedBreaker, openBreaker))

		effect := reportTransition(recorder(&opened), recorder(&closed))(change)(testTime)

		assert.Empty(t, opened, "building the effect must not emit a metric")

		effect()

		assert.Len(t, opened, 1)
	})

	t.Run("emits nothing observable for the noReport branch", func(t *testing.T) {
		var closed []time.Time
		change := readerio.Of[time.Time](pair.MakePair(closedBreaker, openBreaker))

		reportTransition(noReport, recorder(&closed))(change)(testTime)()

		assert.Empty(t, closed, "the opening branch was taken, so onClosed stays untouched")
	})
}

// TestTrackPrevious tests the update shape that keeps the state before an update
// alongside the state after it
func TestTrackPrevious(t *testing.T) {
	closedBreaker := createClosedCircuit(MakeClosedStateCounter(3))
	openBreaker := createOpenCircuit(openState{openedAt: testTime, resetAt: testTime.Add(time.Minute)})

	t.Run("stores the new state and reports both states", func(t *testing.T) {
		result := trackPrevious(F.Constant1[BreakerState](openBreaker))(closedBreaker)

		assert.True(t, IsOpen(pair.Head(result)), "the head is the state to store")
		assert.True(t, IsClosed(pair.Head(pair.Tail(result))), "the previous state comes first")
		assert.True(t, IsOpen(pair.Tail(pair.Tail(result))), "the new state comes second")
	})

	t.Run("reports identical states for an update that changes nothing", func(t *testing.T) {
		result := trackPrevious(F.Identity[BreakerState])(closedBreaker)

		assert.True(t, IsClosed(pair.Head(pair.Tail(result))))
		assert.True(t, IsClosed(pair.Tail(pair.Tail(result))))
	})
}

// TestArmCanary tests that starting a canary marks the circuit half-open and pushes the
// reset time to the canary deadline
func TestArmCanary(t *testing.T) {
	t.Run("sets the canary flag", func(t *testing.T) {
		state := openState{openedAt: testTime, resetAt: testTime.Add(-time.Second), retryStatus: retry.DefaultRetryStatus}

		assert.True(t, armCanary(testTime)(state).canaryRequest)
	})

	t.Run("pushes the reset time to the canary deadline", func(t *testing.T) {
		state := openState{
			openedAt:    testTime.Add(-time.Minute),
			resetAt:     testTime.Add(-time.Second),
			retryStatus: retry.PreviousDelayLens.Set(option.Of(30 * time.Second))(retry.DefaultRetryStatus),
		}

		result := armCanary(testTime)(state)

		assert.Equal(t, testTime.Add(30*time.Second), result.resetAt,
			"the deadline is the current time plus the delay of the retry status")
		assert.False(t, canaryAllowed(testTime)(result), "no further canary while the deadline has not passed")
		assert.True(t, canaryAllowed(testTime.Add(31*time.Second))(result),
			"a canary that never reported back is replaced after the deadline")
	})

	t.Run("keeps the remaining fields", func(t *testing.T) {
		status := retry.PreviousDelayLens.Set(option.Of(time.Second))(retry.DefaultRetryStatus)
		state := openState{openedAt: testTime.Add(-time.Hour), resetAt: testTime, retryStatus: status}

		result := armCanary(testTime)(state)

		assert.Equal(t, testTime.Add(-time.Hour), result.openedAt)
		assert.Equal(t, status, result.retryStatus)
	})

	t.Run("does not modify its input", func(t *testing.T) {
		state := openState{openedAt: testTime, resetAt: testTime.Add(-time.Second), retryStatus: retry.DefaultRetryStatus}

		armCanary(testTime)(state)

		assert.False(t, state.canaryRequest)
		assert.Equal(t, testTime.Add(-time.Second), state.resetAt)
	})
}

// TestApplyAndReport tests that a state update is committed to the IORef and that the
// resulting state is handed to the report operator
func TestApplyAndReport(t *testing.T) {
	opened := openState{openedAt: testTime, resetAt: testTime.Add(time.Minute)}
	toOpen := reader.Of[time.Time](F.Constant1[BreakerState](createOpenCircuit(opened)))

	// capturingReport runs the inner effect and records both the state it produced and
	// the time it was resolved for
	capturingReport := func(seen *[]BreakerState, times *[]time.Time) readerio.Operator[time.Time, Pair[BreakerState, BreakerState], Void] {
		return func(rio ReaderIO[time.Time, Pair[BreakerState, BreakerState]]) ReaderIO[time.Time, Void] {
			return func(ct time.Time) IO[Void] {
				return func() Void {
					*times = append(*times, ct)
					*seen = append(*seen, pair.Tail(rio(ct)()))
					return function.VOID
				}
			}
		}
	}

	t.Run("commits the update to the ioref and reports the new state", func(t *testing.T) {
		ref := io.Run(ioref.MakeIORef(createClosedCircuit(MakeClosedStateCounter(3))))
		var seen []BreakerState
		var times []time.Time

		applyAndReport(io.Of(testTime), capturingReport(&seen, &times), modifyPrevV(ref))(toOpen)()

		assert.Equal(t, opened, openStateOf(t, io.Run(ioref.Read(ref))),
			"the new state must be visible in the ioref")
		assert.Len(t, seen, 1)
		assert.Equal(t, opened, openStateOf(t, seen[0]), "the reported state is the updated one")
	})

	t.Run("resolves update and report with the current time", func(t *testing.T) {
		ref := io.Run(ioref.MakeIORef(createClosedCircuit(MakeClosedStateCounter(3))))
		var seen []BreakerState
		var times []time.Time
		var updateTimes []time.Time

		update := func(ct time.Time) Endomorphism[BreakerState] {
			updateTimes = append(updateTimes, ct)
			return F.Identity[BreakerState]
		}

		applyAndReport(io.Of(testTime), capturingReport(&seen, &times), modifyPrevV(ref))(update)()

		assert.Equal(t, []time.Time{testTime}, updateTimes, "the update sees the current time")
		assert.Equal(t, []time.Time{testTime}, times, "and so does the report")
	})

	t.Run("does not touch the ioref before the effect is executed", func(t *testing.T) {
		ref := io.Run(ioref.MakeIORef(createClosedCircuit(MakeClosedStateCounter(3))))
		var seen []BreakerState
		var times []time.Time

		effect := applyAndReport(io.Of(testTime), capturingReport(&seen, &times), modifyPrevV(ref))(toOpen)

		assert.True(t, IsClosed(io.Run(ioref.Read(ref))), "building the effect must not modify the state")
		assert.Empty(t, seen)

		effect()

		assert.True(t, IsOpen(io.Run(ioref.Read(ref))))
	})

	t.Run("applies updates cumulatively across executions", func(t *testing.T) {
		ref := io.Run(ioref.MakeIORef(createClosedCircuit(MakeClosedStateCounter(3))))
		var seen []BreakerState
		var times []time.Time

		addError := func(ct time.Time) Endomorphism[BreakerState] {
			return either.Map[openState](func(cs ClosedState) ClosedState { return cs.AddError(ct) })
		}

		effect := applyAndReport(io.Of(testTime), capturingReport(&seen, &times), modifyPrevV(ref))(addError)
		effect()
		effect()

		assert.Len(t, seen, 2, "each execution reports the state it produced")
		assert.True(t, option.IsNone(closedStateOf(t, io.Run(ioref.Read(ref))).AddError(testTime).Check(testTime)),
			"two failures were committed, so a third one reaches the threshold of 3")
	})

	t.Run("reads the current time on each execution", func(t *testing.T) {
		ref := io.Run(ioref.MakeIORef(createClosedCircuit(MakeClosedStateCounter(3))))
		var seen []BreakerState
		var times []time.Time

		ticks := 0
		clock := func() time.Time {
			ticks++
			return testTime.Add(time.Duration(ticks) * time.Second)
		}

		effect := applyAndReport(clock, capturingReport(&seen, &times), modifyPrevV(ref))(toOpen)
		effect()
		effect()

		assert.Equal(t, []time.Time{testTime.Add(time.Second), testTime.Add(2 * time.Second)}, times,
			"the clock is read per execution, not once when the effect is built")
	})
}
