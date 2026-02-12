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

// Package iso provides utilities for composing isomorphisms with prisms.
//
// This package enables the composition of isomorphisms (bidirectional transformations)
// with prisms (optics for sum types), allowing you to transform the focus type of a prism
// using an isomorphism. This is particularly useful when you need to work with prisms
// that focus on a type that can be bidirectionally converted to another type.
//
// # Key Concepts
//
// An Iso[S, A] represents a bidirectional transformation between types S and A without
// loss of information. A Prism[S, A] is an optic that focuses on a specific variant
// within a sum type S, extracting values of type A.
//
// When you compose an Iso[A, B] with a Prism[S, A], you get a Prism[S, B] that:
//   - Extracts values of type A from S (using the prism)
//   - Transforms them to type B (using the isomorphism's Get)
//   - Can construct S from B by reversing the transformation (using the isomorphism's ReverseGet)
//
// # Example Usage
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/iso"
//	    "github.com/IBM/fp-go/v2/optics/prism"
//	    PI "github.com/IBM/fp-go/v2/optics/prism/iso"
//	    O "github.com/IBM/fp-go/v2/option"
//	)
//
//	// Create an isomorphism between string and []byte
//	stringBytesIso := iso.MakeIso(
//	    func(s string) []byte { return []byte(s) },
//	    func(b []byte) string { return string(b) },
//	)
//
//	// Create a prism that extracts Right values from Either[error, string]
//	rightPrism := prism.FromEither[error, string]()
//
//	// Compose them to get a prism that extracts Right values as []byte
//	bytesPrism := PI.Compose(stringBytesIso)(rightPrism)
//
//	// Use the composed prism
//	either := either.Right[error]("hello")
//	result := bytesPrism.GetOption(either)  // Some([]byte("hello"))
//
// # Type Aliases
//
// This package re-exports key types from the iso and prism packages for convenience:
//   - Iso[S, A]: An isomorphism between types S and A
//   - Prism[S, A]: A prism focusing on type A within sum type S
//   - Operator[S, A, B]: A function that transforms Prism[S, A] to Prism[S, B]
package iso

import (
	I "github.com/IBM/fp-go/v2/optics/iso"
	P "github.com/IBM/fp-go/v2/optics/prism"
)

type (
	// Iso represents an isomorphism between types S and A.
	// It is a bidirectional transformation that converts between two types
	// without any loss of information.
	//
	// Type Parameters:
	//   - S: The source type
	//   - A: The target type
	//
	// See github.com/IBM/fp-go/v2/optics/iso for the full Iso API.
	Iso[S, A any] = I.Iso[S, A]

	// Prism is an optic used to select part of a sum type (tagged union).
	// It provides operations to extract and construct values within sum types.
	//
	// Type Parameters:
	//   - S: The source type (sum type)
	//   - A: The focus type (variant within the sum type)
	//
	// See github.com/IBM/fp-go/v2/optics/prism for the full Prism API.
	Prism[S, A any] = P.Prism[S, A]

	// Operator represents a function that transforms one prism into another.
	// It takes a Prism[S, A] and returns a Prism[S, B], allowing for prism transformations.
	//
	// This is commonly used with the Compose function to create operators that
	// transform the focus type of a prism using an isomorphism.
	//
	// Type Parameters:
	//   - S: The source type (remains constant)
	//   - A: The original focus type
	//   - B: The new focus type
	//
	// Example:
	//
	//	// Create an operator that transforms string prisms to []byte prisms
	//	stringToBytesOp := Compose(stringBytesIso)
	//	// Apply it to a prism
	//	bytesPrism := stringToBytesOp(stringPrism)
	Operator[S, A, B any] = P.Operator[S, A, B]
)
