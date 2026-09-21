package circuitbreaker

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
)

// describe renders the outcome of a [ClosedState.Check] in a readable way: the check
// yields None when the failure threshold is reached and the circuit has to open.
var describe = option.Fold(
	lazy.Of("open"),
	F.Constant1[ClosedState]("closed"),
)

// ExampleMakeClosedStateCounter shows the counter based failure tracking. The circuit
// opens on the n-th *consecutive* failure; a single success resets the counter.
func ExampleMakeClosedStateCounter() {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	initial := MakeClosedStateCounter(3)
	fmt.Println("initial       ", describe(initial.Check(now)))

	twoFailures := F.Pipe2(initial, addErrorAt(now), addErrorAt(now))
	fmt.Println("2 failures    ", describe(twoFailures.Check(now)))

	// a success wipes out the failures that came before it
	afterSuccess := F.Pipe2(twoFailures, addSuccessAt(now), addErrorAt(now))
	fmt.Println("success, 1 err", describe(afterSuccess.Check(now)))

	threeFailures := F.Pipe2(afterSuccess, addErrorAt(now), addErrorAt(now))
	fmt.Println("3 failures    ", describe(threeFailures.Check(now)))

	// Output:
	// initial        closed
	// 2 failures     closed
	// success, 1 err closed
	// 3 failures     open
}

// ExampleMakeClosedStateHistory shows the sliding window based failure tracking. Only
// failures that happened inside the window count towards the threshold, so a circuit
// recovers by the mere passage of time.
func ExampleMakeClosedStateHistory() {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// open the circuit on 3 failures within one minute
	state := F.Pipe3(
		MakeClosedStateHistory(time.Minute, 3),
		addErrorAt(base),
		addErrorAt(base.Add(10*time.Second)),
		addErrorAt(base.Add(20*time.Second)),
	)
	fmt.Println("3 failures in 20s ", describe(state.Check(base.Add(20*time.Second))))

	// two minutes later all three failures have left the window
	fmt.Println("2 minutes later   ", describe(state.Check(base.Add(2*time.Minute))))

	// Output:
	// 3 failures in 20s  open
	// 2 minutes later    closed
}

// ExampleMakeClosedIORef shows how the mutable breaker state is created and inspected.
func ExampleMakeClosedIORef() {
	ref := io.Run(MakeClosedIORef(MakeClosedStateCounter(3)))

	state := io.Run(ioref.Read(ref))
	fmt.Println("open:  ", IsOpen(state))
	fmt.Println("closed:", IsClosed(state))

	// Output:
	// open:   false
	// closed: true
}

// ExampleMakeCircuitBreakerErrorWithName shows the error that a breaker returns for the
// requests it blocks. The name makes it possible to tell several breakers apart.
func ExampleMakeCircuitBreakerErrorWithName() {
	makeError := MakeCircuitBreakerErrorWithName("Database")

	resetAt := time.Date(2024, 1, 1, 12, 0, 30, 0, time.UTC)
	err := makeError(resetAt)

	fmt.Println(err)

	var cbErr *CircuitBreakerError
	if errors.As(err, &cbErr) {
		fmt.Println("resets at:", cbErr.ResetAt)
	}

	// Output:
	// circuit breaker is open [Database], will close at 2024-01-01 12:00:30 +0000 UTC
	// resets at: 2024-01-01 12:00:30 +0000 UTC
}

// ExampleAnyError shows the most permissive error filter: every error counts towards the
// failure threshold of the breaker.
func ExampleAnyError() {
	fmt.Println(option.IsSome(AnyError(errors.New("boom"))))
	fmt.Println(option.IsSome(AnyError(nil)))

	// Output:
	// true
	// false
}

// ExampleInfrastructureError shows the error filter that only lets infrastructure
// failures open the circuit. Application level errors leave the circuit closed, because
// retrying them later would not help.
func ExampleInfrastructureError() {
	// the remote end is not reachable, the service itself is in trouble
	unreachable := &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
	fmt.Println("connection refused:", option.IsSome(InfrastructureError(unreachable)))

	// a validation failure says nothing about the health of the service
	invalid := errors.New("invalid request payload")
	fmt.Println("validation error:  ", option.IsSome(InfrastructureError(invalid)))

	// Output:
	// connection refused: true
	// validation error:   false
}

// ExampleMakeMetricsFromLogger shows the logging metrics sink. Every event is an [IO],
// so nothing is written until the operation is run.
func ExampleMakeMetricsFromLogger() {
	logger := log.New(os.Stdout, "[CB] ", 0)
	metrics := MakeMetricsFromLogger("UserService", logger)

	ct := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// describing the events has no effect yet
	events := []IO[Void]{
		metrics.Accept(ct),
		metrics.Open(ct),
		metrics.Reject(ct),
		metrics.Canary(ct),
		metrics.Close(ct),
	}

	for _, event := range events {
		io.Run(event)
	}

	// Output:
	// [CB] Accept: UserService, 2024-01-01 12:00:00 +0000 UTC
	// [CB] Open: UserService, 2024-01-01 12:00:00 +0000 UTC
	// [CB] Reject: UserService, 2024-01-01 12:00:00 +0000 UTC
	// [CB] Canary: UserService, 2024-01-01 12:00:00 +0000 UTC
	// [CB] Close: UserService, 2024-01-01 12:00:00 +0000 UTC
}

// ExampleMakeVoidMetrics shows the metrics sink that discards every event. Use it when
// the breaker should not produce any observability overhead.
func ExampleMakeVoidMetrics() {
	metrics := MakeVoidMetrics()

	ct := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	io.Run(metrics.Open(ct))

	fmt.Println("nothing was reported")

	// Output:
	// nothing was reported
}
