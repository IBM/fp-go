// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"log"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/logging"
)

func _log[A any](left, right func(string, ...any), prefix string) Kleisli[Option[A], A] {
	return Fold(
		func() Option[A] {
			left("%s", prefix)
			return None[A]()
		},
		func(a A) Option[A] {
			right("%s: %v", prefix, a)
			return Some(a)
		})
}

// Logger creates a logging function for Options that logs the state (None or Some with value)
// and returns the original Option unchanged. This is useful for debugging pipelines.
//
// Parameters:
//   - loggers: optional log.Logger instances to use for logging (defaults to standard logger)
//
// Returns a function that takes a prefix string and returns a function that logs and passes through an Option.
//
// Example:
//
//	logger := Logger[int]()
//	result := F.Pipe2(
//	    Some(42),
//	    logger("step1"), // logs "step1: 42"
//	    Map(N.Mul(2)),
//	) // Some(84)
//
//	result := F.Pipe1(
//	    None[int](),
//	    logger("step1"), // logs "step1"
//	) // None
func Logger[A any](loggers ...*log.Logger) func(string) Kleisli[Option[A], A] {
	left, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) Kleisli[Option[A], A] {
		delegate := _log[A](left, right, prefix)
		return func(ma Option[A]) Option[A] {
			return F.Pipe1(
				delegate(ma),
				ChainTo[A](ma),
			)
		}
	}
}
