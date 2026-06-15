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

package either_test

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	E "github.com/IBM/fp-go/v2/either"
)

// ExampleEither_String demonstrates the fmt.Stringer interface implementation.
func ExampleEither_String() {
	right := E.Right[error](42)
	left := E.Left[int](errors.New("something went wrong"))

	fmt.Println(right.String())
	fmt.Println(left.String())

	// Output:
	// Right[int](42)
	// Left[*errors.errorString](something went wrong)
}

// ExampleEither_GoString demonstrates the fmt.GoStringer interface implementation.
func ExampleEither_GoString() {
	right := E.Right[error](42)
	left := E.Left[int](errors.New("error"))

	fmt.Printf("%#v\n", right)
	fmt.Printf("%#v\n", left)

	// Output:
	// either.Right[error](42)
	// either.Left[int](&errors.errorString{s:"error"})
}

// ExampleEither_Format demonstrates the fmt.Formatter interface implementation.
func ExampleEither_Format() {
	result := E.Right[error](42)

	// Different format verbs
	fmt.Printf("%%s: %s\n", result)
	fmt.Printf("%%v: %v\n", result)
	fmt.Printf("%%+v: %+v\n", result)
	fmt.Printf("%%#v: %#v\n", result)

	// Output:
	// %s: Right[int](42)
	// %v: Right[int](42)
	// %+v: Right[int](42)
	// %#v: either.Right[error](42)
}

// ExampleEither_LogValue demonstrates the slog.LogValuer interface implementation.
func ExampleEither_LogValue() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time for consistent output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// Right value
	rightResult := E.Right[error](42)
	logger.Info("computation succeeded", "result", rightResult)

	// Left value
	leftResult := E.Left[int](errors.New("computation failed"))
	logger.Error("computation failed", "result", leftResult)

	// Output:
	// level=INFO msg="computation succeeded" result.right=42
	// level=ERROR msg="computation failed" result.left="computation failed"
}

// ExampleEither_formatting_comparison demonstrates different formatting options.
func ExampleEither_formatting_comparison() {
	type User struct {
		ID   int
		Name string
	}

	user := User{ID: 123, Name: "Alice"}
	result := E.Right[error](user)

	fmt.Printf("String():   %s\n", result.String())
	fmt.Printf("GoString(): %s\n", result.GoString())
	fmt.Printf("%%v:         %v\n", result)
	fmt.Printf("%%#v:        %#v\n", result)

	// Output:
	// String():   Right[either_test.User]({123 Alice})
	// GoString(): either.Right[error](either_test.User{ID:123, Name:"Alice"})
	// %v:         Right[either_test.User]({123 Alice})
	// %#v:        either.Right[error](either_test.User{ID:123, Name:"Alice"})
}

// ExampleEither_LogValue_structured demonstrates structured logging with Either.
func ExampleEither_LogValue_structured() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// Simulate a computation pipeline
	compute := func(x int) E.Either[error, int] {
		if x < 0 {
			return E.Left[int](errors.New("negative input"))
		}
		return E.Right[error](x * 2)
	}

	// Log successful computation
	result1 := compute(21)
	logger.Info("computation", "input", 21, "output", result1)

	// Log failed computation
	result2 := compute(-5)
	logger.Error("computation", "input", -5, "output", result2)

	// Output:
	// level=INFO msg=computation input=21 output.right=42
	// level=ERROR msg=computation input=-5 output.left="negative input"
}
