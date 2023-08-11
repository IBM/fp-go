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

package errors

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
)

// OnNone generates a nullary function that produces a formatted error
func OnNone(msg string, args ...any) func() error {
	return func() error {
		return fmt.Errorf(msg, args...)
	}
}

// OnSome generates a unary function that produces a formatted error
func OnSome[T any](msg string, args ...any) func(T) error {
	l := len(args)
	if l == 0 {
		return func(value T) error {
			return fmt.Errorf(msg, value)
		}
	}
	return func(value T) error {
		data := make([]any, l)
		copy(data[1:], args)
		data[0] = value
		return fmt.Errorf(msg, data...)
	}
}

// OnError generates a unary function that produces a formatted error. The argument
// to that function is the root cause of the error and the message will be augmented with
// a format string containing %w
func OnError(msg string, args ...any) func(error) error {
	return func(err error) error {
		return fmt.Errorf(msg+", Caused By: %w", A.ArrayConcatAll(args, A.Of[any](err))...)
	}
}

// ToString converts an error to a string
func ToString(err error) string {
	return err.Error()
}
