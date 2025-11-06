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
	"errors"

	O "github.com/IBM/fp-go/v2/option"
)

// As tries to extract an error of the desired type from the given error.
// It returns an Option containing the extracted error if successful, or None if the
// error cannot be converted to the target type.
//
// This function wraps Go's standard errors.As in a functional style, making it
// composable with other functional operations.
//
// Example:
//
//	type MyError struct{ msg string }
//	func (e *MyError) Error() string { return e.msg }
//
//	rootErr := &MyError{msg: "custom error"}
//	wrappedErr := fmt.Errorf("wrapped: %w", rootErr)
//
//	// Extract MyError from the wrapped error
//	extractMyError := As[*MyError]()
//	result := extractMyError(wrappedErr)
//	// result is Some(*MyError) containing the original error
//
//	// Try to extract a different error type
//	extractOther := As[*os.PathError]()
//	result2 := extractOther(wrappedErr)
//	// result2 is None since wrappedErr doesn't contain *os.PathError
func As[A error]() func(error) O.Option[A] {
	return O.FromValidation(func(err error) (A, bool) {
		var a A
		ok := errors.As(err, &a)
		return a, ok
	})
}
