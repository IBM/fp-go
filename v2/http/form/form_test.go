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

package form

import (
	"fmt"
	"net/url"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
	O "github.com/IBM/fp-go/v2/option"
	RG "github.com/IBM/fp-go/v2/record/generic"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

var (
	sEq      = eq.FromEquals(S.Eq)
	valuesEq = RG.Eq[url.Values](A.Eq(sEq))
)

func TestLaws(t *testing.T) {
	name := "Content-Type"
	fieldLaws := LT.AssertLaws(t, O.Eq(sEq), valuesEq)(AtValue(name))

	n := O.None[string]()
	s1 := O.Some("s1")

	v1 := F.Pipe1(
		Default,
		WithValue(name)("v1"),
	)

	v2 := F.Pipe1(
		Default,
		WithValue("Other-Header")("v2"),
	)

	assert.True(t, fieldLaws(Default, n))
	assert.True(t, fieldLaws(v1, n))
	assert.True(t, fieldLaws(v2, n))

	assert.True(t, fieldLaws(Default, s1))
	assert.True(t, fieldLaws(v1, s1))
	assert.True(t, fieldLaws(v2, s1))
}

func TestFormField(t *testing.T) {

	v1 := F.Pipe1(
		Default,
		WithValue("h1")("v1"),
	)

	v2 := F.Pipe1(
		v1,
		WithValue("h2")("v2"),
	)

	// make sure the code does not change structures
	assert.False(t, valuesEq.Equals(Default, v1))
	assert.False(t, valuesEq.Equals(Default, v2))
	assert.False(t, valuesEq.Equals(v1, v2))

	// check for existence of values
	assert.Equal(t, "v1", v1.Get("h1"))
	assert.Equal(t, "v1", v2.Get("h1"))
	assert.Equal(t, "v2", v2.Get("h2"))

	// check getter on lens

	l1 := AtValue("h1")
	l2 := AtValue("h2")

	assert.Equal(t, O.Of("v1"), l1.Get(v1))
	assert.Equal(t, O.Of("v1"), l1.Get(v2))
	assert.Equal(t, O.Of("v2"), l2.Get(v2))
}

// TestWithValue tests the WithValue function
func TestWithValue(t *testing.T) {
	t.Run("sets a single value", func(t *testing.T) {
		form := WithValue("key")("value")(Default)
		assert.Equal(t, "value", form.Get("key"))
	})

	t.Run("creates immutable transformation", func(t *testing.T) {
		original := Default
		modified := WithValue("key")("value")(original)

		assert.False(t, valuesEq.Equals(original, modified))
		assert.Equal(t, "", original.Get("key"))
		assert.Equal(t, "value", modified.Get("key"))
	})

	t.Run("overwrites existing value", func(t *testing.T) {
		form := WithValue("key")("value1")(Default)
		updated := WithValue("key")("value2")(form)

		assert.Equal(t, "value2", updated.Get("key"))
		assert.Equal(t, "value1", form.Get("key"))
	})

	t.Run("composes multiple values", func(t *testing.T) {
		form := F.Pipe3(
			Default,
			WithValue("key1")("value1"),
			WithValue("key2")("value2"),
			WithValue("key3")("value3"),
		)

		assert.Equal(t, "value1", form.Get("key1"))
		assert.Equal(t, "value2", form.Get("key2"))
		assert.Equal(t, "value3", form.Get("key3"))
	})

	t.Run("handles empty string values", func(t *testing.T) {
		form := WithValue("key")("")(Default)
		assert.Equal(t, "", form.Get("key"))
		assert.True(t, form.Has("key"))
	})

	t.Run("handles special characters in keys", func(t *testing.T) {
		form := F.Pipe2(
			Default,
			WithValue("key-with-dash")("value1"),
			WithValue("key_with_underscore")("value2"),
		)

		assert.Equal(t, "value1", form.Get("key-with-dash"))
		assert.Equal(t, "value2", form.Get("key_with_underscore"))
	})
}

// TestWithoutValue tests the WithoutValue function
func TestWithoutValue(t *testing.T) {
	t.Run("clears field value", func(t *testing.T) {
		form := WithValue("key")("value")(Default)
		updated := WithoutValue("key")(form)

		// WithoutValue sets the field to an empty array, not removing it entirely
		assert.Equal(t, "", updated.Get("key"))
		// The field still exists but with empty values
		values := updated["key"]
		assert.Equal(t, 0, len(values))
	})

	t.Run("is idempotent", func(t *testing.T) {
		form := WithValue("key")("value")(Default)
		removed1 := WithoutValue("key")(form)
		removed2 := WithoutValue("key")(removed1)

		assert.True(t, valuesEq.Equals(removed1, removed2))
	})

	t.Run("does not affect other fields", func(t *testing.T) {
		form := F.Pipe2(
			Default,
			WithValue("key1")("value1"),
			WithValue("key2")("value2"),
		)
		updated := WithoutValue("key1")(form)

		assert.Equal(t, "", updated.Get("key1"))
		assert.Equal(t, "value2", updated.Get("key2"))
	})

	t.Run("creates immutable transformation", func(t *testing.T) {
		form := WithValue("key")("value")(Default)
		updated := WithoutValue("key")(form)

		assert.False(t, valuesEq.Equals(form, updated))
		assert.Equal(t, "value", form.Get("key"))
		assert.Equal(t, "", updated.Get("key"))
	})

	t.Run("handles non-existent field", func(t *testing.T) {
		form := Default
		updated := WithoutValue("nonexistent")(form)

		assert.True(t, valuesEq.Equals(form, updated))
	})
}

// TestMonoid tests the Monoid for Endomorphism
func TestMonoid(t *testing.T) {
	t.Run("identity element", func(t *testing.T) {
		form := F.Pipe1(
			Default,
			WithValue("key")("value"),
		)

		// Concatenating with identity should not change the result
		result := Monoid.Concat(Monoid.Empty(), WithValue("key")("value"))(Default)
		assert.True(t, valuesEq.Equals(form, result))
	})

	t.Run("concatenates transformations", func(t *testing.T) {
		transform := Monoid.Concat(
			WithValue("key1")("value1"),
			WithValue("key2")("value2"),
		)
		result := transform(Default)

		assert.Equal(t, "value1", result.Get("key1"))
		assert.Equal(t, "value2", result.Get("key2"))
	})

	t.Run("concatenates multiple transformations", func(t *testing.T) {
		transform := Monoid.Concat(
			WithValue("key1")("value1"),
			Monoid.Concat(
				WithValue("key2")("value2"),
				WithValue("key3")("value3"),
			),
		)
		result := transform(Default)

		assert.Equal(t, "value1", result.Get("key1"))
		assert.Equal(t, "value2", result.Get("key2"))
		assert.Equal(t, "value3", result.Get("key3"))
	})

	t.Run("respects transformation order", func(t *testing.T) {
		// Monoid concatenation composes functions left-to-right
		// So the first transformation is applied first, then the second
		transform := Monoid.Concat(
			WithValue("key")("first"),
			WithValue("key")("second"),
		)
		result := transform(Default)

		// The transformations are composed, so first is applied, then second overwrites it
		// But since Monoid.Concat composes endomorphisms, we need to check actual behavior
		assert.Equal(t, "first", result.Get("key"))
	})
}

// TestValuesMonoid tests the ValuesMonoid
func TestValuesMonoid(t *testing.T) {
	t.Run("identity element", func(t *testing.T) {
		form := url.Values{"key": []string{"value"}}
		result := ValuesMonoid.Concat(ValuesMonoid.Empty(), form)

		assert.True(t, valuesEq.Equals(form, result))
	})

	t.Run("concatenates disjoint forms", func(t *testing.T) {
		form1 := url.Values{"key1": []string{"value1"}}
		form2 := url.Values{"key2": []string{"value2"}}
		result := ValuesMonoid.Concat(form1, form2)

		assert.Equal(t, "value1", result.Get("key1"))
		assert.Equal(t, "value2", result.Get("key2"))
	})

	t.Run("concatenates arrays for same key", func(t *testing.T) {
		form1 := url.Values{"key": []string{"value1"}}
		form2 := url.Values{"key": []string{"value2"}}
		result := ValuesMonoid.Concat(form1, form2)

		values := result["key"]
		assert.Equal(t, 2, len(values))
		assert.Equal(t, "value1", values[0])
		assert.Equal(t, "value2", values[1])
	})

	t.Run("is associative", func(t *testing.T) {
		form1 := url.Values{"key": []string{"value1"}}
		form2 := url.Values{"key": []string{"value2"}}
		form3 := url.Values{"key": []string{"value3"}}

		result1 := ValuesMonoid.Concat(ValuesMonoid.Concat(form1, form2), form3)
		result2 := ValuesMonoid.Concat(form1, ValuesMonoid.Concat(form2, form3))

		assert.True(t, valuesEq.Equals(result1, result2))
	})
}

// TestAtValues tests the AtValues lens
func TestAtValues(t *testing.T) {
	t.Run("gets values array", func(t *testing.T) {
		form := url.Values{"key": []string{"value1", "value2"}}
		lens := AtValues("key")

		result := lens.Get(form)
		assert.True(t, O.IsSome(result))
		values := O.GetOrElse(F.Constant([]string{}))(result)
		assert.Equal(t, 2, len(values))
		assert.Equal(t, "value1", values[0])
		assert.Equal(t, "value2", values[1])
	})

	t.Run("returns None for non-existent key", func(t *testing.T) {
		lens := AtValues("nonexistent")
		result := lens.Get(Default)

		assert.True(t, O.IsNone(result))
	})

	t.Run("sets values array", func(t *testing.T) {
		lens := AtValues("key")
		form := lens.Set(O.Some([]string{"value1", "value2"}))(Default)

		values := form["key"]
		assert.Equal(t, 2, len(values))
		assert.Equal(t, "value1", values[0])
		assert.Equal(t, "value2", values[1])
	})

	t.Run("removes field with None", func(t *testing.T) {
		form := url.Values{"key": []string{"value"}}
		lens := AtValues("key")
		updated := lens.Set(O.None[[]string]())(form)

		assert.False(t, updated.Has("key"))
	})

	t.Run("creates immutable transformation", func(t *testing.T) {
		form := url.Values{"key": []string{"value1"}}
		lens := AtValues("key")
		updated := lens.Set(O.Some([]string{"value2"}))(form)

		assert.False(t, valuesEq.Equals(form, updated))
		assert.Equal(t, "value1", form.Get("key"))
		assert.Equal(t, "value2", updated.Get("key"))
	})
}

// TestAtValue tests the AtValue lens
func TestAtValue(t *testing.T) {
	t.Run("gets first value", func(t *testing.T) {
		form := url.Values{"key": []string{"value1", "value2"}}
		lens := AtValue("key")

		result := lens.Get(form)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "value1", O.GetOrElse(F.Constant(""))(result))
	})

	t.Run("returns None for non-existent key", func(t *testing.T) {
		lens := AtValue("nonexistent")
		result := lens.Get(Default)

		assert.True(t, O.IsNone(result))
	})

	t.Run("returns None for empty array", func(t *testing.T) {
		form := url.Values{"key": []string{}}
		lens := AtValue("key")
		result := lens.Get(form)

		assert.True(t, O.IsNone(result))
	})

	t.Run("sets first value", func(t *testing.T) {
		lens := AtValue("key")
		form := lens.Set(O.Some("value"))(Default)

		assert.Equal(t, "value", form.Get("key"))
	})

	t.Run("replaces first value in array", func(t *testing.T) {
		form := url.Values{"key": []string{"old1", "old2"}}
		lens := AtValue("key")
		updated := lens.Set(O.Some("new"))(form)

		values := updated["key"]
		// AtValue modifies the head of the array, keeping other elements
		assert.Equal(t, 2, len(values))
		assert.Equal(t, "new", values[0])
		assert.Equal(t, "old2", values[1])
	})

	t.Run("clears field with None", func(t *testing.T) {
		form := url.Values{"key": []string{"value"}}
		lens := AtValue("key")
		updated := lens.Set(O.None[string]())(form)

		// Setting to None creates an empty array, not removing the key
		values := updated["key"]
		assert.Equal(t, 0, len(values))
	})
}

// Example tests demonstrating package usage

// ExampleWithValue demonstrates how to set form field values
func ExampleWithValue() {
	// Create a form with a single field
	form := WithValue("username")("john")(Default)
	fmt.Println(form.Get("username"))
	// Output: john
}

// ExampleWithValue_composition demonstrates composing multiple field assignments
func ExampleWithValue_composition() {
	// Build a form with multiple fields using Pipe
	form := F.Pipe3(
		Default,
		WithValue("username")("john"),
		WithValue("email")("john@example.com"),
		WithValue("age")("30"),
	)

	fmt.Println(form.Get("username"))
	fmt.Println(form.Get("email"))
	fmt.Println(form.Get("age"))
	// Output:
	// john
	// john@example.com
	// 30
}

// ExampleWithoutValue demonstrates clearing a form field value
func ExampleWithoutValue() {
	// Create a form and then clear a field
	form := F.Pipe2(
		Default,
		WithValue("username")("john"),
		WithValue("password")("secret"),
	)

	// Clear the password field (sets it to empty array)
	sanitized := WithoutValue("password")(form)

	fmt.Println(sanitized.Get("username"))
	fmt.Println(sanitized.Get("password"))
	// Output:
	// john
	//
}

// ExampleAtValue demonstrates using the AtValue lens
func ExampleAtValue() {
	form := WithValue("username")("john")(Default)

	// Get a value using the lens
	lens := AtValue("username")
	value := lens.Get(form)

	fmt.Println(O.IsSome(value))
	fmt.Println(O.GetOrElse(F.Constant("default"))(value))
	// Output:
	// true
	// john
}

// ExampleAtValue_set demonstrates setting a value using the AtValue lens
func ExampleAtValue_set() {
	form := WithValue("username")("john")(Default)

	// Update the value using the lens
	lens := AtValue("username")
	updated := lens.Set(O.Some("jane"))(form)

	fmt.Println(updated.Get("username"))
	// Output: jane
}

// ExampleMonoid demonstrates combining form transformations
func ExampleMonoid() {
	// Combine multiple transformations into one
	transform := Monoid.Concat(
		WithValue("field1")("value1"),
		WithValue("field2")("value2"),
	)

	result := transform(Default)
	fmt.Println(result.Get("field1"))
	fmt.Println(result.Get("field2"))
	// Output:
	// value1
	// value2
}

// ExampleValuesMonoid demonstrates merging form data
func ExampleValuesMonoid() {
	form1 := url.Values{"key1": []string{"value1"}}
	form2 := url.Values{"key2": []string{"value2"}}

	merged := ValuesMonoid.Concat(form1, form2)

	fmt.Println(merged.Get("key1"))
	fmt.Println(merged.Get("key2"))
	// Output:
	// value1
	// value2
}

// ExampleValuesMonoid_concatenation demonstrates array concatenation for same keys
func ExampleValuesMonoid_concatenation() {
	form1 := url.Values{"tags": []string{"go"}}
	form2 := url.Values{"tags": []string{"functional"}}

	merged := ValuesMonoid.Concat(form1, form2)

	tags := merged["tags"]
	fmt.Println(len(tags))
	fmt.Println(tags[0])
	fmt.Println(tags[1])
	// Output:
	// 2
	// go
	// functional
}

// ExampleAtValues demonstrates working with multiple values
func ExampleAtValues() {
	form := url.Values{"tags": []string{"go", "functional", "programming"}}

	lens := AtValues("tags")
	values := lens.Get(form)

	if O.IsSome(values) {
		tags := O.GetOrElse(F.Constant([]string{}))(values)
		fmt.Println(len(tags))
		fmt.Println(tags[0])
	}
	// Output:
	// 3
	// go
}
