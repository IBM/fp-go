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

package errors

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/endomorphism"
)

// OnNone generates a nullary function that produces a formatted error.
// This is useful when you need to create an error lazily, such as when
// handling the None case in an Option type.
//
// Example:
//
//	getError := OnNone("value not found")
//	err := getError() // returns error: "value not found"
//
//	getErrorWithArgs := OnNone("failed to load %s", "config.json")
//	err2 := getErrorWithArgs() // returns error: "failed to load config.json"
func OnNone(msg string, args ...any) func() error {
	return func() error {
		return fmt.Errorf(msg, args...)
	}
}

// OnSome generates a unary function that produces a formatted error.
// The function takes a value of type T and includes it in the error message.
// If no additional args are provided, the value is used as the only format argument.
// If additional args are provided, the value becomes the first format argument.
//
// This is useful when you need to create an error that includes information
// about a value, such as when handling the Some case in an Option type.
//
// Example:
//
//	// Without additional args - value is the only format argument
//	makeError := OnSome[int]("invalid value: %d")
//	err := makeError(42) // returns error: "invalid value: 42"
//
//	// With additional args - value is the first format argument
//	makeError2 := OnSome[string]("failed to process %s in file %s", "data.txt")
//	err2 := makeError2("record123") // returns error: "failed to process record123 in file data.txt"
func OnSome[T any](msg string, args ...any) func(T) error {
	l := len(args)
	if l == 0 {
		return func(value T) error {
			return fmt.Errorf(msg, value)
		}
	}
	return func(value T) error {
		data := make([]any, l+1)
		data[0] = value
		copy(data[1:], args)
		return fmt.Errorf(msg, data...)
	}
}

// OnError generates a unary function that produces a formatted error with error wrapping.
// The argument to that function is the root cause of the error and the message will be
// augmented with a format string containing %w for proper error wrapping.
//
// This is useful for adding context to errors while preserving the error chain,
// allowing errors.Is and errors.As to work correctly.
//
// Example:
//
//	wrapError := OnError("failed to load configuration from %s", "config.json")
//	rootErr := errors.New("file not found")
//	wrappedErr := wrapError(rootErr)
//	// returns error: "failed to load configuration from config.json, Caused By: file not found"
//	// errors.Is(wrappedErr, rootErr) returns true
func OnError(msg string, args ...any) endomorphism.Endomorphism[error] {
	return func(err error) error {
		return fmt.Errorf(msg+", Caused By: %w", A.ArrayConcatAll(args, A.Of[any](err))...)
	}
}

// ToString converts an error to its string representation by calling the Error() method.
// This is useful in functional pipelines where you need to transform an error into a string.
//
// Example:
//
//	err := errors.New("something went wrong")
//	msg := ToString(err) // returns "something went wrong"
func ToString(err error) string {
	return err.Error()
}
