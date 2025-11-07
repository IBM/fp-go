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

package io

import (
	"fmt"
	"log"

	L "github.com/IBM/fp-go/v2/logging"
)

// Logger constructs a logger function that can be used with ChainFirst or similar operations.
// It logs values using the provided loggers (or the default logger if none provided).
//
// Example:
//
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    io.ChainFirst(io.Logger[User]()("Fetched user")),
//	    processUser,
//	)
func Logger[A any](loggers ...*log.Logger) func(string) Kleisli[A, any] {
	_, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) Kleisli[A, any] {
		return func(a A) IO[any] {
			return FromImpure(func() {
				right("%s: %v", prefix, a)
			})
		}
	}
}

// Logf constructs a logger function that can be used with ChainFirst or similar operations.
// The prefix string contains the format string for the log value.
//
// Example:
//
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    io.ChainFirst(io.Logf[User]("User: %+v")),
//	    processUser,
//	)
func Logf[A any](prefix string) Kleisli[A, any] {
	return func(a A) IO[any] {
		return FromImpure(func() {
			log.Printf(prefix, a)
		})
	}
}

// Printf constructs a printer function that can be used with ChainFirst or similar operations.
// The prefix string contains the format string for the printed value.
// Unlike Logf, this prints to stdout without log prefixes.
//
// Example:
//
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    io.ChainFirst(io.Printf[User]("User: %+v\n")),
//	    processUser,
//	)
func Printf[A any](prefix string) Kleisli[A, any] {
	return func(a A) IO[any] {
		return FromImpure(func() {
			fmt.Printf(prefix, a)
		})
	}
}
