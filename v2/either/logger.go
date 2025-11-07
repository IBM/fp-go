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

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/logging"
)

func _log[E, A any](left func(string, ...any), right func(string, ...any), prefix string) Operator[E, A, A] {
	return Fold(
		func(e E) Either[E, A] {
			left("%s: %v", prefix, e)
			return Left[A, E](e)
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
//	    either.Map(func(x int) int { return x * 2 }),
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
