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

package effect

import (
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/reader"
)

// TapSLog logs the outcome of an effect at Info level and passes it through unchanged.
//
// A success is logged with its value, a failure with its error. The logger is taken from
// the context.Context, falling back to the global logger, see logging.WithLogger and
// logging.SetLogger.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - message: The log message
//
// Returns:
//   - Operator[C, A, A]: A function that adds logging to an effect
//
// See Also:
//   - TapSLogDebug: Logs at Debug level
//   - AnnotateLogs: Adds attributes to all log records of an effect
func TapSLog[C, A any](message string) Operator[C, A, A] {
	return reader.Map[C](thunk.TapSLog[A](message))
}

// TapSLogDebug logs the outcome of an effect at Debug level and passes it through unchanged.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - message: The log message
//
// Returns:
//   - Operator[C, A, A]: A function that adds logging to an effect
//
// See Also:
//   - TapSLog: Logs at Info level
func TapSLogDebug[C, A any](message string) Operator[C, A, A] {
	return reader.Map[C](thunk.TapSLogDebug[A](message))
}

// LogEntryExit logs when an effect starts and when it ends, with its duration and outcome.
//
// Each run gets a correlation ID that appears in both records, so nested and concurrent
// runs can be told apart.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - name: The name of the operation, used in the log records
//
// Returns:
//   - Operator[C, A, A]: A function that adds entry and exit logging to an effect
//
// See Also:
//   - TapSLog: Logs only the outcome
func LogEntryExit[C, A any](name string) Operator[C, A, A] {
	return reader.Map[C](thunk.LogEntryExit[A](name))
}

// AnnotateLogs runs an effect with a context logger that carries additional attributes.
//
// All logging inside the effect that reads the context logger, such as TapSLog and
// LogEntryExit, includes the attributes. Logging outside the effect is unaffected. This
// corresponds to Effect.annotateLogs in effect-ts.
//
// Type Parameters:
//   - C: The dependencies required by the effect
//   - A: The type of the value produced by the effect
//
// Parameters:
//   - args: Alternating keys and values, or slog.Attr values, as accepted by slog.Logger.With
//
// Returns:
//   - Operator[C, A, A]: A function that runs an effect with the annotated logger
//
// See Also:
//   - TapSLog: Logs the outcome of an effect
func AnnotateLogs[C, A any](args ...any) Operator[C, A, A] {
	return reader.Map[C](thunk.AnnotateLogs[A](args...))
}
