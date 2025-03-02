// Copyright (c) 2023 IBM Corp.
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

// Logger constructs a logger function that can be used with ChainXXXIOK
func Logger[A any](loggers ...*log.Logger) func(string) func(A) IO[any] {
	_, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) func(A) IO[any] {
		return func(a A) IO[any] {
			return FromImpure(func() {
				right("%s: %v", prefix, a)
			})
		}
	}
}

// Logf constructs a logger function that can be used with ChainXXXIOK
// the string prefix contains the format string for the log value
func Logf[A any](prefix string) func(A) IO[any] {
	return func(a A) IO[any] {
		return FromImpure(func() {
			log.Printf(prefix, a)
		})
	}
}

// Printf constructs a printer function that can be used with ChainXXXIOK
// the string prefix contains the format string for the log value
func Printf[A any](prefix string) func(A) IO[any] {
	return func(a A) IO[any] {
		return FromImpure(func() {
			fmt.Printf(prefix, a)
		})
	}
}
