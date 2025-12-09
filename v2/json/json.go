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

package json

import (
	"encoding/json"

	E "github.com/IBM/fp-go/v2/either"
)

// Unmarshal parses JSON-encoded data and returns an Either containing the decoded value or an error.
//
// This function wraps the standard json.Unmarshal in an Either monad, converting the traditional
// (value, error) tuple into a functional Either type. If unmarshaling succeeds, it returns Right[A].
// If it fails, it returns Left[error].
//
// Type parameter A specifies the target type for unmarshaling. The type must be compatible with
// the JSON structure in the input data.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	data := []byte(`{"name":"Alice","age":30}`)
//	result := json.Unmarshal[Person](data)
//	// result is Either[error, Person]
//
//	either.Fold(
//	    func(err error) { fmt.Println("Error:", err) },
//	    func(p Person) { fmt.Printf("Success: %s, %d\n", p.Name, p.Age) },
//	)(result)
func Unmarshal[A any](data []byte) Either[A] {
	var result A
	err := json.Unmarshal(data, &result)
	return E.TryCatchError(result, err)
}

// Marshal converts a Go value to JSON-encoded bytes and returns an Either containing the result or an error.
//
// This function wraps the standard json.Marshal in an Either monad, converting the traditional
// (value, error) tuple into a functional Either type. If marshaling succeeds, it returns Right[[]byte].
// If it fails, it returns Left[error].
//
// The function uses the same encoding rules as the standard library's json.Marshal, including
// support for struct tags, custom MarshalJSON methods, and standard type conversions.
//
// Example:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	person := Person{Name: "Bob", Age: 25}
//	result := json.Marshal(person)
//	// result is Either[error, []byte]
//
//	either.Map(func(data []byte) string {
//	    return string(data)
//	})(result)
//	// Returns Either[error, string] with JSON string
func Marshal[A any](a A) Either[[]byte] {
	return E.TryCatchError(json.Marshal(a))
}

// MarshalIndent converts a Go value to pretty-printed JSON-encoded bytes with indentation.
//
// This function wraps the standard json.MarshalIndent in an Either monad, converting the traditional
// (value, error) tuple into a functional Either type. If marshaling succeeds, it returns Right[[]byte]
// containing the formatted JSON. If it fails, it returns Left[error].
//
// The function uses a default indentation of two spaces ("  ") with no prefix, making the output
// human-readable and suitable for display, logging, or configuration files. Each JSON element begins
// on a new line, and nested structures are indented to show their hierarchy.
//
// Type parameter A specifies the type of value to marshal. The type must be compatible with
// JSON encoding rules (same as json.Marshal).
//
// The function uses the same encoding rules as the standard library's json.MarshalIndent, including:
//   - Support for struct tags to control field names and omitempty behavior
//   - Custom MarshalJSON methods for types that implement json.Marshaler
//   - Standard type conversions (strings, numbers, booleans, arrays, slices, maps, structs)
//   - Proper escaping of special characters in strings
//
// Example with a simple struct:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	person := Person{Name: "Alice", Age: 30}
//	result := json.MarshalIndent(person)
//	// result is Either[error, []byte]
//
//	either.Map(func(data []byte) string {
//	    return string(data)
//	})(result)
//	// Returns Either[error, string] with formatted JSON:
//	// {
//	//   "name": "Alice",
//	//   "age": 30
//	// }
//
// Example with nested structures:
//
//	type Address struct {
//	    Street string `json:"street"`
//	    City   string `json:"city"`
//	}
//
//	type Employee struct {
//	    Name    string  `json:"name"`
//	    Address Address `json:"address"`
//	}
//
//	emp := Employee{
//	    Name: "Bob",
//	    Address: Address{Street: "123 Main St", City: "Boston"},
//	}
//	result := json.MarshalIndent(emp)
//	// Produces formatted JSON:
//	// {
//	//   "name": "Bob",
//	//   "address": {
//	//     "street": "123 Main St",
//	//     "city": "Boston"
//	//   }
//	// }
//
// Example with error handling:
//
//	type Config struct {
//	    Settings map[string]interface{} `json:"settings"`
//	}
//
//	config := Config{Settings: map[string]interface{}{"debug": true}}
//	result := json.MarshalIndent(config)
//
//	either.Fold(
//	    func(err error) string {
//	        return fmt.Sprintf("Failed to marshal: %v", err)
//	    },
//	    func(data []byte) string {
//	        return string(data)
//	    },
//	)(result)
//
// Example with functional composition:
//
//	// Chain operations using Either monad
//	result := F.Pipe2(
//	    person,
//	    json.MarshalIndent[Person],
//	    either.Map(func(data []byte) string {
//	        return string(data)
//	    }),
//	)
//	// result is Either[error, string] with formatted JSON
//
// Use MarshalIndent when you need human-readable JSON output for:
//   - Configuration files that humans will read or edit
//   - Debug output and logging
//   - API responses for development/testing
//   - Documentation examples
//
// Use Marshal (without indentation) when:
//   - Minimizing payload size is important (production APIs)
//   - The JSON will be consumed by machines only
//   - Performance is critical (indentation adds overhead)
func MarshalIndent[A any](a A) Either[[]byte] {
	return E.TryCatchError(json.MarshalIndent(a, "", "  "))
}
