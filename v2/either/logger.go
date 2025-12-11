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

package either

import (
	"log"
	"log/slog"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/logging"
)

var (
	// slogError creates a slog.Attr with key "error" for logging error values
	slogError = F.Bind1st(slog.Any, "error")
	// slogValue creates a slog.Attr with key "value" for logging success values
	slogValue = F.Bind1st(slog.Any, "value")
)

func _log[E, A any](left func(string, ...any), right func(string, ...any), prefix string) Operator[E, A, A] {
	return Fold(
		func(e E) Either[E, A] {
			left("%s: %v", prefix, e)
			return Left[A](e)
		},
		func(a A) Either[E, A] {
			right("%s: %v", prefix, a)
			return Right[E](a)
		})
}

// Logger creates a logging function for Either values that logs both Left and Right cases.
// The function logs the value and then returns the original Either unchanged.
//
// Parameters:
//   - loggers: Optional log.Logger instances. If none provided, uses default logger.
//
// Example:
//
//	logger := either.Logger[error, int]()
//	result := F.Pipe2(
//	    either.Right[error](42),
//	    logger("Processing"),
//	    either.Map(N.Mul(2)),
//	)
//	// Logs: "Processing: 42"
//	// result is Right(84)
func Logger[E, A any](loggers ...*log.Logger) func(string) Operator[E, A, A] {
	left, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) Operator[E, A, A] {
		delegate := _log[E, A](left, right, prefix)
		return func(ma Either[E, A]) Either[E, A] {
			return F.Pipe1(
				delegate(ma),
				ChainTo[A](ma),
			)
		}
	}
}

// ToSLogAttr converts an Either value to a structured logging attribute (slog.Attr).
//
// This function creates a converter that transforms Either values into slog.Attr for use
// with Go's structured logging (log/slog). It maps:
//   - Left values to an "error" attribute
//   - Right values to a "value" attribute
//
// This is particularly useful when integrating Either-based error handling with structured
// logging systems, allowing you to log both successful values and errors in a consistent,
// structured format.
//
// Type Parameters:
//   - E: The Left (error) type of the Either
//   - A: The Right (success) type of the Either
//
// Returns:
//   - A function that converts Either[E, A] to slog.Attr
//
// Example with Left (error):
//
//	converter := either.ToSLogAttr[error, int]()
//	leftValue := either.Left[int](errors.New("connection failed"))
//	attr := converter(leftValue)
//	// attr is: slog.Any("error", errors.New("connection failed"))
//
//	logger.LogAttrs(ctx, slog.LevelError, "Operation failed", attr)
//	// Logs: {"level":"error","msg":"Operation failed","error":"connection failed"}
//
// Example with Right (success):
//
//	converter := either.ToSLogAttr[error, User]()
//	rightValue := either.Right[error](User{ID: 123, Name: "Alice"})
//	attr := converter(rightValue)
//	// attr is: slog.Any("value", User{ID: 123, Name: "Alice"})
//
//	logger.LogAttrs(ctx, slog.LevelInfo, "User fetched", attr)
//	// Logs: {"level":"info","msg":"User fetched","value":{"ID":123,"Name":"Alice"}}
//
// Example in a pipeline with structured logging:
//
//	toAttr := either.ToSLogAttr[error, Data]()
//
//	result := F.Pipe2(
//	    fetchData(id),
//	    either.Map(processData),
//	    either.Map(validateData),
//	)
//
//	attr := toAttr(result)
//	logger.LogAttrs(ctx, slog.LevelInfo, "Data processing complete", attr)
//	// Logs success: {"level":"info","msg":"Data processing complete","value":{...}}
//	// Or error: {"level":"info","msg":"Data processing complete","error":"validation failed"}
//
// Example with custom log levels based on Either:
//
//	toAttr := either.ToSLogAttr[error, Response]()
//	result := callAPI(endpoint)
//
//	level := either.Fold(
//	    func(error) slog.Level { return slog.LevelError },
//	    func(Response) slog.Level { return slog.LevelInfo },
//	)(result)
//
//	logger.LogAttrs(ctx, level, "API call completed", toAttr(result))
//
// Use Cases:
//   - Structured logging: Convert Either results to structured log attributes
//   - Error tracking: Log errors with consistent "error" key in structured logs
//   - Success monitoring: Log successful values with consistent "value" key
//   - Observability: Integrate Either-based error handling with logging systems
//   - Debugging: Inspect Either values in logs with proper structure
//   - Metrics: Extract Either values for metrics collection in logging pipelines
//
// Note: The returned slog.Attr uses "error" for Left values and "value" for Right values.
// These keys are consistent with common structured logging conventions.
func ToSLogAttr[E, A any]() func(Either[E, A]) slog.Attr {
	return Fold(
		F.Flow2(
			F.ToAny[E],
			slogError,
		),
		F.Flow2(
			F.ToAny[A],
			slogValue,
		),
	)
}
