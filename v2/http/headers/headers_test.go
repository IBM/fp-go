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

package headers

import (
	"net/http"
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
	valuesEq = RG.Eq[http.Header](A.Eq(sEq))
)

func TestLaws(t *testing.T) {
	name := "Content-Type"
	fieldLaws := LT.AssertLaws(t, O.Eq(sEq), valuesEq)(AtValue(name))

	n := O.None[string]()
	s1 := O.Some("s1")

	def := make(http.Header)

	v1 := make(http.Header)
	v1.Set(name, "v1")

	v2 := make(http.Header)
	v2.Set("Other-Header", "v2")

	assert.True(t, fieldLaws(def, n))
	assert.True(t, fieldLaws(v1, n))
	assert.True(t, fieldLaws(v2, n))

	assert.True(t, fieldLaws(def, s1))
	assert.True(t, fieldLaws(v1, s1))
	assert.True(t, fieldLaws(v2, s1))
}

// TestMonoidEmpty tests the Monoid empty (identity) element
func TestMonoidEmpty(t *testing.T) {
	empty := Monoid.Empty()
	assert.NotNil(t, empty)
	assert.Equal(t, 0, len(empty))
}

// TestMonoidConcat tests concatenating two header maps
func TestMonoidConcat(t *testing.T) {
	h1 := make(http.Header)
	h1.Set("X-Custom-1", "value1")
	h1.Set("Authorization", "Bearer token1")

	h2 := make(http.Header)
	h2.Set("X-Custom-2", "value2")
	h2.Set("Content-Type", "application/json")

	result := Monoid.Concat(h1, h2)

	assert.Equal(t, "value1", result.Get("X-Custom-1"))
	assert.Equal(t, "value2", result.Get("X-Custom-2"))
	assert.Equal(t, "Bearer token1", result.Get("Authorization"))
	assert.Equal(t, "application/json", result.Get("Content-Type"))
}

// TestMonoidConcatWithOverlap tests concatenating headers with overlapping keys
func TestMonoidConcatWithOverlap(t *testing.T) {
	h1 := make(http.Header)
	h1.Set("X-Custom", "value1")

	h2 := make(http.Header)
	h2.Add("X-Custom", "value2")

	result := Monoid.Concat(h1, h2)

	// Both values should be present
	values := result.Values("X-Custom")
	assert.Contains(t, values, "value1")
	assert.Contains(t, values, "value2")
}

// TestMonoidIdentity tests that concatenating with empty is identity
func TestMonoidIdentity(t *testing.T) {
	h := make(http.Header)
	h.Set("X-Test", "value")

	empty := Monoid.Empty()

	// Left identity: empty + h = h
	leftResult := Monoid.Concat(empty, h)
	assert.Equal(t, "value", leftResult.Get("X-Test"))

	// Right identity: h + empty = h
	rightResult := Monoid.Concat(h, empty)
	assert.Equal(t, "value", rightResult.Get("X-Test"))
}

// TestAtValuesGet tests getting header values using AtValues lens
func TestAtValuesGet(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Add("Accept", "application/json")
	headers.Add("Accept", "text/html")

	// Get Content-Type values
	ctLens := AtValues("Content-Type")
	ctValuesOpt := ctLens.Get(headers)
	assert.True(t, O.IsSome(ctValuesOpt))
	ctValues := O.GetOrElse(F.Constant([]string{}))(ctValuesOpt)
	assert.Equal(t, []string{"application/json"}, ctValues)

	// Get Accept values (multiple)
	acceptLens := AtValues("Accept")
	acceptValuesOpt := acceptLens.Get(headers)
	assert.True(t, O.IsSome(acceptValuesOpt))
	acceptValues := O.GetOrElse(F.Constant([]string{}))(acceptValuesOpt)
	assert.Equal(t, 2, len(acceptValues))
	assert.Contains(t, acceptValues, "application/json")
	assert.Contains(t, acceptValues, "text/html")
}

// TestAtValuesSet tests setting header values using AtValues lens
func TestAtValuesSet(t *testing.T) {
	headers := make(http.Header)
	headers.Set("X-Old", "old-value")

	lens := AtValues("Content-Type")
	newHeaders := lens.Set(O.Some([]string{"application/json", "text/plain"}))(headers)

	// New header should be set
	values := newHeaders.Values("Content-Type")
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, "application/json")
	assert.Contains(t, values, "text/plain")

	// Old header should still exist
	assert.Equal(t, "old-value", newHeaders.Get("X-Old"))
}

// TestAtValuesCanonical tests that header names are canonicalized
func TestAtValuesCanonical(t *testing.T) {
	headers := make(http.Header)
	headers.Set("content-type", "application/json")

	// Access with different casing
	lens := AtValues("Content-Type")
	valuesOpt := lens.Get(headers)

	assert.True(t, O.IsSome(valuesOpt))
	values := O.GetOrElse(F.Constant([]string{}))(valuesOpt)
	assert.Equal(t, []string{"application/json"}, values)
}

// TestAtValueGet tests getting first header value using AtValue lens
func TestAtValueGet(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer token123")

	lens := AtValue("Authorization")
	value := lens.Get(headers)

	assert.True(t, O.IsSome(value))
	token := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, "Bearer token123", token)
}

// TestAtValueGetNone tests getting non-existent header returns None
func TestAtValueGetNone(t *testing.T) {
	headers := make(http.Header)

	lens := AtValue("X-Non-Existent")
	value := lens.Get(headers)

	assert.True(t, O.IsNone(value))
}

// TestAtValueSet tests setting header value using AtValue lens
func TestAtValueSet(t *testing.T) {
	headers := make(http.Header)

	lens := AtValue("Content-Type")
	newHeaders := lens.Set(O.Some("application/json"))(headers)

	value := lens.Get(newHeaders)
	assert.True(t, O.IsSome(value))

	ct := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, "application/json", ct)
}

// TestAtValueSetNone tests removing header using AtValue lens
func TestAtValueSetNone(t *testing.T) {
	headers := make(http.Header)
	headers.Set("X-Custom", "value")

	lens := AtValue("X-Custom")
	newHeaders := lens.Set(O.None[string]())(headers)

	value := lens.Get(newHeaders)
	assert.True(t, O.IsNone(value))
}

// TestAtValueMultipleValues tests AtValue with multiple header values
func TestAtValueMultipleValues(t *testing.T) {
	headers := make(http.Header)
	headers.Add("Accept", "application/json")
	headers.Add("Accept", "text/html")

	lens := AtValue("Accept")
	value := lens.Get(headers)

	assert.True(t, O.IsSome(value))
	// Should get the first value
	first := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, "application/json", first)
}

// TestHeaderConstants tests that header constants are correct
func TestHeaderConstants(t *testing.T) {
	assert.Equal(t, "Accept", Accept)
	assert.Equal(t, "Authorization", Authorization)
	assert.Equal(t, "Content-Type", ContentType)
	assert.Equal(t, "Content-Length", ContentLength)
}

// TestHeaderConstantsUsage tests using header constants with http.Header
func TestHeaderConstantsUsage(t *testing.T) {
	headers := make(http.Header)

	headers.Set(Accept, "application/json")
	headers.Set(Authorization, "Bearer token")
	headers.Set(ContentType, "application/json")
	headers.Set(ContentLength, "1234")

	assert.Equal(t, "application/json", headers.Get(Accept))
	assert.Equal(t, "Bearer token", headers.Get(Authorization))
	assert.Equal(t, "application/json", headers.Get(ContentType))
	assert.Equal(t, "1234", headers.Get(ContentLength))
}

// TestAtValueWithConstants tests using AtValue with header constants
func TestAtValueWithConstants(t *testing.T) {
	headers := make(http.Header)
	headers.Set(ContentType, "application/json")

	lens := AtValue(ContentType)
	value := lens.Get(headers)

	assert.True(t, O.IsSome(value))
	ct := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, "application/json", ct)
}

// TestMonoidAssociativity tests that Monoid concatenation is associative
func TestMonoidAssociativity(t *testing.T) {
	h1 := make(http.Header)
	h1.Set("X-1", "value1")

	h2 := make(http.Header)
	h2.Set("X-2", "value2")

	h3 := make(http.Header)
	h3.Set("X-3", "value3")

	// (h1 + h2) + h3
	left := Monoid.Concat(Monoid.Concat(h1, h2), h3)

	// h1 + (h2 + h3)
	right := Monoid.Concat(h1, Monoid.Concat(h2, h3))

	// Both should have all three headers
	assert.Equal(t, "value1", left.Get("X-1"))
	assert.Equal(t, "value2", left.Get("X-2"))
	assert.Equal(t, "value3", left.Get("X-3"))

	assert.Equal(t, "value1", right.Get("X-1"))
	assert.Equal(t, "value2", right.Get("X-2"))
	assert.Equal(t, "value3", right.Get("X-3"))
}

// TestAtValuesEmptyHeader tests AtValues with empty headers
func TestAtValuesEmptyHeader(t *testing.T) {
	headers := make(http.Header)

	lens := AtValues("X-Non-Existent")
	valuesOpt := lens.Get(headers)

	assert.True(t, O.IsNone(valuesOpt))
}

// TestComplexHeaderOperations tests complex operations combining lenses and monoid
func TestComplexHeaderOperations(t *testing.T) {
	// Create initial headers
	h1 := make(http.Header)
	h1.Set("X-Initial", "initial")

	// Use lens to add Content-Type
	ctLens := AtValue(ContentType)
	h2 := ctLens.Set(O.Some("application/json"))(h1)

	// Use lens to add Authorization
	authLens := AtValue(Authorization)
	h3 := authLens.Set(O.Some("Bearer token"))(h2)

	// Create additional headers
	h4 := make(http.Header)
	h4.Set("X-Additional", "additional")

	// Combine using Monoid
	final := Monoid.Concat(h3, h4)

	// Verify all headers are present
	assert.Equal(t, "initial", final.Get("X-Initial"))
	assert.Equal(t, "application/json", final.Get(ContentType))
	assert.Equal(t, "Bearer token", final.Get(Authorization))
	assert.Equal(t, "additional", final.Get("X-Additional"))
}
