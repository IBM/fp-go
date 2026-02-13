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

package assert

import (
	"testing"

	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/readerio"
)

// Logf creates a logging function that outputs formatted test messages using Go's testing.T.Logf.
//
// This function provides a functional programming approach to test logging, returning a
// [ReaderIO] that can be composed with other test operations. It's particularly useful
// for debugging tests, tracing execution flow, or documenting test behavior without
// affecting test outcomes.
//
// The function uses a curried design pattern:
//  1. First, you provide a format string (prefix) with format verbs (like %v, %d, %s)
//  2. This returns a function that takes a value of type T
//  3. That function returns a ReaderIO that performs the logging when executed
//
// # Parameters
//
//   - prefix: A format string compatible with fmt.Printf (e.g., "Value: %v", "Count: %d")
//     The format string should contain exactly one format verb that matches type T
//
// # Returns
//
//   - A function that takes a value of type T and returns a [ReaderIO][*testing.T, Void]
//     When executed, this ReaderIO logs the formatted message to the test output
//
// # Type Parameters
//
//   - T: The type of value to be logged. Can be any type that can be formatted by fmt
//
// # Use Cases
//
//   - Debugging test execution by logging intermediate values
//   - Tracing the flow of complex test scenarios
//   - Documenting test behavior in the test output
//   - Logging values in functional pipelines without breaking the chain
//   - Creating reusable logging operations for specific types
//
// # Example - Basic Logging
//
//	func TestBasicLogging(t *testing.T) {
//	    // Create a logger for integers
//	    logInt := assert.Logf[int]("Processing value: %d")
//
//	    // Use it to log a value
//	    value := 42
//	    logInt(value)(t)()  // Outputs: "Processing value: 42"
//	}
//
// # Example - Logging in Test Pipeline
//
//	func TestPipelineWithLogging(t *testing.T) {
//	    type User struct {
//	        Name string
//	        Age  int
//	    }
//
//	    user := User{Name: "Alice", Age: 30}
//
//	    // Create a logger for User
//	    logUser := assert.Logf[User]("Testing user: %+v")
//
//	    // Log the user being tested
//	    logUser(user)(t)()
//
//	    // Continue with assertions
//	    assert.StringNotEmpty(user.Name)(t)
//	    assert.That(func(age int) bool { return age > 0 })(user.Age)(t)
//	}
//
// # Example - Multiple Loggers for Different Types
//
//	func TestMultipleLoggers(t *testing.T) {
//	    // Create type-specific loggers
//	    logString := assert.Logf[string]("String value: %s")
//	    logInt := assert.Logf[int]("Integer value: %d")
//	    logFloat := assert.Logf[float64]("Float value: %.2f")
//
//	    // Use them throughout the test
//	    logString("hello")(t)()      // Outputs: "String value: hello"
//	    logInt(42)(t)()               // Outputs: "Integer value: 42"
//	    logFloat(3.14159)(t)()        // Outputs: "Float value: 3.14"
//	}
//
// # Example - Logging Complex Structures
//
//	func TestComplexStructureLogging(t *testing.T) {
//	    type Config struct {
//	        Host    string
//	        Port    int
//	        Timeout int
//	    }
//
//	    config := Config{Host: "localhost", Port: 8080, Timeout: 30}
//
//	    // Use %+v to include field names
//	    logConfig := assert.Logf[Config]("Configuration: %+v")
//	    logConfig(config)(t)()
//	    // Outputs: "Configuration: {Host:localhost Port:8080 Timeout:30}"
//
//	    // Or use %#v for Go-syntax representation
//	    logConfigGo := assert.Logf[Config]("Config (Go syntax): %#v")
//	    logConfigGo(config)(t)()
//	    // Outputs: "Config (Go syntax): assert.Config{Host:"localhost", Port:8080, Timeout:30}"
//	}
//
// # Example - Debugging Test Failures
//
//	func TestWithDebugLogging(t *testing.T) {
//	    numbers := []int{1, 2, 3, 4, 5}
//	    logSlice := assert.Logf[[]int]("Testing slice: %v")
//
//	    // Log the input data
//	    logSlice(numbers)(t)()
//
//	    // Perform assertions
//	    assert.ArrayNotEmpty(numbers)(t)
//	    assert.ArrayLength[int](5)(numbers)(t)
//
//	    // Log intermediate results
//	    sum := 0
//	    for _, n := range numbers {
//	        sum += n
//	    }
//	    logInt := assert.Logf[int]("Sum: %d")
//	    logInt(sum)(t)()
//
//	    assert.Equal(15)(sum)(t)
//	}
//
// # Example - Conditional Logging
//
//	func TestConditionalLogging(t *testing.T) {
//	    logDebug := assert.Logf[string]("DEBUG: %s")
//
//	    values := []int{1, 2, 3, 4, 5}
//	    for _, v := range values {
//	        if v%2 == 0 {
//	            logDebug(fmt.Sprintf("Found even number: %d", v))(t)()
//	        }
//	    }
//	    // Outputs:
//	    // DEBUG: Found even number: 2
//	    // DEBUG: Found even number: 4
//	}
//
// # Format Verbs
//
// Common format verbs you can use in the prefix string:
//   - %v: Default format
//   - %+v: Default format with field names for structs
//   - %#v: Go-syntax representation
//   - %T: Type of the value
//   - %d: Integer in base 10
//   - %s: String
//   - %f: Floating point number
//   - %t: Boolean (true/false)
//   - %p: Pointer address
//
// See the fmt package documentation for a complete list of format verbs.
//
// # Notes
//
//   - Logging does not affect test pass/fail status
//   - Log output appears in test results when running with -v flag or when tests fail
//   - The function returns Void, indicating it's used for side effects only
//   - The ReaderIO pattern allows logging to be composed with other operations
//
// # Related Functions
//
//   - [FromReaderIO]: Converts ReaderIO operations into test assertions
//   - testing.T.Logf: The underlying Go testing log function
//
// # References
//
//   - Go testing package: https://pkg.go.dev/testing
//   - fmt package format verbs: https://pkg.go.dev/fmt
//   - ReaderIO pattern: Combines Reader (context dependency) with IO (side effects)
func Logf[T any](prefix string) func(T) readerio.ReaderIO[*testing.T, Void] {
	return func(a T) readerio.ReaderIO[*testing.T, Void] {
		return func(t *testing.T) IO[Void] {
			return io.FromImpure(func() {
				t.Logf(prefix, a)
			})
		}
	}
}
