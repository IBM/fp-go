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

// Package logging provides utilities for creating logging callbacks from standard log.Logger instances.
// It offers a convenient way to configure logging for functional programming patterns where separate
// loggers for success and error cases are needed.
package logging

import (
	"context"
	"log"
	"log/slog"
	"sync/atomic"
)

// LoggingCallbacks creates a pair of logging callback functions from the provided loggers.
// It returns two functions that can be used for logging messages, typically one for success
// cases and one for error cases.
//
// The behavior depends on the number of loggers provided:
//   - 0 loggers: Returns two callbacks using log.Default() for both success and error logging
//   - 1 logger: Returns two callbacks both using the provided logger
//   - 2+ loggers: Returns callbacks using the first logger for success and second for errors
//
// Parameters:
//   - loggers: Variable number of *log.Logger instances (0, 1, or more)
//
// Returns:
//   - First function: Callback for success/info logging (signature: func(string, ...any))
//   - Second function: Callback for error logging (signature: func(string, ...any))
//
// Example:
//
//	// Using default logger for both
//	infoLog, errLog := LoggingCallbacks()
//
//	// Using custom logger for both
//	customLogger := log.New(os.Stdout, "APP: ", log.LstdFlags)
//	infoLog, errLog := LoggingCallbacks(customLogger)
//
//	// Using separate loggers for info and errors
//	infoLogger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
//	errorLogger := log.New(os.Stderr, "ERROR: ", log.LstdFlags)
//	infoLog, errLog := LoggingCallbacks(infoLogger, errorLogger)
func LoggingCallbacks(loggers ...*log.Logger) (func(string, ...any), func(string, ...any)) {
	switch len(loggers) {
	case 0:
		def := log.Default()
		return def.Printf, def.Printf
	case 1:
		log0 := loggers[0]
		return log0.Printf, log0.Printf
	default:
		return loggers[0].Printf, loggers[1].Printf
	}
}

var globalLogger atomic.Pointer[slog.Logger]

func init() {
	globalLogger.Store(slog.Default())
}

// SetLogger sets the global logger instance and returns the previous logger.
// This function is useful for configuring application-wide logging behavior.
//
// Parameters:
//   - l: The new *slog.Logger to set as the global logger
//
// Returns:
//   - The previous *slog.Logger that was set as the global logger
//
// Example:
//
//	oldLogger := SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
//	defer SetLogger(oldLogger) // Restore previous logger
func SetLogger(l *slog.Logger) *slog.Logger {
	return globalLogger.Swap(l)
}

// GetLogger returns the current global logger instance.
// If no logger has been set via SetLogger, it returns slog.Default().
//
// Returns:
//   - The current global *slog.Logger instance
//
// Example:
//
//	logger := GetLogger()
//	logger.Info("Application started")
func GetLogger() *slog.Logger {
	return globalLogger.Load()
}

type loggerInContextType int

var loggerInContextKey loggerInContextType

// GetLoggerFromContext retrieves a logger from the provided context.
// If no logger is found in the context, it returns the global logger.
//
// This function is useful in applications where different parts of the code
// need access to context-specific loggers, such as in request handlers where
// each request might have its own logger with specific attributes.
//
// Parameters:
//   - ctx: The context.Context from which to retrieve the logger
//
// Returns:
//   - A *slog.Logger instance, either from the context or the global logger
//
// Example:
//
//	func handleRequest(ctx context.Context) {
//	    logger := GetLoggerFromContext(ctx)
//	    logger.Info("Processing request")
//	}
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	value, ok := ctx.Value(loggerInContextKey).(*slog.Logger)
	if !ok {
		return globalLogger.Load()
	}
	return value
}

// WithLogger returns an endomorphism that adds a logger to a context.
// An endomorphism is a function that takes a value and returns a value of the same type.
// This function creates a context transformation that embeds the provided logger.
//
// This is particularly useful in functional programming patterns where you want to
// compose context transformations, or when working with middleware that needs to
// inject loggers into request contexts.
//
// Parameters:
//   - l: The *slog.Logger to embed in the context
//
// Returns:
//   - An Endomorphism[context.Context] function that adds the logger to a context
//
// Example:
//
//	// Create a logger transformation
//	addLogger := WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
//
//	// Apply it to a context
//	ctx := context.Background()
//	ctxWithLogger := addLogger(ctx)
//
//	// Retrieve the logger later
//	logger := GetLoggerFromContext(ctxWithLogger)
//	logger.Info("Using context logger")
func WithLogger(l *slog.Logger) Endomorphism[context.Context] {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, loggerInContextKey, l)
	}
}
