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

// Package readerresult provides logging utilities for the ReaderResult monad,
// which combines the Reader monad (for dependency injection via context.Context)
// with the Result monad (for error handling).
//
// The logging functions in this package allow you to log Result values (both
// successes and errors) while preserving the functional composition style.
package readerresult

import (
	"context"
	"log/slog"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/logging"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

// curriedLog creates a curried logging function that takes an slog.Attr and a context,
// then logs the attribute with the specified log level and message.
//
// This is an internal helper function used to create the logging pipeline in a
// point-free style. The currying allows for partial application in functional
// composition.
//
// Parameters:
//   - logLevel: The slog.Level at which to log (e.g., LevelInfo, LevelError)
//   - cb: A callback function that retrieves a logger from the context
//   - message: The log message to display
//
// Returns:
//   - A curried function that takes an slog.Attr, then a context, and performs logging
func curriedLog(
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	message string) func(slog.Attr) Reader[context.Context, struct{}] {
	return F.Curry2(func(a slog.Attr, ctx context.Context) struct{} {
		cb(ctx).LogAttrs(ctx, logLevel, message, a)
		return struct{}{}
	})
}

// SLogWithCallback creates a Kleisli arrow that logs a Result value using a custom
// logger callback and log level. The Result value is logged and then returned unchanged,
// making this function suitable for use in functional pipelines.
//
// This function logs both successful values and errors:
//   - Success values are logged with the key "value"
//   - Error values are logged with the key "error"
//
// The logging is performed as a side effect while preserving the Result value,
// allowing it to be used in the middle of a computation pipeline without
// interrupting the flow.
//
// Type Parameters:
//   - A: The type of the success value in the Result
//
// Parameters:
//   - logLevel: The slog.Level at which to log (e.g., LevelInfo, LevelDebug, LevelError)
//   - cb: A callback function that retrieves a *slog.Logger from the context
//   - message: The log message to display
//
// Returns:
//   - A Kleisli arrow that takes a Result[A] and returns a ReaderResult[A]
//     The returned ReaderResult, when executed with a context, logs the Result
//     and returns it unchanged
//
// Example:
//
//	type User struct {
//	    ID   int
//	    Name string
//	}
//
//	// Custom logger callback
//	getLogger := func(ctx context.Context) *slog.Logger {
//	    return slog.Default()
//	}
//
//	// Create a logging function for debug level
//	logDebug := SLogWithCallback[User](slog.LevelDebug, getLogger, "User data")
//
//	// Use in a pipeline
//	ctx := t.Context()
//	user := result.Of(User{ID: 123, Name: "Alice"})
//	logged := logDebug(user)(ctx) // Logs: level=DEBUG msg="User data" value={ID:123 Name:Alice}
//	// logged still contains the User value
//
// Example with error:
//
//	err := errors.New("user not found")
//	userResult := result.Left[User](err)
//	logged := logDebug(userResult)(ctx) // Logs: level=DEBUG msg="User data" error="user not found"
//	// logged still contains the error
func SLogWithCallback[A any](
	logLevel slog.Level,
	cb func(context.Context) *slog.Logger,
	message string) Kleisli[Result[A], A] {

	return F.Pipe1(
		F.Flow2(
			result.ToSLogAttr[A](),
			curriedLog(logLevel, cb, message),
		),
		reader.Chain(reader.Sequence(F.Flow2( // this flow is basically the `MapTo` function with side effects
			reader.Of[struct{}, Result[A]],
			reader.Map[context.Context, struct{}, Result[A]],
		))),
	)

}

// SLog creates a Kleisli arrow that logs a Result value at INFO level using the
// logger from the context. This is a convenience function that uses SLogWithCallback
// with default settings.
//
// The Result value is logged and then returned unchanged, making this function
// suitable for use in functional pipelines for debugging or monitoring purposes.
//
// This function logs both successful values and errors:
//   - Success values are logged with the key "value"
//   - Error values are logged with the key "error"
//
// Type Parameters:
//   - A: The type of the success value in the Result
//
// Parameters:
//   - message: The log message to display
//
// Returns:
//   - A Kleisli arrow that takes a Result[A] and returns a ReaderResult[A]
//     The returned ReaderResult, when executed with a context, logs the Result
//     at INFO level and returns it unchanged
//
// Example - Logging a successful computation:
//
//	ctx := t.Context()
//
//	// Simple value logging
//	res := result.Of(42)
//	logged := SLog[int]("Processing number")(res)(ctx)
//	// Logs: level=INFO msg="Processing number" value=42
//	// logged == result.Of(42)
//
// Example - Logging in a pipeline:
//
//	type User struct {
//	    ID   int
//	    Name string
//	}
//
//	fetchUser := func(id int) result.Result[User] {
//	    return result.Of(User{ID: id, Name: "Alice"})
//	}
//
//	processUser := func(user User) result.Result[string] {
//	    return result.Of(fmt.Sprintf("Processed: %s", user.Name))
//	}
//
//	ctx := t.Context()
//
//	// Log at each step
//	userResult := fetchUser(123)
//	logged1 := SLog[User]("Fetched user")(userResult)(ctx)
//	// Logs: level=INFO msg="Fetched user" value={ID:123 Name:Alice}
//
//	processed := result.Chain(processUser)(logged1)
//	logged2 := SLog[string]("Processed user")(processed)(ctx)
//	// Logs: level=INFO msg="Processed user" value="Processed: Alice"
//
// Example - Logging errors:
//
//	err := errors.New("database connection failed")
//	errResult := result.Left[User](err)
//	logged := SLog[User]("Database operation")(errResult)(ctx)
//	// Logs: level=INFO msg="Database operation" error="database connection failed"
//	// logged still contains the error
//
// Example - Using with context logger:
//
//	// Set up a custom logger in the context
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	ctx := logging.WithLogger(logger)(t.Context())
//
//	res := result.Of("important data")
//	logged := SLog[string]("Critical operation")(res)(ctx)
//	// Uses the logger from context to log the message
//
// Note: The function uses logging.GetLoggerFromContext to retrieve the logger,
// which falls back to the global logger if no logger is found in the context.
//
//go:inline
func SLog[A any](message string) Kleisli[Result[A], A] {
	return SLogWithCallback[A](slog.LevelInfo, logging.GetLoggerFromContext, message)
}

//go:inline
func TapSLog[A any](message string) Operator[A, A] {
	return reader.Chain(SLog[A](message))
}
