// Package readerioresult provides logging utilities for ReaderIOResult computations.
// It includes functions for entry/exit logging with timing, correlation IDs, and context management.
package readerioresult

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/IBM/fp-go/v2/context/readerio"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// loggingContextKeyType is the type used as a key for storing logging information in context.Context
	loggingContextKeyType int

	// LoggingID is a unique identifier assigned to each logged operation for correlation
	LoggingID uint64
)

var (
	// loggingContextKey is the singleton key used to store/retrieve logging data from context
	loggingContextKey loggingContextKeyType

	// loggingCounter is an atomic counter that generates unique LoggingIDs
	loggingCounter atomic.Uint64

	// getLoggingContext retrieves the logging information (start time and ID) from the context.
	// It returns a Pair containing the start time and the logging ID.
	// This function assumes the context contains logging information; it will panic if not present.
	getLoggingContext = function.Flow3(
		function.Bind2nd(context.Context.Value, any(loggingContextKey)),
		option.ToType[pair.Pair[time.Time, LoggingID]],
		option.GetOrElse(function.Zero[pair.Pair[time.Time, LoggingID]]),
	)

	// getLoggingID extracts just the LoggingID from the context, discarding the start time.
	// This is a convenience function composed from getLoggingContext and pair.Tail.
	getLoggingID = function.Flow2(
		getLoggingContext,
		pair.Tail,
	)
)

// WithLoggingID wraps a value with its associated LoggingID from the current context.
//
// This function retrieves the LoggingID from the context and pairs it with the provided value,
// creating a ReaderIOResult that produces a Pair[LoggingID, A]. This is useful when you need
// to correlate a value with the logging ID of the operation that produced it.
//
// Type Parameters:
//   - A: The type of the value to be paired with the logging ID
//
// Parameters:
//   - src: The value to be paired with the logging ID
//
// Returns:
//   - A ReaderIOResult that produces a Pair containing the LoggingID and the source value
//
// Example:
//
//	fetchUser := func(id int) ReaderIOResult[User] {
//	    return Of(User{ID: id, Name: "Alice"})
//	}
//
//	// Wrap the result with its logging ID
//	withID := F.Pipe2(
//	    fetchUser(123),
//	    LogEntryExit[User]("fetchUser"),
//	    Chain(WithLoggingID[User]),
//	)
//
//	result := withID(ctx)() // Returns Result[Pair[LoggingID, User]]
//	// Can now correlate the user with the operation that fetched it
//
// Use Cases:
//   - Correlating results with the operations that produced them
//   - Tracking data lineage through complex pipelines
//   - Debugging by associating values with their source operations
//   - Audit logging with operation correlation
func WithLoggingID[A any](src A) ReaderIOResult[pair.Pair[LoggingID, A]] {
	return function.Pipe1(
		Ask(),
		Map(function.Flow2(
			getLoggingID,
			pair.FromTail[LoggingID](src),
		)),
	)
}

// LogEntryExitF creates a customizable operator that wraps a ReaderIOResult computation with entry/exit callbacks.
//
// This is a more flexible version of LogEntryExit that allows you to provide custom callbacks for
// entry and exit events. The onEntry callback receives the current context and can return a modified
// context (e.g., with additional logging information). The onExit callback receives the computation
// result and can perform custom logging, metrics collection, or cleanup.
//
// The function uses the bracket pattern to ensure that:
//   - The onEntry callback is executed before the computation starts
//   - The computation runs with the context returned by onEntry
//   - The onExit callback is executed after the computation completes (success or failure)
//   - The original result is preserved and returned unchanged
//   - Cleanup happens even if the computation fails
//
// Type Parameters:
//   - A: The success type of the ReaderIOResult
//   - ANY: The return type of the onExit callback (typically any)
//
// Parameters:
//   - onEntry: A ReaderIO that receives the current context and returns a (possibly modified) context.
//     This is executed before the computation starts. Use this for logging entry, adding context values,
//     starting timers, or initialization logic.
//   - onExit: A Kleisli function that receives the Result[A] and returns a ReaderIO[ANY].
//     This is executed after the computation completes, regardless of success or failure.
//     Use this for logging exit, recording metrics, cleanup, or finalization logic.
//
// Returns:
//   - An Operator that wraps the ReaderIOResult computation with the custom entry/exit callbacks
//
// Example with custom context modification:
//
//	type RequestID string
//
//	logOp := LogEntryExitF[User, any](
//	    func(ctx context.Context) IO[context.Context] {
//	        return func() context.Context {
//	            reqID := RequestID(uuid.New().String())
//	            log.Printf("[%s] Starting operation", reqID)
//	            return context.WithValue(ctx, "requestID", reqID)
//	        }
//	    },
//	    func(res Result[User]) ReaderIO[any] {
//	        return func(ctx context.Context) IO[any] {
//	            return func() any {
//	                reqID := ctx.Value("requestID").(RequestID)
//	                return function.Pipe1(
//	                    res,
//	                    result.Fold(
//	                        func(err error) any {
//	                            log.Printf("[%s] Operation failed: %v", reqID, err)
//	                            return nil
//	                        },
//	                        func(_ User) any {
//	                            log.Printf("[%s] Operation succeeded", reqID)
//	                            return nil
//	                        },
//	                    ),
//	                )
//	            }
//	        }
//	    },
//	)
//
//	wrapped := logOp(fetchUser(123))
//
// Example with metrics collection:
//
//	import "github.com/prometheus/client_golang/prometheus"
//
//	metricsOp := LogEntryExitF[Response, any](
//	    func(ctx context.Context) IO[context.Context] {
//	        return func() context.Context {
//	            requestCount.WithLabelValues("api_call", "started").Inc()
//	            return context.WithValue(ctx, "startTime", time.Now())
//	        }
//	    },
//	    func(res Result[Response]) ReaderIO[any] {
//	        return func(ctx context.Context) IO[any] {
//	            return func() any {
//	                startTime := ctx.Value("startTime").(time.Time)
//	                duration := time.Since(startTime).Seconds()
//
//	                return function.Pipe1(
//	                    res,
//	                    result.Fold(
//	                        func(err error) any {
//	                            requestCount.WithLabelValues("api_call", "error").Inc()
//	                            requestDuration.WithLabelValues("api_call", "error").Observe(duration)
//	                            return nil
//	                        },
//	                        func(_ Response) any {
//	                            requestCount.WithLabelValues("api_call", "success").Inc()
//	                            requestDuration.WithLabelValues("api_call", "success").Observe(duration)
//	                            return nil
//	                        },
//	                    ),
//	                )
//	            }
//	        }
//	    },
//	)
//
// Use Cases:
//   - Custom context modification: Adding request IDs, trace IDs, or other context values
//   - Structured logging: Integration with zap, logrus, or other structured loggers
//   - Metrics collection: Recording operation durations, success/failure rates
//   - Distributed tracing: OpenTelemetry, Jaeger integration
//   - Custom monitoring: Application-specific monitoring and alerting
//
// Note: LogEntryExit is implemented using LogEntryExitF with standard logging and context management.
// Use LogEntryExitF when you need more control over the entry/exit behavior or context modification.
func LogEntryExitF[A, ANY any](
	onEntry ReaderIO[context.Context],
	onExit readerio.Kleisli[Result[A], ANY],
) Operator[A, A] {
	bracket := function.Bind13of3(readerio.Bracket[context.Context, Result[A], ANY])(onEntry, func(newCtx context.Context, res Result[A]) ReaderIO[ANY] {
		return readerio.FromIO(onExit(res)(newCtx)) // Get the exit callback for this result
	})

	return func(src ReaderIOResult[A]) ReaderIOResult[A] {
		return bracket(function.Flow2(
			src,
			FromIOResult,
		))
	}
}

// LogEntryExit creates an operator that logs the entry and exit of a ReaderIOResult computation with timing and correlation IDs.
//
// This function wraps a ReaderIOResult computation with automatic logging that tracks:
//   - Entry: Logs when the computation starts with "[entering <id>] <name>"
//   - Exit: Logs when the computation completes successfully with "[exiting  <id>] <name> [duration]"
//   - Error: Logs when the computation fails with "[throwing <id>] <name> [duration]: <error>"
//
// Each logged operation is assigned a unique LoggingID (a monotonically increasing counter) that
// appears in all log messages for that operation. This ID enables correlation of entry and exit
// logs, even when multiple operations are running concurrently or are interleaved.
//
// The logging information (start time and ID) is stored in the context and can be retrieved using
// getLoggingContext or getLoggingID. This allows nested operations to access the parent operation's
// logging information.
//
// Type Parameters:
//   - A: The success type of the ReaderIOResult
//
// Parameters:
//   - name: A descriptive name for the computation, used in log messages to identify the operation
//
// Returns:
//   - An Operator that wraps the ReaderIOResult computation with entry/exit logging
//
// The function uses the bracket pattern to ensure that:
//   - Entry is logged before the computation starts
//   - A unique LoggingID is assigned and stored in the context
//   - Exit/error is logged after the computation completes, regardless of success or failure
//   - Timing is accurate, measuring from entry to exit
//   - The original result is preserved and returned unchanged
//
// Log Format:
//   - Entry:   "[entering <id>] <name>"
//   - Success: "[exiting  <id>] <name> [<duration>s]"
//   - Error:   "[throwing <id>] <name> [<duration>s]: <error>"
//
// Example with successful computation:
//
//	fetchUser := func(id int) ReaderIOResult[User] {
//	    return Of(User{ID: id, Name: "Alice"})
//	}
//
//	// Wrap with logging
//	loggedFetch := LogEntryExit[User]("fetchUser")(fetchUser(123))
//
//	// Execute
//	result := loggedFetch(context.Background())()
//	// Logs:
//	// [entering 1] fetchUser
//	// [exiting  1] fetchUser [0.1s]
//
// Example with error:
//
//	failingOp := func() ReaderIOResult[string] {
//	    return Left[string](errors.New("connection timeout"))
//	}
//
//	logged := LogEntryExit[string]("failingOp")(failingOp())
//	result := logged(context.Background())()
//	// Logs:
//	// [entering 2] failingOp
//	// [throwing 2] failingOp [0.0s]: connection timeout
//
// Example with nested operations:
//
//	fetchOrders := func(userID int) ReaderIOResult[[]Order] {
//	    return Of([]Order{{ID: 1}})
//	}
//
//	pipeline := F.Pipe3(
//	    fetchUser(123),
//	    LogEntryExit[User]("fetchUser"),
//	    Chain(func(user User) ReaderIOResult[[]Order] {
//	        return fetchOrders(user.ID)
//	    }),
//	    LogEntryExit[[]Order]("fetchOrders"),
//	)
//
//	result := pipeline(context.Background())()
//	// Logs:
//	// [entering 3] fetchUser
//	// [exiting  3] fetchUser [0.1s]
//	// Fetching orders for user (parent operation: 3)
//	// [entering 4] fetchOrders
//	// [exiting  4] fetchOrders [0.2s]
//
// Example with concurrent operations:
//
//	// Multiple operations can run concurrently, each with unique IDs
//	op1 := LogEntryExit[Data]("operation1")(fetchData(1))
//	op2 := LogEntryExit[Data]("operation2")(fetchData(2))
//
//	go op1(context.Background())()
//	go op2(context.Background())()
//	// Logs (order may vary):
//	// [entering 5] operation1
//	// [entering 6] operation2
//	// [exiting  5] operation1 [0.1s]
//	// [exiting  6] operation2 [0.2s]
//	// The IDs allow correlation even when logs are interleaved
//
// Use Cases:
//   - Debugging: Track execution flow through complex ReaderIOResult chains with correlation IDs
//   - Performance monitoring: Identify slow operations with timing information
//   - Production logging: Monitor critical operations with unique identifiers
//   - Concurrent operations: Correlate logs from multiple concurrent operations
//   - Nested operations: Track parent-child relationships in operation hierarchies
//   - Troubleshooting: Quickly identify where errors occur and correlate with entry logs
//
// Note: This function uses Go's standard log package and a global atomic counter for IDs.
// For production systems, consider using a structured logging library and adapting this
// pattern to support different log levels, structured fields, and distributed tracing.
func LogEntryExit[A any](name string) Operator[A, A] {
	return LogEntryExitF(
		func(ctx context.Context) IO[context.Context] {
			return func() context.Context {
				// Generate unique logging ID and capture start time
				counter := LoggingID(loggingCounter.Add(1))
				tStart := time.Now()

				// Log entry with unique ID
				log.Printf("[entering %d] %s", counter, name)

				// Store logging information in context for later retrieval
				return context.WithValue(ctx, loggingContextKey, pair.MakePair(tStart, counter))
			}
		},
		func(res Result[A]) ReaderIO[any] {
			return func(ctx context.Context) IO[any] {
				value := getLoggingContext(ctx)
				counter := pair.Tail(value)

				return func() any {
					// Retrieve logging information from context
					duration := time.Since(pair.Head(value)).Seconds()

					// Log error with ID and duration
					onError := func(err error) any {
						log.Printf("[throwing %d] %s [%.1fs]: %v", counter, name, duration, err)
						return nil
					}

					// Log success with ID and duration
					onSuccess := func(_ A) any {
						log.Printf("[exiting  %d] %s [%.1fs]", counter, name, duration)
						return nil
					}

					return function.Pipe1(
						res,
						result.Fold(onError, onSuccess),
					)
				}
			}
		},
	)
}
