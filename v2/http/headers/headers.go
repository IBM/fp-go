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

// Package headers provides constants and utilities for working with HTTP headers
// in a functional programming style. It offers type-safe header name constants,
// monoid operations for combining headers, and lens-based access to header values.
//
// The package follows functional programming principles by providing:
//   - Immutable operations through lenses
//   - Monoid for combining header maps
//   - Type-safe header name constants
//   - Functional composition of header operations
//
// Constants:
//
// The package defines commonly used HTTP header names as constants:
//   - Accept: The Accept request header
//   - Authorization: The Authorization request header
//   - ContentType: The Content-Type header
//   - ContentLength: The Content-Length header
//
// Monoid:
//
// The Monoid provides a way to combine multiple http.Header maps:
//
//	headers1 := make(http.Header)
//	headers1.Set("X-Custom", "value1")
//
//	headers2 := make(http.Header)
//	headers2.Set("Authorization", "Bearer token")
//
//	combined := Monoid.Concat(headers1, headers2)
//	// combined now contains both headers
//
// Lenses:
//
// AtValues and AtValue provide lens-based access to header values:
//
//	// AtValues focuses on all values of a header ([]string)
//	contentTypeLens := AtValues("Content-Type")
//	values := contentTypeLens.Get(headers)
//
//	// AtValue focuses on the first value of a header (Option[string])
//	authLens := AtValue("Authorization")
//	token := authLens.Get(headers) // Returns Option[string]
//
// The lenses support functional updates:
//
//	// Set a header value
//	newHeaders := AtValue("Content-Type").Set(O.Some("application/json"))(headers)
//
//	// Remove a header
//	newHeaders := AtValue("X-Custom").Set(O.None[string]())(headers)
package headers

import (
	"net/http"
	"net/textproto"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	LA "github.com/IBM/fp-go/v2/optics/lens/array"
	LO "github.com/IBM/fp-go/v2/optics/lens/option"
	LRG "github.com/IBM/fp-go/v2/optics/lens/record/generic"
	RG "github.com/IBM/fp-go/v2/record/generic"
)

// Common HTTP header name constants.
// These constants provide type-safe access to standard HTTP header names.
const (
	// Accept specifies the media types that are acceptable for the response.
	// Example: "Accept: application/json"
	Accept = "Accept"

	// Authorization contains credentials for authenticating the client with the server.
	// Example: "Authorization: Bearer token123"
	Authorization = "Authorization"

	// ContentType indicates the media type of the resource or data.
	// Example: "Content-Type: application/json"
	ContentType = "Content-Type"

	// ContentLength indicates the size of the entity-body in bytes.
	// Example: "Content-Length: 348"
	ContentLength = "Content-Length"
)

var (
	// Monoid is a Monoid for combining http.Header maps.
	// It uses a union operation where values from both headers are preserved.
	// When the same header exists in both maps, the values are concatenated.
	//
	// Example:
	//   h1 := make(http.Header)
	//   h1.Set("X-Custom", "value1")
	//
	//   h2 := make(http.Header)
	//   h2.Set("Authorization", "Bearer token")
	//
	//   combined := Monoid.Concat(h1, h2)
	//   // combined contains both X-Custom and Authorization headers
	Monoid = RG.UnionMonoid[http.Header](A.Semigroup[string]())

	// AtValues is a Lens that focuses on all values of a specific header.
	// It returns a lens that accesses the []string slice of header values.
	// The header name is automatically canonicalized using MIME header key rules.
	//
	// Parameters:
	//   - name: The header name (will be canonicalized)
	//
	// Returns:
	//   - A Lens[http.Header, []string] focusing on the header's values
	//
	// Example:
	//   lens := AtValues("Content-Type")
	//   values := lens.Get(headers) // Returns []string
	//   newHeaders := lens.Set([]string{"application/json"})(headers)
	AtValues = F.Flow2(
		textproto.CanonicalMIMEHeaderKey,
		LRG.AtRecord[http.Header, []string],
	)

	// composeHead is an internal helper that composes a lens to focus on the first
	// element of a string array, returning an Option[string].
	composeHead = F.Pipe1(
		LA.AtHead[string](),
		LO.Compose[http.Header, string](A.Empty[string]()),
	)

	// AtValue is a Lens that focuses on the first value of a specific header.
	// It returns a lens that accesses an Option[string] representing the first
	// header value, or None if the header doesn't exist.
	// The header name is automatically canonicalized using MIME header key rules.
	//
	// Parameters:
	//   - name: The header name (will be canonicalized)
	//
	// Returns:
	//   - A Lens[http.Header, Option[string]] focusing on the first header value
	//
	// Example:
	//   lens := AtValue("Authorization")
	//   token := lens.Get(headers) // Returns Option[string]
	//
	//   // Set a header value
	//   newHeaders := lens.Set(O.Some("Bearer token"))(headers)
	//
	//   // Remove a header
	//   newHeaders := lens.Set(O.None[string]())(headers)
	AtValue = F.Flow2(
		AtValues,
		composeHead,
	)
)
