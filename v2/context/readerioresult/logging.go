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

// Package readerioresult provides logging utilities for ReaderIOResult computations.
// It includes functions for entry/exit logging with timing, correlation IDs, and context management.
package readerioresult

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/IBM/fp-go/v2/context/readerio"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/logging"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// loggingContextKeyType is the type used as a key for storing logging information in context.Context
	loggingContextKeyType int

	// LoggingID is a unique identifier assigned to each logged operation for correlation
	LoggingID uint64

	// loggingContext holds the logging state for a computation, including timing,
	// correlation ID, logger instance, and whether logging is enabled.
	loggingContext struct {
		contextID LoggingID    // Unique identifier for this logged operation
		startTime time.Time    // When the operation started (for duration calculation)
		logger    *slog.Logger // The logger instance to use for this operation
		isEnabled bool         // Whether logging is enabled for this operation
	}
)

var (
	// loggingContextKey is the singleton key used to store/retrieve logging data from context
	loggingContextKey loggingContextKeyType

	// loggingCounter is an atomic counter that generates unique LoggingIDs
	loggingCounter atomic.Uint64

	loggingContextValue = F.Bind2nd(context.Context.Value, any(loggingContextKey))

	withLoggingContextValue = F.Bind2of3(context.WithValue)(any(loggingContextKey))

	// getLoggingContext retrieves the logging information (start time and ID) from the context.
	// It returns a Pair containing the start time and the logging ID.
	// This function assumes the context contains logging information; it will panic if not present.
	getLoggingContext = F.Flow3(
		loggingContextValue,
		option.InstanceOf[loggingContext],
		option.GetOrElse(getDefaultLoggingContext),
	)
)

// getDefaultLoggingContext returns a default logging context with the global logger.
// This is used when no logging context is found in the context.Context.
func getDefaultLoggingContext() loggingContext {
	return loggingContext{
		logger: logging.GetLogger(),
	}
}

// withLoggingContext creates an endomorphism that adds a logging context to a context.Context.
// This is used internally to store logging state in the context for retrieval by nested operations.
//
// Parameters:
//   - lctx: The logging context to store
//
// Returns:
//   - An endomorphism that adds the logging context to a context.Context
func withLoggingContext(lctx loggingContext) Endomorphism[context.Context] {
	return F.Bind2nd(withLoggingContextValue, any(lctx))
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
//	                return F.Pipe1(
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
//	                return F.Pipe1(
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
	bracket := F.Bind13of3(readerio.Bracket[context.Context, Result[A], ANY])(onEntry, func(newCtx context.Context, res Result[A]) ReaderIO[ANY] {
		return readerio.FromIO(onExit(res)(newCtx)) // Get the exit callback for this result
	})

	return func(src ReaderIOResult[A]) ReaderIOResult[A] {
		return bracket(F.Flow2(
			src,
			FromIOResult,
		))
	}
}

// onEntry creates a ReaderIO that handles the entry logging for an operation.
// It generates a unique logging ID, captures the start time, and logs the entry message.
// The logging context is stored in the context.Context for later retrieval.
//
// Parameters:
//   - logLevel: The slog.Level to use for logging (e.g., slog.LevelInfo, slog.LevelDebug)
//   - cb: Callback function to retrieve the logger from the context
//   - nameAttr: The slog.Attr containing the operation name
//
// Returns:
//   - A ReaderIO that prepares the context with logging information and logs the entry
func onEntry(
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	nameAttr slog.Attr,
) ReaderIO[context.Context] {

	return func(ctx context.Context) IO[context.Context] {
		// logger
		logger := cb(ctx)

		return func() context.Context {
			// check if the logger is enabled
			if logger.Enabled(ctx, logLevel) {
				// Generate unique logging ID and capture start time
				contextID := LoggingID(loggingCounter.Add(1))
				startTime := time.Now()

				newLogger := logger.With("ID", contextID)

				// log using ID
				newLogger.LogAttrs(ctx, logLevel, "[entering]", nameAttr)

				withCtx := withLoggingContext(loggingContext{
					contextID: contextID,
					startTime: startTime,
					logger:    newLogger,
					isEnabled: true,
				})
				withLogger := logging.WithLogger(newLogger)

				return withCtx(withLogger(ctx))
			}
			// logging disabled
			withCtx := withLoggingContext(loggingContext{
				logger:    logger,
				isEnabled: false,
			})
			return withCtx(ctx)
		}
	}
}

// onExitAny creates a Kleisli function that handles exit logging for an operation.
// It logs either success or error based on the Result, including the operation duration.
// Only logs if logging was enabled during entry (checked via loggingContext.isEnabled).
//
// Parameters:
//   - logLevel: The slog.Level to use for logging
//   - nameAttr: The slog.Attr containing the operation name
//
// Returns:
//   - A Kleisli function that logs the exit/error and returns nil
func onExitAny(
	logLevel slog.Level,
	nameAttr slog.Attr,
) readerio.Kleisli[Result[any], any] {
	return func(res Result[any]) ReaderIO[any] {
		return func(ctx context.Context) IO[any] {
			value := getLoggingContext(ctx)

			if value.isEnabled {

				return func() any {
					// Retrieve logging information from context
					durationAttr := slog.Duration("duration", time.Since(value.startTime))

					// Log error with ID and duration
					onError := func(err error) any {
						value.logger.LogAttrs(ctx, logLevel, "[throwing]",
							nameAttr,
							durationAttr,
							slog.Any("error", err))
						return nil
					}

					// Log success with ID and duration
					onSuccess := func(_ any) any {
						value.logger.LogAttrs(ctx, logLevel, "[exiting ]", nameAttr, durationAttr)
						return nil
					}

					return F.Pipe1(
						res,
						result.Fold(onError, onSuccess),
					)
				}
			}
			// nothing to do
			return io.Of[any](nil)
		}
	}
}

// LogEntryExitWithCallback creates an operator that logs entry and exit of a ReaderIOResult computation
// using a custom logger callback and log level. This provides more control than LogEntryExit.
//
// This function allows you to:
//   - Use a custom log level (Debug, Info, Warn, Error)
//   - Retrieve the logger from the context using a custom callback
//   - Control whether logging is enabled based on the logger's configuration
//
// Type Parameters:
//   - A: The success type of the ReaderIOResult
//
// Parameters:
//   - logLevel: The slog.Level to use for all log messages (entry, exit, error)
//   - cb: Callback function to retrieve the *slog.Logger from the context
//   - name: A descriptive name for the operation
//
// Returns:
//   - An Operator that wraps the ReaderIOResult with customizable logging
//
// Example with custom log level:
//
//	// Log at debug level
//	debugOp := LogEntryExitWithCallback[User](
//	    slog.LevelDebug,
//	    logging.GetLoggerFromContext,
//	    "fetchUser",
//	)
//	result := debugOp(fetchUser(123))
//
// Example with custom logger callback:
//
//	type loggerKey int
//	const myLoggerKey loggerKey = 0
//
//	getMyLogger := func(ctx context.Context) *slog.Logger {
//	    if logger := ctx.Value(myLoggerKey); logger != nil {
//	        return logger.(*slog.Logger)
//	    }
//	    return slog.Default()
//	}
//
//	customOp := LogEntryExitWithCallback[Data](
//	    slog.LevelInfo,
//	    getMyLogger,
//	    "processData",
//	)
func LogEntryExitWithCallback[A any](
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	name string) Operator[A, A] {

	nameAttr := slog.String("name", name)

	return LogEntryExitF(
		onEntry(logLevel, cb, nameAttr),
		F.Flow2(
			result.MapTo[A, any](nil),
			onExitAny(logLevel, nameAttr),
		),
	)
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
//	result := loggedFetch(t.Context())()
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
//	result := logged(t.Context())()
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
//	result := pipeline(t.Context())()
//	// Logs:
//	// [entering 3] fetchUser
//	// [exiting  3] fetchUser [0.1s]
//	// [entering 4] fetchOrders
//	// [exiting  4] fetchOrders [0.2s]
//
// Example with concurrent operations:
//
//	// Multiple operations can run concurrently, each with unique IDs
//	op1 := LogEntryExit[Data]("operation1")(fetchData(1))
//	op2 := LogEntryExit[Data]("operation2")(fetchData(2))
//
//	go op1(t.Context())()
//	go op2(t.Context())()
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
//go:inline
func LogEntryExit[A any](name string) Operator[A, A] {
	return LogEntryExitWithCallback[A](slog.LevelInfo, logging.GetLoggerFromContext, name)
}

func curriedLog(
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	message string) func(slog.Attr) func(context.Context) func() struct{} {
	return F.Curry2(func(a slog.Attr, ctx context.Context) func() struct{} {
		logger := cb(ctx)
		return func() struct{} {
			logger.LogAttrs(ctx, logLevel, message, a)
			return struct{}{}
		}
	})
}

// SLogWithCallback creates a Kleisli arrow that logs a Result value (success or error) with a custom logger and log level.
//
// This function logs both successful values and errors, making it useful for debugging and monitoring
// Result values as they flow through a computation. Unlike TapSLog which only logs successful values,
// SLogWithCallback logs the Result regardless of whether it contains a value or an error.
//
// The logged output includes:
//   - For success: The message with the value as a structured "value" attribute
//   - For error: The message with the error as a structured "error" attribute
//
// The Result is passed through unchanged after logging.
//
// Type Parameters:
//   - A: The success type of the Result
//
// Parameters:
//   - logLevel: The slog.Level to use for logging (e.g., slog.LevelInfo, slog.LevelDebug)
//   - cb: Callback function to retrieve the *slog.Logger from the context
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - A Kleisli arrow that logs the Result (value or error) and returns it unchanged
//
// Example with custom log level:
//
//	debugLog := SLogWithCallback[User](
//	    slog.LevelDebug,
//	    logging.GetLoggerFromContext,
//	    "User result",
//	)
//
//	pipeline := F.Pipe2(
//	    fetchUser(123),
//	    Chain(debugLog),
//	    Map(func(u User) string { return u.Name }),
//	)
//
// Example with custom logger:
//
//	type loggerKey int
//	const myLoggerKey loggerKey = 0
//
//	getMyLogger := func(ctx context.Context) *slog.Logger {
//	    if logger := ctx.Value(myLoggerKey); logger != nil {
//	        return logger.(*slog.Logger)
//	    }
//	    return slog.Default()
//	}
//
//	customLog := SLogWithCallback[Data](
//	    slog.LevelWarn,
//	    getMyLogger,
//	    "Data processing result",
//	)
//
// Use Cases:
//   - Debugging: Log both successful and failed Results in a pipeline
//   - Error tracking: Monitor error occurrences with custom log levels
//   - Custom logging: Use application-specific loggers and log levels
//   - Conditional logging: Enable/disable logging based on logger configuration
func SLogWithCallback[A any](
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	message string) Kleisli[Result[A], A] {

	return F.Pipe1(
		F.Flow2(
			// create the attribute to log depending on the condition
			result.ToSLogAttr[A](),
			// create an `IO` that logs the attribute
			curriedLog(logLevel, cb, message),
		),
		// preserve the original context
		reader.Chain(reader.Sequence(readerio.MapTo[struct{}, Result[A]])),
	)
}

// SLog creates a Kleisli arrow that logs a Result value (success or error) with a message.
//
// This function logs both successful values and errors at Info level using the logger from the context.
// It's a convenience wrapper around SLogWithCallback with standard settings.
//
// The logged output includes:
//   - For success: The message with the value as a structured "value" attribute
//   - For error: The message with the error as a structured "error" attribute
//
// The Result is passed through unchanged after logging, making this function transparent in the
// computation pipeline.
//
// Type Parameters:
//   - A: The success type of the Result
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - A Kleisli arrow that logs the Result (value or error) and returns it unchanged
//
// Example with successful Result:
//
//	pipeline := F.Pipe2(
//	    fetchUser(123),
//	    Chain(SLog[User]("Fetched user")),
//	    Map(func(u User) string { return u.Name }),
//	)
//
//	result := pipeline(t.Context())()
//	// If successful, logs: "Fetched user" value={ID:123 Name:"Alice"}
//	// If error, logs: "Fetched user" error="user not found"
//
// Example in error handling pipeline:
//
//	pipeline := F.Pipe3(
//	    fetchData(id),
//	    Chain(SLog[Data]("Data fetched")),
//	    Chain(validateData),
//	    Chain(SLog[Data]("Data validated")),
//	    Chain(processData),
//	)
//
//	// Logs each step, including errors:
//	// "Data fetched" value={...} or error="..."
//	// "Data validated" value={...} or error="..."
//
// Use Cases:
//   - Debugging: Track both successful and failed Results in a pipeline
//   - Error monitoring: Log errors as they occur in the computation
//   - Flow tracking: See the progression of Results through a pipeline
//   - Troubleshooting: Identify where errors are introduced or propagated
//
// Note: This function logs the Result itself (which may contain an error), not just successful values.
// For logging only successful values, use TapSLog instead.
//
//go:inline
func SLog[A any](message string) Kleisli[Result[A], A] {
	return SLogWithCallback[A](slog.LevelInfo, logging.GetLoggerFromContext, message)
}

// TapSLog creates an operator that logs only successful values with a message and passes them through unchanged.
//
// This function is useful for debugging and monitoring values as they flow through a ReaderIOResult
// computation chain. Unlike SLog which logs both successes and errors, TapSLog only logs when the
// computation is successful. If the computation contains an error, no logging occurs and the error
// is propagated unchanged.
//
// The logged output includes:
//   - The provided message
//   - The value being passed through (as a structured "value" attribute)
//
// Type Parameters:
//   - A: The type of the value to log and pass through
//
// Parameters:
//   - message: A descriptive message to include in the log entry
//
// Returns:
//   - An Operator that logs successful values and returns them unchanged
//
// Example with simple value logging:
//
//	fetchUser := func(id int) ReaderIOResult[User] {
//	    return Of(User{ID: id, Name: "Alice"})
//	}
//
//	pipeline := F.Pipe2(
//	    fetchUser(123),
//	    TapSLog[User]("Fetched user"),
//	    Map(func(u User) string { return u.Name }),
//	)
//
//	result := pipeline(t.Context())()
//	// Logs: "Fetched user" value={ID:123 Name:"Alice"}
//	// Returns: result.Of("Alice")
//
// Example in a processing pipeline:
//
//	processOrder := F.Pipe4(
//	    fetchOrder(orderId),
//	    TapSLog[Order]("Order fetched"),
//	    Chain(validateOrder),
//	    TapSLog[Order]("Order validated"),
//	    Chain(processPayment),
//	    TapSLog[Payment]("Payment processed"),
//	)
//
//	result := processOrder(t.Context())()
//	// Logs each successful step with the intermediate values
//	// If any step fails, subsequent TapSLog calls don't log
//
// Example with error handling:
//
//	pipeline := F.Pipe3(
//	    fetchData(id),
//	    TapSLog[Data]("Data fetched"),
//	    Chain(func(d Data) ReaderIOResult[Result] {
//	        if d.IsValid() {
//	            return Of(processData(d))
//	        }
//	        return Left[Result](errors.New("invalid data"))
//	    }),
//	    TapSLog[Result]("Data processed"),
//	)
//
//	// If fetchData succeeds: logs "Data fetched" with the data
//	// If processing succeeds: logs "Data processed" with the result
//	// If processing fails: "Data processed" is NOT logged (error propagates)
//
// Use Cases:
//   - Debugging: Inspect intermediate successful values in a computation pipeline
//   - Monitoring: Track successful data flow through complex operations
//   - Troubleshooting: Identify where successful computations stop (last logged value before error)
//   - Auditing: Log important successful values for compliance or security
//   - Development: Understand data transformations during development
//
// Note: This function only logs successful values. Errors are silently propagated without logging.
// For logging both successes and errors, use SLog instead.
//
//go:inline
func TapSLog[A any](message string) Operator[A, A] {
	return readerio.ChainFirst(SLog[A](message))
}
