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
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
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
	name := ContentType
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
	h1.Set(Authorization, "Bearer token1")

	h2 := make(http.Header)
	h2.Set("X-Custom-2", "value2")
	h2.Set(ContentType, C.JSON)

	result := Monoid.Concat(h1, h2)

	assert.Equal(t, "value1", result.Get("X-Custom-1"))
	assert.Equal(t, "value2", result.Get("X-Custom-2"))
	assert.Equal(t, "Bearer token1", result.Get(Authorization))
	assert.Equal(t, C.JSON, result.Get(ContentType))
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
	headers.Set(ContentType, C.JSON)
	headers.Add(Accept, C.JSON)
	headers.Add(Accept, C.TextHTML)

	// Get Content-Type values
	ctLens := AtValues(ContentType)
	ctValuesOpt := ctLens.Get(headers)
	assert.True(t, O.IsSome(ctValuesOpt))
	ctValues := O.GetOrElse(F.Constant([]string{}))(ctValuesOpt)
	assert.Equal(t, []string{C.JSON}, ctValues)

	// Get Accept values (multiple)
	acceptLens := AtValues(Accept)
	acceptValuesOpt := acceptLens.Get(headers)
	assert.True(t, O.IsSome(acceptValuesOpt))
	acceptValues := O.GetOrElse(F.Constant([]string{}))(acceptValuesOpt)
	assert.Equal(t, 2, len(acceptValues))
	assert.Contains(t, acceptValues, C.JSON)
	assert.Contains(t, acceptValues, C.TextHTML)
}

// TestAtValuesSet tests setting header values using AtValues lens
func TestAtValuesSet(t *testing.T) {
	headers := make(http.Header)
	headers.Set("X-Old", "old-value")

	lens := AtValues(ContentType)
	newHeaders := lens.Set(O.Some([]string{C.JSON, C.TextPlain}))(headers)

	// New header should be set
	values := newHeaders.Values(ContentType)
	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, C.JSON)
	assert.Contains(t, values, C.TextPlain)

	// Old header should still exist
	assert.Equal(t, "old-value", newHeaders.Get("X-Old"))
}

// TestAtValuesCanonical tests that header names are canonicalized
func TestAtValuesCanonical(t *testing.T) {
	headers := make(http.Header)
	headers.Set("content-type", C.JSON)

	// Access with different casing
	lens := AtValues("Content-Type")
	valuesOpt := lens.Get(headers)

	assert.True(t, O.IsSome(valuesOpt))
	values := O.GetOrElse(F.Constant([]string{}))(valuesOpt)
	assert.Equal(t, []string{C.JSON}, values)
}

// TestAtValueGet tests getting first header value using AtValue lens
func TestAtValueGet(t *testing.T) {
	headers := make(http.Header)
	headers.Set(Authorization, "Bearer token123")

	lens := AtValue(Authorization)
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

	lens := AtValue(ContentType)
	newHeaders := lens.Set(O.Some(C.JSON))(headers)

	value := lens.Get(newHeaders)
	assert.True(t, O.IsSome(value))

	ct := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, C.JSON, ct)
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
	headers.Add(Accept, C.JSON)
	headers.Add(Accept, C.TextHTML)

	lens := AtValue(Accept)
	value := lens.Get(headers)

	assert.True(t, O.IsSome(value))
	// Should get the first value
	first := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, C.JSON, first)
}

// TestHeaderConstants tests that header constants are correct
func TestHeaderConstants(t *testing.T) {
	assert.Equal(t, "accept", Accept)
	assert.Equal(t, "authorization", Authorization)
	assert.Equal(t, "content-type", ContentType)
	assert.Equal(t, "content-length", ContentLength)
}

// allHeaderConstants lists every header name constant defined by the package
var allHeaderConstants = []string{
	ContentType, ContentLength, ContentEncoding, ContentLanguage, ContentDisposition,
	ContentRange, ContentLocation, TransferEncoding, Date, Link,
	Accept, AcceptCharset, AcceptEncoding, AcceptLanguage, Authorization,
	ProxyAuthorization, Cookie, Host, UserAgent, Referer, Origin, Range, Expect, Forwarded,
	Location, Server, SetCookie, WWWAuthenticate, ProxyAuthenticate, RetryAfter, Allow, AcceptRanges,
	CacheControl, ETag, LastModified, Expires, Age, Vary, Pragma,
	IfMatch, IfNoneMatch, IfModifiedSince, IfUnmodifiedSince, IfRange,
	AccessControlAllowOrigin, AccessControlAllowMethods, AccessControlAllowHeaders,
	AccessControlAllowCredentials, AccessControlExposeHeaders, AccessControlMaxAge,
	AccessControlRequestMethod, AccessControlRequestHeaders,
	StrictTransportSecurity, ContentSecurityPolicy, XContentTypeOptions, XFrameOptions, ReferrerPolicy,
	XForwardedFor, XForwardedHost, XForwardedProto, XRequestID, XCorrelationID, Traceparent, Tracestate,
}

// TestAllHeaderConstantsLowerCase verifies that every header name is lower case
// (HTTP/2 and HTTP/3 requirement), unique, and usable with http.Header and the lenses
func TestAllHeaderConstantsLowerCase(t *testing.T) {
	seen := make(map[string]bool)
	for _, name := range allHeaderConstants {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, strings.ToLower(name), name, "header name must be lower case")
			assert.False(t, seen[name], "header name must be unique")
			seen[name] = true

			headers := make(http.Header)
			headers.Set(name, "value")
			assert.Equal(t, "value", headers.Get(name))
			assert.Equal(t, O.Some("value"), AtValue(name).Get(headers))
		})
	}
}

// TestHeaderConstantsUsage tests using header constants with http.Header
func TestHeaderConstantsUsage(t *testing.T) {
	headers := make(http.Header)

	headers.Set(Accept, C.JSON)
	headers.Set(Authorization, "Bearer token")
	headers.Set(ContentType, C.JSON)
	headers.Set(ContentLength, "1234")

	assert.Equal(t, C.JSON, headers.Get(Accept))
	assert.Equal(t, "Bearer token", headers.Get(Authorization))
	assert.Equal(t, C.JSON, headers.Get(ContentType))
	assert.Equal(t, "1234", headers.Get(ContentLength))
}

// TestAtValueWithConstants tests using AtValue with header constants
func TestAtValueWithConstants(t *testing.T) {
	headers := make(http.Header)
	headers.Set(ContentType, C.JSON)

	lens := AtValue(ContentType)
	value := lens.Get(headers)

	assert.True(t, O.IsSome(value))
	ct := O.GetOrElse(F.Constant(""))(value)
	assert.Equal(t, C.JSON, ct)
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
	h2 := ctLens.Set(O.Some(C.JSON))(h1)

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
	assert.Equal(t, C.JSON, final.Get(ContentType))
	assert.Equal(t, "Bearer token", final.Get(Authorization))
	assert.Equal(t, "additional", final.Get("X-Additional"))
}
