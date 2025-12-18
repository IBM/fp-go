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

package option_test

import (
	"fmt"
	"log/slog"
	"os"

	O "github.com/IBM/fp-go/v2/option"
)

// ExampleOption_String demonstrates the fmt.Stringer interface implementation.
func ExampleOption_String() {
	some := O.Some(42)
	none := O.None[int]()

	fmt.Println(some.String())
	fmt.Println(none.String())

	// Output:
	// Some[int](42)
	// None[int]
}

// ExampleOption_GoString demonstrates the fmt.GoStringer interface implementation.
func ExampleOption_GoString() {
	some := O.Some(42)
	none := O.None[int]()

	fmt.Printf("%#v\n", some)
	fmt.Printf("%#v\n", none)

	// Output:
	// option.Some[int](42)
	// option.None[int]
}

// ExampleOption_Format demonstrates the fmt.Formatter interface implementation.
func ExampleOption_Format() {
	result := O.Some(42)

	// Different format verbs
	fmt.Printf("%%s: %s\n", result)
	fmt.Printf("%%v: %v\n", result)
	fmt.Printf("%%+v: %+v\n", result)
	fmt.Printf("%%#v: %#v\n", result)

	// Output:
	// %s: Some[int](42)
	// %v: Some[int](42)
	// %+v: Some[int](42)
	// %#v: option.Some[int](42)
}

// ExampleOption_LogValue demonstrates the slog.LogValuer interface implementation.
func ExampleOption_LogValue() {
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

	// Some value
	someResult := O.Some(42)
	logger.Info("computation succeeded", "result", someResult)

	// None value
	noneResult := O.None[int]()
	logger.Info("computation failed", "result", noneResult)

	// Output:
	// level=INFO msg="computation succeeded" result.some=42
	// level=INFO msg="computation failed" result.none={}
}

// ExampleOption_formatting_comparison demonstrates different formatting options.
func ExampleOption_formatting_comparison() {
	type User struct {
		ID   int
		Name string
	}

	user := User{ID: 123, Name: "Alice"}
	result := O.Some(user)

	fmt.Printf("String():   %s\n", result.String())
	fmt.Printf("GoString(): %s\n", result.GoString())
	fmt.Printf("%%v:         %v\n", result)
	fmt.Printf("%%#v:        %#v\n", result)

	// Output:
	// String():   Some[option_test.User]({123 Alice})
	// GoString(): option.Some[option_test.User](option_test.User{ID:123, Name:"Alice"})
	// %v:         Some[option_test.User]({123 Alice})
	// %#v:        option.Some[option_test.User](option_test.User{ID:123, Name:"Alice"})
}

// ExampleOption_LogValue_structured demonstrates structured logging with Option.
func ExampleOption_LogValue_structured() {
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
	compute := func(x int) O.Option[int] {
		if x < 0 {
			return O.None[int]()
		}
		return O.Some(x * 2)
	}

	// Log successful computation
	result1 := compute(21)
	logger.Info("computation", "input", 21, "output", result1)

	// Log failed computation
	result2 := compute(-5)
	logger.Warn("computation", "input", -5, "output", result2)

	// Output:
	// level=INFO msg=computation input=21 output.some=42
	// level=WARN msg=computation input=-5 output.none={}
}

// Example_none_formatting demonstrates formatting of None values.
func Example_none_formatting() {
	none := O.None[string]()

	fmt.Printf("String():   %s\n", none.String())
	fmt.Printf("GoString(): %s\n", none.GoString())
	fmt.Printf("%%v:         %v\n", none)
	fmt.Printf("%%#v:        %#v\n", none)

	// Output:
	// String():   None[string]
	// GoString(): option.None[string]
	// %v:         None[string]
	// %#v:        option.None[string]
}
