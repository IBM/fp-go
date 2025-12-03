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

package readerio

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
)

// Logf constructs a logger function that can be used with ChainFirst or similar operations.
// The prefix string contains the format string for both the reader context (R) and the value (A).
// It uses log.Printf to output the formatted message.
//
// Type Parameters:
//   - R: Reader context type
//   - A: Value type
//
// Parameters:
//   - prefix: Format string that accepts two arguments: the reader context and the value
//
// Returns:
//   - A Kleisli arrow that logs the context and value, then returns the original value
//
// Example:
//
//	type Config struct {
//	    AppName string
//	}
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    readerio.ChainFirst(readerio.Logf[Config, User]("[%v] User: %+v")),
//	    processUser,
//	)(Config{AppName: "MyApp"})()
func Logf[R, A any](prefix string) Kleisli[R, A, A] {
	return func(a A) ReaderIO[R, A] {
		return func(r R) IO[A] {
			return func() A {
				log.Printf(prefix, r, a)
				return a
			}
		}
	}
}

// Printf constructs a printer function that can be used with ChainFirst or similar operations.
// The prefix string contains the format string for both the reader context (R) and the value (A).
// Unlike Logf, this prints to stdout without log prefixes.
//
// Type Parameters:
//   - R: Reader context type
//   - A: Value type
//
// Parameters:
//   - prefix: Format string that accepts two arguments: the reader context and the value
//
// Returns:
//   - A Kleisli arrow that prints the context and value, then returns the original value
//
// Example:
//
//	type Config struct {
//	    Debug bool
//	}
//	result := pipe.Pipe2(
//	    fetchData(),
//	    readerio.ChainFirst(readerio.Printf[Config, Data]("[%v] Data: %+v\n")),
//	    processData,
//	)(Config{Debug: true})()
func Printf[R, A any](prefix string) Kleisli[R, A, A] {
	return func(a A) ReaderIO[R, A] {
		return func(r R) IO[A] {
			return func() A {
				fmt.Printf(prefix, r, a)
				return a
			}
		}
	}
}

// handleLoggingG is a generic helper function that creates a Kleisli arrow for logging/printing
// values using Go template syntax. It lazily compiles the template on first use and
// executes it with a context struct containing both the reader context (R) and value (A).
//
// Parameters:
//   - onSuccess: callback function to handle successfully formatted output
//   - onError: callback function to handle template parsing or execution errors
//   - prefix: Go template string to format the context and value
//
// The template is compiled lazily using sync.Once to ensure it's only parsed once.
// The template receives a context struct with fields R (reader context) and A (value).
// The function always returns the original value unchanged, making it suitable for
// use with ChainFirst or similar operations.
func handleLoggingG(onSuccess func(string), onError func(error), prefix string) Kleisli[any, any, any] {
	var tmp *template.Template
	var err error
	var once sync.Once

	type context struct {
		R any
		A any
	}

	init := func() {
		tmp, err = template.New("").Parse(prefix)
	}
	return func(a any) ReaderIO[any, any] {
		return func(r any) IO[any] {
			return func() any {
				// make sure to compile lazily
				once.Do(init)
				if err == nil {
					var buffer strings.Builder
					tmpErr := tmp.Execute(&buffer, context{r, a})
					if tmpErr != nil {
						onError(tmpErr)
						onSuccess(fmt.Sprintf("%v: %v", r, a))
					} else {
						onSuccess(buffer.String())
					}
				} else {
					onError(err)
					onSuccess(fmt.Sprintf("%v: %v", r, a))
				}
				// in any case return the original value
				return a
			}
		}
	}
}

// handleLogging is a typed wrapper around handleLoggingG that creates a Kleisli arrow
// for logging/printing values using Go template syntax.
//
// Parameters:
//   - onSuccess: callback function to handle successfully formatted output
//   - onError: callback function to handle template parsing or execution errors
//   - prefix: Go template string to format the context and value
//
// Returns:
//   - A Kleisli arrow that formats and outputs the value, then returns it unchanged
func handleLogging[R, A any](onSuccess func(string), onError func(error), prefix string) Kleisli[R, A, A] {
	generic := handleLoggingG(onSuccess, onError, prefix)
	return func(a A) ReaderIO[R, A] {
		ga := generic(a)
		return func(r R) IO[A] {
			gr := ga(r)
			return func() A {
				gr()
				return a
			}
		}
	}
}

// LogGo constructs a logger function using Go template syntax for formatting.
// The prefix string is parsed as a Go template and executed with a context struct
// containing both the reader context (R) and the value (A) as fields .R and .A.
// Both successful output and template errors are logged using log.Println.
//
// Type Parameters:
//   - R: Reader context type
//   - A: Value type
//
// Parameters:
//   - prefix: Go template string with access to .R (context) and .A (value)
//
// Returns:
//   - A Kleisli arrow that logs the formatted output and returns the original value
//
// Example:
//
//	type Config struct {
//	    AppName string
//	}
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	result := pipe.Pipe2(
//	    fetchUser(),
//	    readerio.ChainFirst(readerio.LogGo[Config, User]("[{{.R.AppName}}] User: {{.A.Name}}, Age: {{.A.Age}}")),
//	    processUser,
//	)(Config{AppName: "MyApp"})()
func LogGo[R, A any](prefix string) Kleisli[R, A, A] {
	return handleLogging[R, A](func(value string) {
		log.Println(value)
	}, func(err error) {
		log.Println(err)
	}, prefix)
}

// PrintGo constructs a printer function using Go template syntax for formatting.
// The prefix string is parsed as a Go template and executed with a context struct
// containing both the reader context (R) and the value (A) as fields .R and .A.
// Successful output is printed to stdout using fmt.Println, while template errors
// are printed to stderr using fmt.Fprintln.
//
// Type Parameters:
//   - R: Reader context type
//   - A: Value type
//
// Parameters:
//   - prefix: Go template string with access to .R (context) and .A (value)
//
// Returns:
//   - A Kleisli arrow that prints the formatted output and returns the original value
//
// Example:
//
//	type Config struct {
//	    Verbose bool
//	}
//	type Data struct {
//	    ID    int
//	    Value string
//	}
//	result := pipe.Pipe2(
//	    fetchData(),
//	    readerio.ChainFirst(readerio.PrintGo[Config, Data]("{{if .R.Verbose}}[VERBOSE] {{end}}Data: {{.A.ID}} - {{.A.Value}}")),
//	    processData,
//	)(Config{Verbose: true})()
func PrintGo[R, A any](prefix string) Kleisli[R, A, A] {
	return handleLogging[R, A](func(value string) {
		fmt.Println(value)
	}, func(err error) {
		fmt.Fprintln(os.Stderr, err)
	}, prefix)
}
