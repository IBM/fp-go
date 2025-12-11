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

package consumer

// Local transforms a Consumer by preprocessing its input through a function.
// This is the contravariant map operation for Consumers, analogous to reader.Local
// but operating on the input side rather than the output side.
//
// Given a Consumer[R1] that consumes values of type R1, and a function f that
// converts R2 to R1, Local creates a new Consumer[R2] that:
//   1. Takes a value of type R2
//   2. Applies f to convert it to R1
//   3. Passes the result to the original Consumer[R1]
//
// This is particularly useful for adapting consumers to work with different input types,
// similar to how reader.Local adapts readers to work with different environment types.
//
// Comparison with reader.Local:
//   - reader.Local: Transforms the environment BEFORE passing it to a Reader (preprocessing input)
//   - consumer.Local: Transforms the value BEFORE passing it to a Consumer (preprocessing input)
//   - Both are contravariant operations on the input type
//   - Reader produces output, Consumer performs side effects
//
// Type Parameters:
//   - R2: The input type of the new Consumer (what you have)
//   - R1: The input type of the original Consumer (what it expects)
//
// Parameters:
//   - f: A function that converts R2 to R1 (preprocessing function)
//
// Returns:
//   - An Operator that transforms Consumer[R1] into Consumer[R2]
//
// Example - Basic type adaptation:
//
//	// Consumer that logs integers
//	logInt := func(x int) {
//	    fmt.Printf("Value: %d\n", x)
//	}
//
//	// Adapt it to consume strings by parsing them first
//	parseToInt := func(s string) int {
//	    n, _ := strconv.Atoi(s)
//	    return n
//	}
//
//	logString := consumer.Local(parseToInt)(logInt)
//	logString("42") // Logs: "Value: 42"
//
// Example - Extracting fields from structs:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	// Consumer that logs names
//	logName := func(name string) {
//	    fmt.Printf("Name: %s\n", name)
//	}
//
//	// Adapt it to consume User structs
//	extractName := func(u User) string {
//	    return u.Name
//	}
//
//	logUser := consumer.Local(extractName)(logName)
//	logUser(User{Name: "Alice", Age: 30}) // Logs: "Name: Alice"
//
// Example - Simplifying complex types:
//
//	type DetailedConfig struct {
//	    Host     string
//	    Port     int
//	    Timeout  time.Duration
//	    MaxRetry int
//	}
//
//	type SimpleConfig struct {
//	    Host string
//	    Port int
//	}
//
//	// Consumer that logs simple configs
//	logSimple := func(c SimpleConfig) {
//	    fmt.Printf("Server: %s:%d\n", c.Host, c.Port)
//	}
//
//	// Adapt it to consume detailed configs
//	simplify := func(d DetailedConfig) SimpleConfig {
//	    return SimpleConfig{Host: d.Host, Port: d.Port}
//	}
//
//	logDetailed := consumer.Local(simplify)(logSimple)
//	logDetailed(DetailedConfig{
//	    Host:     "localhost",
//	    Port:     8080,
//	    Timeout:  time.Second,
//	    MaxRetry: 3,
//	}) // Logs: "Server: localhost:8080"
//
// Example - Composing multiple transformations:
//
//	type Response struct {
//	    StatusCode int
//	    Body       string
//	}
//
//	// Consumer that logs status codes
//	logStatus := func(code int) {
//	    fmt.Printf("Status: %d\n", code)
//	}
//
//	// Extract status code from response
//	getStatus := func(r Response) int {
//	    return r.StatusCode
//	}
//
//	// Adapt to consume responses
//	logResponse := consumer.Local(getStatus)(logStatus)
//	logResponse(Response{StatusCode: 200, Body: "OK"}) // Logs: "Status: 200"
//
// Example - Using with multiple consumers:
//
//	type Event struct {
//	    Type      string
//	    Timestamp time.Time
//	    Data      map[string]any
//	}
//
//	// Consumers for different aspects
//	logType := func(t string) { fmt.Printf("Type: %s\n", t) }
//	logTime := func(t time.Time) { fmt.Printf("Time: %v\n", t) }
//
//	// Adapt them to consume events
//	logEventType := consumer.Local(func(e Event) string { return e.Type })(logType)
//	logEventTime := consumer.Local(func(e Event) time.Time { return e.Timestamp })(logTime)
//
//	event := Event{Type: "UserLogin", Timestamp: time.Now(), Data: nil}
//	logEventType(event) // Logs: "Type: UserLogin"
//	logEventTime(event) // Logs: "Time: ..."
//
// Use Cases:
//   - Type adaptation: Convert between different input types
//   - Field extraction: Extract specific fields from complex structures
//   - Data transformation: Preprocess data before consumption
//   - Interface adaptation: Adapt consumers to work with different interfaces
//   - Logging pipelines: Transform data before logging
//   - Event handling: Extract relevant data from events before processing
//
// Relationship to Reader:
// Consumer is the dual of Reader in category theory:
//   - Reader[R, A] = R -> A (produces output from environment)
//   - Consumer[A] = A -> () (consumes input, produces side effects)
//   - reader.Local transforms the environment before reading
//   - consumer.Local transforms the input before consuming
//   - Both are contravariant functors on their input type
func Local[R2, R1 any](f func(R2) R1) Operator[R1, R2] {
	return func(c Consumer[R1]) Consumer[R2] {
		return func(r2 R2) {
			c(f(r2))
		}
	}
}
