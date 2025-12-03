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
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/IBM/fp-go/v2/function"
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
func Logger[A any](loggers ...*log.Logger) func(string) Kleisli[A, A] {
	_, right := L.LoggingCallbacks(loggers...)
	return func(prefix string) Kleisli[A, A] {
		return func(a A) IO[A] {
			return func() A {
				right("%s: %v", prefix, a)
				return a
			}
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
func Logf[A any](prefix string) Kleisli[A, A] {
	return func(a A) IO[A] {
		return func() A {
			log.Printf(prefix, a)
			return a
		}
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
func Printf[A any](prefix string) Kleisli[A, A] {
	return func(a A) IO[A] {
		return func() A {
			fmt.Printf(prefix, a)
			return a
		}
	}
}

// handleLogging is a helper function that creates a Kleisli arrow for logging/printing
// values using Go template syntax. It lazily compiles the template on first use and
// executes it with the provided value as data.
//
// Parameters:
//   - onSuccess: callback function to handle successfully formatted output
//   - onError: callback function to handle template parsing or execution errors
//   - prefix: Go template string to format the value
//
// The template is compiled lazily using sync.Once to ensure it's only parsed once.
// The function always returns the original value unchanged, making it suitable for
// use with ChainFirst or similar operations.
func handleLoggingG(onSuccess func(string), onError func(error), prefix string) Kleisli[any, any] {
	var tmp *template.Template
	var err error
	var once sync.Once

	init := func() {
		tmp, err = template.New("").Parse(prefix)
	}

	return func(a any) IO[any] {
		return func() any {
			// make sure to compile lazily
			once.Do(init)
			if err == nil {
				var buffer strings.Builder
				tmpErr := tmp.Execute(&buffer, a)
				if tmpErr != nil {
					onError(tmpErr)
					onSuccess(fmt.Sprintf("%v", a))
				} else {
					onSuccess(buffer.String())
				}
			} else {
				onError(err)
				onSuccess(fmt.Sprintf("%v", a))
			}
			// in any case return the original value
			return a
		}
	}
}

// handleLogging is a helper function that creates a Kleisli arrow for logging/printing
// values using Go template syntax. It lazily compiles the template on first use and
// executes it with the provided value as data.
//
// Parameters:
//   - onSuccess: callback function to handle successfully formatted output
//   - onError: callback function to handle template parsing or execution errors
//   - prefix: Go template string to format the value
//
// The template is compiled lazily using sync.Once to ensure it's only parsed once.
// The function always returns the original value unchanged, making it suitable for
// use with ChainFirst or similar operations.
func handleLogging[A any](onSuccess func(string), onError func(error), prefix string) Kleisli[A, A] {
	generic := handleLoggingG(onSuccess, onError, prefix)
	return func(a A) IO[A] {
		return function.Pipe1(
			generic(a),
			MapTo[any](a),
		)
	}
}

// LogGo constructs a logger function using Go template syntax for formatting.
// The prefix string is parsed as a Go template and executed with the value as data.
// Both successful output and template errors are logged using log.Println.
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    io.ChainFirst(io.LogGo[User]("User: {{.Name}}, Age: {{.Age}}")),
//	    processUser,
//	)
func LogGo[A any](prefix string) Kleisli[A, A] {
	return handleLogging[A](func(value string) {
		log.Println(value)
	}, func(err error) {
		log.Println(err)
	}, prefix)
}

// PrintGo constructs a printer function using Go template syntax for formatting.
// The prefix string is parsed as a Go template and executed with the value as data.
// Successful output is printed to stdout using fmt.Println, while template errors
// are printed to stderr using fmt.Fprintln.
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    io.ChainFirst(io.PrintGo[User]("User: {{.Name}}, Age: {{.Age}}")),
//	    processUser,
//	)
func PrintGo[A any](prefix string) Kleisli[A, A] {
	return handleLogging[A](func(value string) {
		fmt.Println(value)
	}, func(err error) {
		fmt.Fprintln(os.Stderr, err)
	}, prefix)
}
