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

// Package prism provides utilities for composing prisms with isomorphisms.
//
// This package enables the composition of prisms (optics for sum types) with
// isomorphisms (bidirectional transformations), allowing you to transform the
// source type of a prism using an isomorphism. This is the inverse operation
// of optics/prism/iso, where we transform the focus type instead of the source type.
//
// # Key Concepts
//
// A Prism[S, A] is an optic that focuses on a specific variant within a sum type S,
// extracting values of type A. An Iso[S, A] represents a bidirectional transformation
// between types S and A without loss of information.
//
// When you compose a Prism[A, B] with an Iso[S, A], you get a Prism[S, B] that:
//   - Transforms S to A using the isomorphism's Get
//   - Extracts values of type B from A (using the prism)
//   - Can construct S from B by first using the prism's ReverseGet to get A, then the iso's ReverseGet
//
// # Example Usage
//
//	import (
//	    "github.com/IBM/fp-go/v2/optics/iso"
//	    "github.com/IBM/fp-go/v2/optics/prism"
//	    IP "github.com/IBM/fp-go/v2/optics/iso/prism"
//	    O "github.com/IBM/fp-go/v2/option"
//	)
//
//	// Create an isomorphism between []byte and string
//	bytesStringIso := iso.MakeIso(
//	    func(b []byte) string { return string(b) },
//	    func(s string) []byte { return []byte(s) },
//	)
//
//	// Create a prism that extracts Right values from Either[error, string]
//	rightPrism := prism.FromEither[error, string]()
//
//	// Compose them to get a prism that works with []byte as source
//	bytesPrism := IP.Compose(rightPrism)(bytesStringIso)
//
//	// Use the composed prism
//	bytes := []byte(`{"status":"ok"}`)
//	// First converts bytes to string via iso, then extracts Right value
//	result := bytesPrism.GetOption(either.Right[error](string(bytes)))
//
// # Comparison with optics/prism/iso
//
// This package (optics/iso/prism) is the dual of optics/prism/iso:
//   - optics/prism/iso: Composes Iso[A, B] with Prism[S, A] → Prism[S, B] (transforms focus type)
//   - optics/iso/prism: Composes Prism[A, B] with Iso[S, A] → Prism[S, B] (transforms source type)
//
// # Type Aliases
//
// This package re-exports key types from the iso and prism packages for convenience:
//   - Iso[S, A]: An isomorphism between types S and A
//   - Prism[S, A]: A prism focusing on type A within sum type S
package prism

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
)
