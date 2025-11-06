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
	"log"
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
