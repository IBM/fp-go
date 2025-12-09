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

package ioeither

import (
	"log"
	"time"

	"github.com/IBM/fp-go/v2/bytes"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/json"
	"github.com/IBM/fp-go/v2/pair"
)

// LogJSON converts the argument to pretty printed JSON and then logs it via the format string
// Can be used with [ChainFirst] and [Tap]
func LogJSON[A any](prefix string) Kleisli[error, A, string] {
	return function.Flow4(
		json.MarshalIndent[A],
		either.Map[error](bytes.ToString),
		FromEither[error, string],
		ChainIOK[error](io.Logf[string](prefix)),
	)
}

// LogEntryExitF creates a customizable operator that wraps an IOEither computation with entry/exit callbacks.
//
// This is a more flexible version of LogEntryExit that allows you to provide custom callbacks for
// entry and exit events. The onEntry callback is executed before the computation starts and can
// return a "start token" (such as a timestamp, trace ID, or any context data). This token is then
// passed to the onExit callback along with the computation result, enabling correlation between
// entry and exit events.
//
// The function uses the bracket pattern to ensure that:
//   - The onEntry callback is executed before the computation starts
//   - The onExit callback is executed after the computation completes (success or failure)
//   - The start token from onEntry is available in onExit for correlation
//   - The original result is preserved and returned unchanged
//   - Cleanup happens even if the computation fails
//
// Type Parameters:
//   - E: The error type (Left value) of the IOEither
//   - A: The success type (Right value) of the IOEither
//   - STARTTOKEN: The type of the token returned by onEntry (e.g., time.Time, string, trace.Span)
//   - ANY: The return type of the onExit callback (typically any or a specific type)
//
// Parameters:
//   - onEntry: An IO action executed when the computation starts. Returns a STARTTOKEN that will
//     be passed to onExit. Use this for logging entry, starting timers, creating trace spans, etc.
//   - onExit: A Kleisli function that receives a Pair containing:
//   - Head: STARTTOKEN - the token returned by onEntry
//   - Tail: Either[E, A] - the result of the computation (Left for error, Right for success)
//     Use this for logging exit, recording metrics, closing spans, or cleanup logic.
//
// Returns:
//   - An Operator that wraps the IOEither computation with the custom entry/exit callbacks
//
// Example with timing (as used by LogEntryExit):
//
//	logOp := LogEntryExitF[error, User, time.Time, any](
//	    func() time.Time {
//	        log.Printf("[entering] fetchUser")
//	        return time.Now()  // Start token is the start time
//	    },
//	    func(res pair.Pair[time.Time, Either[error, User]]) IO[any] {
//	        startTime := pair.Head(res)
//	        result := pair.Tail(res)
//	        duration := time.Since(startTime).Seconds()
//
//	        return func() any {
//	            if either.IsLeft(result) {
//	                log.Printf("[throwing] fetchUser [%.1fs]: %v", duration, either.GetLeft(result))
//	            } else {
//	                log.Printf("[exiting] fetchUser [%.1fs]", duration)
//	            }
//	            return nil
//	        }
//	    },
//	)
//
//	wrapped := logOp(fetchUser(123))
//
// Example with distributed tracing:
//
//	import "go.opentelemetry.io/otel/trace"
//
//	tracer := otel.Tracer("my-service")
//
//	traceOp := LogEntryExitF[error, Data, trace.Span, any](
//	    func() trace.Span {
//	        _, span := tracer.Start(ctx, "fetchData")
//	        return span  // Start token is the span
//	    },
//	    func(res pair.Pair[trace.Span, Either[error, Data]]) IO[any] {
//	        span := pair.Head(res)  // Get the span from entry
//	        result := pair.Tail(res)
//
//	        return func() any {
//	            if either.IsLeft(result) {
//	                span.RecordError(either.GetLeft(result))
//	                span.SetStatus(codes.Error, "operation failed")
//	            } else {
//	                span.SetStatus(codes.Ok, "operation succeeded")
//	            }
//	            span.End()  // Close the span
//	            return nil
//	        }
//	    },
//	)
//
// Example with correlation ID:
//
//	type RequestContext struct {
//	    CorrelationID string
//	    StartTime     time.Time
//	}
//
//	correlationOp := LogEntryExitF[error, Response, RequestContext, any](
//	    func() RequestContext {
//	        ctx := RequestContext{
//	            CorrelationID: uuid.New().String(),
//	            StartTime:     time.Now(),
//	        }
//	        log.Printf("[%s] Request started", ctx.CorrelationID)
//	        return ctx
//	    },
//	    func(res pair.Pair[RequestContext, Either[error, Response]]) IO[any] {
//	        ctx := pair.Head(res)
//	        result := pair.Tail(res)
//	        duration := time.Since(ctx.StartTime)
//
//	        return func() any {
//	            if either.IsLeft(result) {
//	                log.Printf("[%s] Request failed after %v: %v",
//	                    ctx.CorrelationID, duration, either.GetLeft(result))
//	            } else {
//	                log.Printf("[%s] Request completed after %v",
//	                    ctx.CorrelationID, duration)
//	            }
//	            return nil
//	        }
//	    },
//	)
//
// Example with metrics collection:
//
//	import "github.com/prometheus/client_golang/prometheus"
//
//	type MetricsToken struct {
//	    StartTime time.Time
//	    OpName    string
//	}
//
//	metricsOp := LogEntryExitF[error, Result, MetricsToken, any](
//	    func() MetricsToken {
//	        token := MetricsToken{
//	            StartTime: time.Now(),
//	            OpName:    "api_call",
//	        }
//	        requestCount.WithLabelValues(token.OpName, "started").Inc()
//	        return token
//	    },
//	    func(res pair.Pair[MetricsToken, Either[error, Result]]) IO[any] {
//	        token := pair.Head(res)
//	        result := pair.Tail(res)
//	        duration := time.Since(token.StartTime).Seconds()
//
//	        return func() any {
//	            if either.IsLeft(result) {
//	                requestCount.WithLabelValues(token.OpName, "error").Inc()
//	                requestDuration.WithLabelValues(token.OpName, "error").Observe(duration)
//	            } else {
//	                requestCount.WithLabelValues(token.OpName, "success").Inc()
//	                requestDuration.WithLabelValues(token.OpName, "success").Observe(duration)
//	            }
//	            return nil
//	        }
//	    },
//	)
//
// Use Cases:
//   - Structured logging: Integration with zap, logrus, or other structured loggers
//   - Distributed tracing: OpenTelemetry, Jaeger, Zipkin integration with span management
//   - Metrics collection: Recording operation durations, success/failure rates with Prometheus
//   - Request correlation: Tracking requests across service boundaries with correlation IDs
//   - Custom monitoring: Application-specific monitoring and alerting
//   - Audit logging: Recording detailed operation information for compliance
//
// Note: LogEntryExit is implemented using LogEntryExitF with time.Time as the start token.
// Use LogEntryExitF when you need more control over the entry/exit behavior or need to
// pass custom context between entry and exit callbacks.
func LogEntryExitF[E, A, STARTTOKEN, ANY any](
	onEntry IO[STARTTOKEN],
	onExit io.Kleisli[pair.Pair[STARTTOKEN, Either[E, A]], ANY],
) Operator[E, A, A] {

	// release: Invokes the onExit callback with the start token and computation result
	// This function is called by the bracket pattern after the computation completes,
	// regardless of whether it succeeded or failed. It pairs the start token (from onEntry)
	// with the computation result and passes them to the onExit callback.
	release := func(start pair.Pair[STARTTOKEN, IOEither[E, A]], result Either[E, A]) IO[ANY] {
		return function.Pipe1(
			pair.MakePair(pair.Head(start), result), // Pair the start token with the result
			onExit,                                  // Pass to the exit callback
		)
	}

	return func(src IOEither[E, A]) IOEither[E, A] {
		return io.Bracket(
			// Acquire: Execute onEntry to get the start token, then pair it with the source IOEither
			function.Pipe1(
				onEntry,                                // Execute entry callback to get start token
				io.Map(pair.FromTail[STARTTOKEN](src)), // Pair the token with the source computation
			),
			// Use: Extract and execute the IOEither computation from the pair
			pair.Tail[STARTTOKEN, IOEither[E, A]],
			// Release: Call onExit with the start token and result (always executed)
			release,
		)

	}
}

// LogEntryExit creates an operator that logs the entry and exit of an IOEither computation with timing information.
//
// This function wraps an IOEither computation with automatic logging that tracks:
//   - Entry: Logs when the computation starts with "[entering] <name>"
//   - Exit: Logs when the computation completes successfully with "[exiting ] <name> [duration]"
//   - Error: Logs when the computation fails with "[throwing] <name> [duration]: <error>"
//
// The duration is measured in seconds with one decimal place precision (e.g., "2.5s").
// This is particularly useful for debugging, performance monitoring, and understanding the
// execution flow of complex IOEither chains.
//
// Type Parameters:
//   - E: The error type (Left value) of the IOEither
//   - A: The success type (Right value) of the IOEither
//
// Parameters:
//   - name: A descriptive name for the computation, used in log messages to identify the operation
//
// Returns:
//   - An Operator that wraps the IOEither computation with entry/exit logging
//
// The function uses the bracket pattern to ensure that:
//   - Entry is logged before the computation starts
//   - Exit/error is logged after the computation completes, regardless of success or failure
//   - Timing is accurate, measuring from entry to exit
//   - The original result is preserved and returned unchanged
//
// Log Format:
//   - Entry:   "[entering] <name>"
//   - Success: "[exiting ] <name> [<duration>s]"
//   - Error:   "[throwing] <name> [<duration>s]: <error>"
//
// Example with successful computation:
//
//	fetchUser := func(id int) IOEither[error, User] {
//	    return TryCatch(func() (User, error) {
//	        // Simulate database query
//	        time.Sleep(100 * time.Millisecond)
//	        return User{ID: id, Name: "Alice"}, nil
//	    })
//	}
//
//	// Wrap with logging
//	loggedFetch := LogEntryExit[error, User]("fetchUser")(fetchUser(123))
//
//	// Execute
//	result := loggedFetch()
//	// Logs:
//	// [entering] fetchUser
//	// [exiting ] fetchUser [0.1s]
//
// Example with error:
//
//	failingOp := func() IOEither[error, string] {
//	    return TryCatch(func() (string, error) {
//	        time.Sleep(50 * time.Millisecond)
//	        return "", errors.New("connection timeout")
//	    })
//	}
//
//	logged := LogEntryExit[error, string]("failingOp")(failingOp())
//	result := logged()
//	// Logs:
//	// [entering] failingOp
//	// [throwing] failingOp [0.1s]: connection timeout
//
// Example with chained operations:
//
//	pipeline := F.Pipe3(
//	    fetchUser(123),
//	    LogEntryExit[error, User]("fetchUser"),
//	    Chain(func(user User) IOEither[error, []Order] {
//	        return fetchOrders(user.ID)
//	    }),
//	    LogEntryExit[error, []Order]("fetchOrders"),
//	)
//	// Logs each step with timing:
//	// [entering] fetchUser
//	// [exiting ] fetchUser [0.1s]
//	// [entering] fetchOrders
//	// [exiting ] fetchOrders [0.2s]
//
// Example for performance monitoring:
//
//	slowQuery := func() IOEither[error, []Record] {
//	    return TryCatch(func() ([]Record, error) {
//	        // Simulate slow database query
//	        time.Sleep(2 * time.Second)
//	        return []Record{{ID: 1}}, nil
//	    })
//	}
//
//	monitored := LogEntryExit[error, []Record]("slowQuery")(slowQuery())
//	result := monitored()
//	// Logs:
//	// [entering] slowQuery
//	// [exiting ] slowQuery [2.0s]
//	// Helps identify performance bottlenecks
//
// Example with custom error types:
//
//	type AppError struct {
//	    Code    int
//	    Message string
//	}
//
//	func (e AppError) Error() string {
//	    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
//	}
//
//	operation := func() IOEither[AppError, Data] {
//	    return Left[Data](AppError{Code: 404, Message: "Not Found"})
//	}
//
//	logged := LogEntryExit[AppError, Data]("operation")(operation())
//	result := logged()
//	// Logs:
//	// [entering] operation
//	// [throwing] operation [0.0s]: Error 404: Not Found
//
// Use Cases:
//   - Debugging: Track execution flow through complex IOEither chains
//   - Performance monitoring: Identify slow operations with timing information
//   - Production logging: Monitor critical operations in production systems
//   - Testing: Verify that operations are executed in the expected order
//   - Troubleshooting: Quickly identify where errors occur in a pipeline
//
// Note: This function uses Go's standard log package. For production systems,
// consider using a structured logging library and adapting this pattern to
// support different log levels and structured fields.
func LogEntryExit[E, A any](name string) Operator[E, A, A] {

	return LogEntryExitF(
		func() time.Time {
			log.Printf("[entering] %s", name)
			return time.Now()
		},
		func(res pair.Pair[time.Time, Either[E, A]]) IO[any] {

			duration := time.Since(pair.Head(res)).Seconds()

			return func() any {

				onError := func(err E) any {
					log.Printf("[throwing] %s [%.1fs]: %v", name, duration, err)
					return nil
				}

				onSuccess := func(_ A) any {
					log.Printf("[exiting ] %s [%.1fs]", name, duration)
					return nil
				}

				return function.Pipe2(
					res,
					pair.Tail,
					either.Fold(onError, onSuccess),
				)
			}
		},
	)
}
