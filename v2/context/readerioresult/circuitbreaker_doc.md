# Circuit Breaker Documentation

## Overview

The `circuitbreaker.go` file provides a circuit breaker implementation for the `readerioresult` package. A circuit breaker is a design pattern used to detect failures and prevent cascading failures in distributed systems by temporarily blocking operations that are likely to fail.

## Package

```go
package readerioresult
```

This is part of the `context/readerioresult` package, which provides functional programming abstractions for operations that:
- Depend on a `context.Context` (Reader aspect)
- Perform side effects (IO aspect)  
- Can fail with an `error` (Result/Either aspect)

## Type Definitions

### ClosedState

```go
type ClosedState = circuitbreaker.ClosedState
```

A type alias for the circuit breaker's closed state. When the circuit is closed, requests are allowed to pass through normally. The closed state tracks success and failure counts to determine when to open the circuit.

### Env[T any]

```go
type Env[T any] = Pair[IORef[circuitbreaker.BreakerState], ReaderIOResult[T]]
```

The environment type for the circuit breaker state machine. It contains:
- `IORef[circuitbreaker.BreakerState]`: A mutable reference to the current breaker state
- `ReaderIOResult[T]`: The computation to be protected by the circuit breaker

### CircuitBreaker[T any]

```go
type CircuitBreaker[T any] = State[Env[T], ReaderIOResult[T]]
```

The main circuit breaker type. It's a state monad that:
- Takes an environment containing the breaker state and the protected computation
- Returns a new environment and a wrapped computation that respects the circuit breaker logic

## Functions

### MakeCircuitBreaker

```go
func MakeCircuitBreaker[T any](
    currentTime IO[time.Time],
    closedState ClosedState,
    checkError option.Kleisli[error, error],
    policy retry.RetryPolicy,
    logger io.Kleisli[string, string],
) CircuitBreaker[T]
```

Creates a new circuit breaker with the specified configuration.

#### Parameters

- **currentTime** `IO[time.Time]`: A function that returns the current time. This can be a virtual timer for testing purposes, allowing you to control time progression in tests.

- **closedState** `ClosedState`: The initial closed state configuration. This defines:
  - Maximum number of failures before opening the circuit
  - Time window for counting failures
  - Other closed state parameters

- **checkError** `option.Kleisli[error, error]`: A function that determines whether an error should be counted as a failure. Returns:
  - `Some(error)`: The error should be counted as a failure
  - `None`: The error should be ignored (not counted as a failure)
  
  This allows you to distinguish between transient errors (that should trigger circuit breaking) and permanent errors (that shouldn't).

- **policy** `retry.RetryPolicy`: The retry policy that determines:
  - How long to wait before attempting to close the circuit (reset time)
  - Exponential backoff or other delay strategies
  - Maximum number of retry attempts

- **logger** `io.Kleisli[string, string]`: A logging function for circuit breaker events. Receives log messages and performs side effects (like writing to a log file or console).

#### Returns

A `CircuitBreaker[T]` that wraps computations with circuit breaker logic.

#### Circuit Breaker States

The circuit breaker operates in three states:

1. **Closed**: Normal operation. Requests pass through. Failures are counted.
   - If failure threshold is exceeded, transitions to Open state

2. **Open**: Circuit is broken. Requests fail immediately without executing.
   - After reset time expires, transitions to Half-Open state

3. **Half-Open** (Canary): Testing if the service has recovered.
   - Allows a single test request (canary request)
   - If canary succeeds, transitions to Closed state
   - If canary fails, transitions back to Open state with extended reset time

#### Implementation Details

The function delegates to the generic `circuitbreaker.MakeCircuitBreaker` function, providing the necessary type-specific operations:

- **Left**: Creates a failed computation from an error
- **ChainFirstIOK**: Chains an IO operation that runs for side effects on success
- **ChainFirstLeftIOK**: Chains an IO operation that runs for side effects on failure
- **FromIO**: Lifts an IO computation into ReaderIOResult
- **Flap**: Applies a computation to a function
- **Flatten**: Flattens nested ReaderIOResult structures

These operations allow the generic circuit breaker to work with the `ReaderIOResult` monad.

## Usage Example

```go
import (
    "context"
    "fmt"
    "time"
    
    "github.com/IBM/fp-go/v2/circuitbreaker"
    "github.com/IBM/fp-go/v2/context/readerioresult"
    "github.com/IBM/fp-go/v2/io"
    "github.com/IBM/fp-go/v2/ioref"
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/retry"
)

// Create a circuit breaker configuration
func createCircuitBreaker() readerioresult.CircuitBreaker[string] {
    // Use real time
    currentTime := func() time.Time { return time.Now() }
    
    // Configure closed state: open after 5 failures in 10 seconds
    closedState := circuitbreaker.MakeClosedState(5, 10*time.Second)
    
    // Check all errors (count all as failures)
    checkError := func(err error) option.Option[error] {
        return option.Some(err)
    }
    
    // Retry policy: exponential backoff with max 5 retries
    policy := retry.Monoid.Concat(
        retry.LimitRetries(5),
        retry.ExponentialBackoff(100*time.Millisecond),
    )
    
    // Simple logger
    logger := func(msg string) io.IO[string] {
        return func() string {
            fmt.Println("Circuit Breaker:", msg)
            return msg
        }
    }
    
    return readerioresult.MakeCircuitBreaker[string](
        currentTime,
        closedState,
        checkError,
        policy,
        logger,
    )
}

// Use the circuit breaker
func main() {
    cb := createCircuitBreaker()
    
    // Create initial state
    stateRef := ioref.NewIORef(circuitbreaker.InitialState())
    
    // Your protected operation
    operation := func(ctx context.Context) readerioresult.IOResult[string] {
        return func() readerioresult.Result[string] {
            // Your actual operation here
            return result.Of("success")
        }
    }
    
    // Apply circuit breaker
    env := pair.MakePair(stateRef, operation)
    result := cb(env)
    
    // Execute the protected operation
    ctx := t.Context()
    protectedOp := pair.Tail(result)
    outcome := protectedOp(ctx)()
}
```

## Testing with Virtual Timer

For testing, you can provide a virtual timer instead of `time.Now()`:

```go
// Virtual timer for testing
type VirtualTimer struct {
    current time.Time
}

func (vt *VirtualTimer) Now() time.Time {
    return vt.current
}

func (vt *VirtualTimer) Advance(d time.Duration) {
    vt.current = vt.current.Add(d)
}

// Use in tests
func TestCircuitBreaker(t *testing.T) {
    vt := &VirtualTimer{current: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
    
    currentTime := func() time.Time { return vt.Now() }
    
    cb := readerioresult.MakeCircuitBreaker[string](
        currentTime,
        closedState,
        checkError,
        policy,
        logger,
    )
    
    // Test circuit breaker behavior
    // Advance time as needed
    vt.Advance(5 * time.Second)
}
```

## Related Types

- `circuitbreaker.BreakerState`: The internal state of the circuit breaker (closed or open)
- `circuitbreaker.ClosedState`: Configuration for the closed state
- `retry.RetryPolicy`: Policy for retry delays and limits
- `option.Kleisli[error, error]`: Function type for error checking
- `io.Kleisli[string, string]`: Function type for logging

## See Also

- `circuitbreaker` package: Generic circuit breaker implementation
- `retry` package: Retry policies and strategies
- `readerioresult` package: Core ReaderIOResult monad operations