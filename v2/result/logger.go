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

package result

import (
	"log"
	"log/slog"

	"github.com/IBM/fp-go/v2/either"
)

// Logger creates a logging function for Result values that logs both error and success cases.
// The function logs the value and then returns the original Result unchanged.
//
// This is a specialized version of either.Logger where the Left (error) type is fixed to error.
// It provides a convenient way to add logging to Result-based computations without affecting
// the computation's outcome. The logger is particularly useful for debugging and monitoring
// Result pipelines.
//
// Type Parameters:
//   - A: The type of the success value (Right side of the Result)
//
// Parameters:
//   - loggers: Optional log.Logger instances. If none provided, uses the default logger.
//
// Returns:
//   - A function that takes a prefix string and returns an Operator[A, A] that logs and passes through the Result
//
// Behavior:
//   - For Ok(value): Logs the success value with the given prefix
//   - For Err(error): Logs the error with the given prefix
//   - Always returns the original Result unchanged
//
// Example with success value:
//
//	logger := result.Logger[int]()
//	result := F.Pipe2(
//	    result.Of(42),
//	    logger("Processing"),  // Logs: "Processing: 42"
//	    result.Map(N.Mul(2)),
//	)
//	// result is Ok(84)
//
// Example with error:
//
//	logger := result.Logger[User]()
//	result := F.Pipe2(
//	    result.Error[User](errors.New("database connection failed")),
//	    logger("Fetching user"),  // Logs: "Fetching user: database connection failed"
//	    result.Map(processUser),
//	)
//	// result is Err(error), Map is not executed
//
// Example with custom logger:
//
//	customLogger := log.New(os.Stderr, "APP: ", log.LstdFlags)
//	logger := result.Logger[Data](customLogger)
//
//	result := F.Pipe3(
//	    fetchData(id),
//	    logger("Fetched"),      // Logs to custom logger
//	    result.Map(transform),
//	    logger("Transformed"),  // Logs to custom logger
//	)
//
// Example in a pipeline with multiple logging points:
//
//	logger := result.Logger[Response]()
//
//	result := F.Pipe4(
//	    validateInput(input),
//	    logger("Validated"),
//	    result.Chain(processData),
//	    logger("Processed"),
//	    result.Chain(saveToDatabase),
//	    logger("Saved"),
//	)
//	// Logs at each step, showing the progression or where an error occurred
//
// Use Cases:
//   - Debugging: Track values flowing through Result pipelines
//   - Monitoring: Log successful operations and errors for observability
//   - Auditing: Record operations without affecting the computation
//   - Development: Inspect intermediate values during development
//   - Error tracking: Log errors as they occur in the pipeline
//
// Note: The logging is a side effect that doesn't modify the Result. The original
// Result is always returned, making this function safe to insert anywhere in a
// Result pipeline without changing the computation's semantics.
//
//go:inline
func Logger[A any](loggers ...*log.Logger) func(string) Operator[A, A] {
	return either.Logger[error, A](loggers...)
}

// ToSLogAttr converts a Result value to a structured logging attribute (slog.Attr).
//
// This function creates a converter that transforms Result values into slog.Attr for use
// with Go's structured logging (log/slog). It maps:
//   - Err(error) values to an "error" attribute
//   - Ok(value) values to a "value" attribute
//
// This is a specialized version of either.ToSLogAttr where the error type is fixed to error,
// making it particularly convenient for integrating Result-based error handling with
// structured logging systems. It allows you to log both successful values and errors in a
// consistent, structured format.
//
// Type Parameters:
//   - A: The type of the success value (Right side of the Result)
//
// Returns:
//   - A function that converts Result[A] to slog.Attr
//
// Example with error:
//
//	converter := result.ToSLogAttr[int]()
//	errResult := result.Error[int](errors.New("connection failed"))
//	attr := converter(errResult)
//	// attr is: slog.Any("error", errors.New("connection failed"))
//
//	logger.LogAttrs(ctx, slog.LevelError, "Operation failed", attr)
//	// Logs: {"level":"error","msg":"Operation failed","error":"connection failed"}
//
// Example with success value:
//
//	converter := result.ToSLogAttr[User]()
//	okResult := result.Of(User{ID: 123, Name: "Alice"})
//	attr := converter(okResult)
//	// attr is: slog.Any("value", User{ID: 123, Name: "Alice"})
//
//	logger.LogAttrs(ctx, slog.LevelInfo, "User fetched", attr)
//	// Logs: {"level":"info","msg":"User fetched","value":{"ID":123,"Name":"Alice"}}
//
// Example in a pipeline with structured logging:
//
//	toAttr := result.ToSLogAttr[Data]()
//
//	res := F.Pipe2(
//	    fetchData(id),
//	    result.Map(processData),
//	    result.Map(validateData),
//	)
//
//	attr := toAttr(res)
//	logger.LogAttrs(ctx, slog.LevelInfo, "Data processing complete", attr)
//	// Logs success: {"level":"info","msg":"Data processing complete","value":{...}}
//	// Or error: {"level":"info","msg":"Data processing complete","error":"validation failed"}
//
// Example with custom log levels based on Result:
//
//	toAttr := result.ToSLogAttr[Response]()
//	res := callAPI(endpoint)
//
//	level := result.Fold(
//	    func(error) slog.Level { return slog.LevelError },
//	    func(Response) slog.Level { return slog.LevelInfo },
//	)(res)
//
//	logger.LogAttrs(ctx, level, "API call completed", toAttr(res))
//
// Example with multiple attributes:
//
//	toAttr := result.ToSLogAttr[Order]()
//	res := processOrder(orderID)
//
//	logger.LogAttrs(ctx, slog.LevelInfo, "Order processed",
//	    slog.String("order_id", orderID),
//	    slog.String("user_id", userID),
//	    toAttr(res),  // Adds either "error" or "value" attribute
//	)
//
// Use Cases:
//   - Structured logging: Convert Result outcomes to structured log attributes
//   - Error tracking: Log errors with consistent "error" key in structured logs
//   - Success monitoring: Log successful values with consistent "value" key
//   - Observability: Integrate Result-based error handling with logging systems
//   - Debugging: Inspect Result values in logs with proper structure
//   - Metrics: Extract Result values for metrics collection in logging pipelines
//   - Audit trails: Create structured audit logs from Result computations
//
// Note: The returned slog.Attr uses "error" for Err values and "value" for Ok values.
// These keys are consistent with common structured logging conventions and make it easy
// to query and filter logs based on success or failure.
//
//go:inline
func ToSLogAttr[A any]() func(Result[A]) slog.Attr {
	return either.ToSLogAttr[error, A]()
}
