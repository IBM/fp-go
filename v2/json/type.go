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
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Either is a type alias for either.Either[error, A], representing a computation that may fail.
	// By convention, Left contains an error and Right contains the successful result of type A.
	Either[A any] = E.Either[error, A]

	// Option is a type alias for option.Option[A], representing an optional value.
	// It can be either Some(value) or None.
	Option[A any] = option.Option[A]
)

// ToTypeE converts a value from one type to another using JSON as an intermediate format,
// returning an Either that contains the converted value or an error.
//
// This function performs a round-trip conversion: src → JSON → target type A.
// It's useful for converting between compatible types (e.g., map[string]any to a struct)
// or for deep copying values.
//
// The conversion will fail if:
//   - The source value cannot be marshaled to JSON
//   - The JSON cannot be unmarshaled into the target type A
//   - The JSON structure doesn't match the target type's structure
//
// Example:
//
//	type Source struct {
//	    Name  string
//	    Value int
//	}
//
//	type Target struct {
//	    Name  string `json:"name"`
//	    Value int    `json:"value"`
//	}
//
//	src := Source{Name: "test", Value: 42}
//	result := json.ToTypeE[Target](src)
//	// result is Either[error, Target]
//
//	// Converting from map to struct
//	data := map[string]any{"name": "Alice", "value": 100}
//	person := json.ToTypeE[Target](data)
func ToTypeE[A any](src any) Either[A] {
	return function.Pipe2(
		src,
		Marshal[any],
		E.Chain(Unmarshal[A]),
	)
}

// ToTypeO converts a value from one type to another using JSON as an intermediate format,
// returning an Option that contains the converted value or None if conversion fails.
//
// This is a convenience wrapper around ToTypeE that discards error details and returns
// an Option instead. Use this when you only care about success/failure and don't need
// the specific error message.
//
// The conversion follows the same rules as ToTypeE, performing a round-trip through JSON.
//
// Example:
//
//	type Config struct {
//	    Host string `json:"host"`
//	    Port int    `json:"port"`
//	}
//
//	data := map[string]any{"host": "localhost", "port": 8080}
//	maybeConfig := json.ToTypeO[Config](data)
//	// maybeConfig is Option[Config]
//
//	option.Fold(
//	    func() { fmt.Println("Conversion failed") },
//	    func(cfg Config) { fmt.Printf("Config: %s:%d\n", cfg.Host, cfg.Port) },
//	)(maybeConfig)
func ToTypeO[A any](src any) Option[A] {
	return function.Pipe1(
		ToTypeE[A](src),
		E.ToOption[error, A],
	)
}
