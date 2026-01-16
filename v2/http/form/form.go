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

// Package form provides functional utilities for working with HTTP form data (url.Values).
//
// This package offers a functional approach to building and manipulating HTTP form data
// using lenses, endomorphisms, and monoids. It enables immutable transformations of
// url.Values through composable operations.
//
// # Core Concepts
//
// The package is built around several key abstractions:
//   - Endomorphism: A function that transforms url.Values immutably
//   - Lenses: Optics for focusing on specific form fields
//   - Monoids: For combining form transformations and values
//
// # Basic Usage
//
// Create form data by composing endomorphisms:
//
//	form := F.Pipe3(
//	    form.Default,
//	    form.WithValue("username")("john"),
//	    form.WithValue("email")("john@example.com"),
//	    form.WithValue("age")("30"),
//	)
//
// Remove fields from forms:
//
//	updated := F.Pipe1(
//	    form,
//	    form.WithoutValue("age"),
//	)
//
// # Lenses
//
// The package provides two main lenses:
//   - AtValues: Focuses on all values of a form field ([]string)
//   - AtValue: Focuses on the first value of a form field (Option[string])
//
// Use lenses to read and update form fields:
//
//	lens := form.AtValue("username")
//	value := lens.Get(form)  // Returns Option[string]
//	updated := lens.Set(O.Some("jane"))(form)
//
// # Monoids
//
// Combine multiple form transformations:
//
//	transform := form.Monoid.Concat(
//	    form.WithValue("field1")("value1"),
//	    form.WithValue("field2")("value2"),
//	)
//	result := transform(form.Default)
//
// Merge form values:
//
//	merged := form.ValuesMonoid.Concat(form1, form2)
package form

import (
	"net/url"

	A "github.com/IBM/fp-go/v2/array"
	ENDO "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	LA "github.com/IBM/fp-go/v2/optics/lens/array"
	LO "github.com/IBM/fp-go/v2/optics/lens/option"
	LRG "github.com/IBM/fp-go/v2/optics/lens/record/generic"
	O "github.com/IBM/fp-go/v2/option"
	RG "github.com/IBM/fp-go/v2/record/generic"
)

type (
	// Endomorphism is a function that transforms url.Values immutably.
	// It represents a transformation from url.Values to url.Values,
	// enabling functional composition of form modifications.
	//
	// Example:
	//	transform := form.WithValue("key")("value")
	//	result := transform(form.Default)
	Endomorphism = ENDO.Endomorphism[url.Values]
)

var (
	// Default is an empty url.Values that serves as the starting point
	// for building form data. Use this with Pipe operations to construct
	// forms functionally.
	//
	// Example:
	//	form := F.Pipe2(
	//	    form.Default,
	//	    form.WithValue("key1")("value1"),
	//	    form.WithValue("key2")("value2"),
	//	)
	Default = make(url.Values)

	noField = O.None[string]()

	// Monoid is a Monoid for Endomorphism that allows combining multiple
	// form transformations into a single transformation. The identity element
	// is the identity function, and concatenation composes transformations.
	//
	// Example:
	//	transform := form.Monoid.Concat(
	//	    form.WithValue("field1")("value1"),
	//	    form.WithValue("field2")("value2"),
	//	)
	//	result := transform(form.Default)
	Monoid = ENDO.Monoid[url.Values]()

	// ValuesMonoid is a Monoid for url.Values that concatenates form data.
	// When two forms are combined, arrays of values for the same key are
	// concatenated using the array Semigroup.
	//
	// Example:
	//	form1 := url.Values{"key": []string{"value1"}}
	//	form2 := url.Values{"key": []string{"value2"}}
	//	merged := form.ValuesMonoid.Concat(form1, form2)
	//	// Result: url.Values{"key": []string{"value1", "value2"}}
	ValuesMonoid = RG.UnionMonoid[url.Values](A.Semigroup[string]())

	// AtValues is a Lens that focuses on all values of a form field as a slice.
	// It provides access to the complete []string array for a given field name.
	//
	// Example:
	//	lens := form.AtValues("tags")
	//	values := lens.Get(form)  // Returns Option[[]string]
	//	updated := lens.Set(O.Some([]string{"tag1", "tag2"}))(form)
	AtValues = LRG.AtRecord[url.Values, []string]

	composeHead = F.Pipe1(
		LA.AtHead[string](),
		LO.Compose[url.Values, string](A.Empty[string]()),
	)

	// AtValue is a Lens that focuses on the first value of a form field.
	// It returns an Option[string] representing the first value if present,
	// or None if the field doesn't exist or has no values.
	//
	// Example:
	//	lens := form.AtValue("username")
	//	value := lens.Get(form)  // Returns Option[string]
	//	updated := lens.Set(O.Some("newuser"))(form)
	AtValue = F.Flow2(
		AtValues,
		composeHead,
	)
)

// WithValue creates an Endomorphism that sets a form field to a specific value.
// It returns a curried function that takes the field name first, then the value,
// and finally returns a transformation function.
//
// The transformation is immutable - it creates a new url.Values rather than
// modifying the input.
//
// Example:
//
//	// Set a single field
//	form := form.WithValue("username")("john")(form.Default)
//
//	// Compose multiple fields
//	form := F.Pipe3(
//	    form.Default,
//	    form.WithValue("username")("john"),
//	    form.WithValue("email")("john@example.com"),
//	    form.WithValue("age")("30"),
//	)
func WithValue(name string) func(value string) Endomorphism {
	return F.Flow2(
		O.Of[string],
		AtValue(name).Set,
	)
}

// WithoutValue creates an Endomorphism that removes a form field.
// The transformation is immutable - it creates a new url.Values rather than
// modifying the input.
//
// Example:
//
//	// Remove a field
//	updated := form.WithoutValue("age")(form)
//
//	// Compose with other operations
//	form := F.Pipe2(
//	    existingForm,
//	    form.WithValue("username")("john"),
//	    form.WithoutValue("password"),
//	)
func WithoutValue(name string) Endomorphism {
	return AtValue(name).Set(noField)
}
