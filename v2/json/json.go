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
