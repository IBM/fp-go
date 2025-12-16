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

// Package builder provides a functional, immutable HTTP request builder with composable operations.
// It follows functional programming principles to construct HTTP requests in a type-safe,
// testable, and maintainable way.
//
// The Builder type is immutable - all operations return a new builder instance rather than
// modifying the existing one. This ensures thread-safety and makes the code easier to reason about.
//
// Key Features:
//   - Immutable builder pattern with method chaining
//   - Lens-based access to builder properties
//   - Support for headers, query parameters, request body, and HTTP methods
//   - JSON and form data encoding
//   - URL construction with query parameter merging
//   - Hash generation for caching
//   - Bearer token authentication helpers
//
// Basic Usage:
//
//	import (
//	    B "github.com/IBM/fp-go/v2/http/builder"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Build a simple GET request
//	builder := F.Pipe2(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithHeader("Accept")("application/json"),
//	)
//
//	// Build a POST request with JSON body
//	builder := F.Pipe3(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithMethod("POST"),
//	    B.WithJSON(map[string]string{"name": "John"}),
//	)
//
//	// Build a request with query parameters
//	builder := F.Pipe3(
//	    B.Default,
//	    B.WithURL("https://api.example.com/search"),
//	    B.WithQueryArg("q")("golang"),
//	    B.WithQueryArg("limit")("10"),
//	)
//
// The package provides several convenience functions for common HTTP methods:
//   - WithGet, WithPost, WithPut, WithDelete for setting HTTP methods
//   - WithBearer for adding Bearer token authentication
//   - WithJSON for JSON payloads
//   - WithFormData for form-encoded payloads
//
// Lenses are provided for advanced use cases:
//   - URL, Method, Body, Headers, Query for accessing builder properties
//   - Header(name) for accessing individual headers
//   - QueryArg(name) for accessing individual query parameters
package builder

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"

	A "github.com/IBM/fp-go/v2/array"
	B "github.com/IBM/fp-go/v2/bytes"
	ENDO "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
	FM "github.com/IBM/fp-go/v2/http/form"
	H "github.com/IBM/fp-go/v2/http/headers"
	J "github.com/IBM/fp-go/v2/json"
	LZ "github.com/IBM/fp-go/v2/lazy"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/record"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
)

type (
	Builder struct {
		method  Option[string]
		url     string
		headers http.Header
		body    Option[Result[[]byte]]
		query   url.Values
	}

	// Endomorphism returns an [ENDO.Endomorphism] that transforms a builder
	Endomorphism = ENDO.Endomorphism[*Builder]
)

var (
	// Default is the default builder
	Default = &Builder{method: O.Some(defaultMethod()), headers: make(http.Header), body: noBody}

	defaultMethod = F.Constant(http.MethodGet)

	// Monoid is the [M.Monoid] for the [Endomorphism]
	Monoid = ENDO.Monoid[*Builder]()

	// Url is a [Lens] for the URL
	//
	// Deprecated: use [URL] instead
	Url = L.MakeLensRef((*Builder).GetURL, (*Builder).SetURL)
	// URL is a [Lens] for the URL
	URL = L.MakeLensRef((*Builder).GetURL, (*Builder).SetURL)
	// Method is a [Lens] for the HTTP method
	Method = L.MakeLensRef((*Builder).GetMethod, (*Builder).SetMethod)
	// Body is a [Lens] for the request body
	Body = L.MakeLensRef((*Builder).GetBody, (*Builder).SetBody)
	// Headers is a [Lens] for the complete set of request headers
	Headers = L.MakeLensRef((*Builder).GetHeaders, (*Builder).SetHeaders)
	// Query is a [Lens] for the set of query parameters
	Query = L.MakeLensRef((*Builder).GetQuery, (*Builder).SetQuery)

	rawQuery = L.MakeLensRef(getRawQuery, setRawQuery)

	getHeader = F.Bind2of2((*Builder).GetHeader)
	delHeader = F.Bind2of2((*Builder).DelHeader)
	setHeader = F.Bind2of3((*Builder).SetHeader)

	noHeader   = O.None[string]()
	noBody     = O.None[Result[[]byte]]()
	noQueryArg = O.None[string]()

	parseURL   = result.Eitherize1(url.Parse)
	parseQuery = result.Eitherize1(url.ParseQuery)

	// WithQuery creates a [Endomorphism] for a complete set of query parameters
	WithQuery = Query.Set
	// WithMethod creates a [Endomorphism] for a certain method
	WithMethod = Method.Set
	// WithUrl creates a [Endomorphism] for the URL
	//
	// Deprecated: use [WithURL] instead
	WithUrl = URL.Set
	// WithURL creates a [Endomorphism] for the URL
	WithURL = URL.Set
	// WithHeaders creates a [Endomorphism] for a set of headers
	WithHeaders = Headers.Set
	// WithBody creates a [Endomorphism] for a request body
	WithBody = F.Flow2(
		O.Of[Result[[]byte]],
		Body.Set,
	)
	// WithBytes creates a [Endomorphism] for a request body using bytes
	WithBytes = F.Flow2(
		result.Of[[]byte],
		WithBody,
	)
	// WithContentType adds the [H.ContentType] header
	WithContentType = WithHeader(H.ContentType)
	// WithAuthorization adds the [H.Authorization] header
	WithAuthorization = WithHeader(H.Authorization)

	// WithGet adds the [http.MethodGet] method
	WithGet = WithMethod(http.MethodGet)
	// WithPost adds the [http.MethodPost] method
	WithPost = WithMethod(http.MethodPost)
	// WithPut adds the [http.MethodPut] method
	WithPut = WithMethod(http.MethodPut)
	// WithDelete adds the [http.MethodDelete] method
	WithDelete = WithMethod(http.MethodDelete)

	// WithBearer creates a [Endomorphism] to add a Bearer [H.Authorization] header
	WithBearer = F.Flow2(
		S.Format[string]("Bearer %s"),
		WithAuthorization,
	)

	// WithoutBody creates a [Endomorphism] to remove the body
	WithoutBody = F.Pipe1(
		noBody,
		Body.Set,
	)

	// WithFormData creates a [Endomorphism] to send form data payload
	WithFormData = F.Flow4(
		url.Values.Encode,
		S.ToBytes,
		WithBytes,
		ENDO.Chain(WithContentType(C.FormEncoded)),
	)

	// bodyAsBytes returns a []byte with a fallback to the empty array
	bodyAsBytes = O.Fold(B.Empty, result.Fold(F.Ignore1of1[error](B.Empty), F.Identity[[]byte]))
)

func setRawQuery(u *url.URL, raw string) *url.URL {
	u.RawQuery = raw
	return u
}

func getRawQuery(u *url.URL) string {
	return u.RawQuery
}

func (builder *Builder) clone() *Builder {
	cpy := *builder
	cpy.headers = cpy.headers.Clone()
	return &cpy
}

// GetTargetUrl constructs a full URL with query parameters on top of the provided URL string
//
// Deprecated: use [GetTargetURL] instead
func (builder *Builder) GetTargetUrl() Result[string] {
	return builder.GetTargetURL()
}

// GetTargetURL constructs a full URL with query parameters on top of the provided URL string
func (builder *Builder) GetTargetURL() Result[string] {
	// construct the final URL
	return F.Pipe3(
		builder,
		Url.Get,
		parseURL,
		result.Chain(F.Flow4(
			T.Replicate2[*url.URL],
			T.Map2(
				F.Flow2(
					F.Curry2(setRawQuery),
					result.Of[func(string) *url.URL],
				),
				F.Flow3(
					rawQuery.Get,
					parseQuery,
					result.Map(F.Flow2(
						F.Curry2(FM.ValuesMonoid.Concat)(builder.GetQuery()),
						url.Values.Encode,
					)),
				),
			),
			T.Tupled2(result.MonadAp[*url.URL, string]),
			result.Map((*url.URL).String),
		)),
	)
}

// Deprecated: use [GetURL] instead
func (builder *Builder) GetUrl() string {
	return builder.url
}

func (builder *Builder) GetURL() string {
	return builder.url
}

func (builder *Builder) GetMethod() string {
	return F.Pipe1(
		builder.method,
		O.GetOrElse(defaultMethod),
	)
}

func (builder *Builder) GetHeaders() http.Header {
	return builder.headers
}

func (builder *Builder) GetQuery() url.Values {
	return builder.query
}

func (builder *Builder) SetQuery(query url.Values) *Builder {
	builder.query = query
	return builder
}

func (builder *Builder) GetBody() Option[Result[[]byte]] {
	return builder.body
}

func (builder *Builder) SetMethod(method string) *Builder {
	builder.method = O.Some(method)
	return builder
}

// Deprecated: use [SetURL] instead
func (builder *Builder) SetUrl(url string) *Builder {
	builder.url = url
	return builder
}

func (builder *Builder) SetURL(url string) *Builder {
	builder.url = url
	return builder
}

func (builder *Builder) SetHeaders(headers http.Header) *Builder {
	builder.headers = headers
	return builder
}

func (builder *Builder) SetBody(body Option[Result[[]byte]]) *Builder {
	builder.body = body
	return builder
}

func (builder *Builder) SetHeader(name, value string) *Builder {
	builder.headers.Set(name, value)
	return builder
}

func (builder *Builder) DelHeader(name string) *Builder {
	builder.headers.Del(name)
	return builder
}

func (builder *Builder) GetHeader(name string) Option[string] {
	return F.Pipe2(
		name,
		builder.headers.Get,
		O.FromPredicate(S.IsNonEmpty),
	)
}

func (builder *Builder) GetHeaderValues(name string) []string {
	return builder.headers.Values(name)
}

// GetHash returns a hash value for the builder that can be used as a cache key
func (builder *Builder) GetHash() string {
	return MakeHash(builder)
}

// Header returns a [Lens] for a single header
func Header(name string) Lens[*Builder, Option[string]] {
	get := getHeader(name)
	set := F.Bind1of2(setHeader(name))
	del := F.Flow2(
		LZ.Of[*Builder],
		LZ.Map(delHeader(name)),
	)

	return L.MakeLensWithName(get, func(b *Builder, value Option[string]) *Builder {
		cpy := b.clone()
		return F.Pipe1(
			value,
			O.Fold(del(cpy), set(cpy)),
		)
	}, fmt.Sprintf("HttpHeader[%s]", name))
}

// WithHeader creates a [Endomorphism] for a certain header
func WithHeader(name string) func(value string) Endomorphism {
	return F.Flow2(
		O.Of[string],
		Header(name).Set,
	)
}

// WithoutHeader creates a [Endomorphism] to remove a certain header
func WithoutHeader(name string) Endomorphism {
	return Header(name).Set(noHeader)
}

// WithJson creates a [Endomorphism] to send JSON payload
//
// Deprecated: use [WithJSON] instead
func WithJson[T any](data T) Endomorphism {
	return WithJSON(data)
}

// WithJSON creates a [Endomorphism] to send JSON payload
func WithJSON[T any](data T) Endomorphism {
	return Monoid.Concat(
		F.Pipe2(
			data,
			J.Marshal[T],
			WithBody,
		),
		WithContentType(C.JSON),
	)
}

// QueryArg is a [Lens] for the first value of a query argument
func QueryArg(name string) Lens[*Builder, Option[string]] {
	return F.Pipe1(
		Query,
		L.Compose[*Builder](FM.AtValue(name)),
	)
}

// WithQueryArg creates a [Endomorphism] for a certain query argument
func WithQueryArg(name string) func(value string) Endomorphism {
	return F.Flow2(
		O.Of[string],
		QueryArg(name).Set,
	)
}

// WithoutQueryArg creates a [Endomorphism] that removes a query argument
func WithoutQueryArg(name string) Endomorphism {
	return QueryArg(name).Set(noQueryArg)
}

func hashWriteValue(buf *bytes.Buffer, value string) *bytes.Buffer {
	buf.WriteString(value)
	return buf
}

func hashWriteQuery(name string, buf *bytes.Buffer, values []string) *bytes.Buffer {
	buf.WriteString(name)
	return A.Reduce(hashWriteValue, buf)(values)
}

func makeBytes(b *Builder) []byte {
	var buf bytes.Buffer

	buf.WriteString(b.GetMethod())
	buf.WriteString(b.GetURL())
	b.GetHeaders().Write(&buf) // #nosec: G104

	R.ReduceOrdWithIndex[[]string, *bytes.Buffer](S.Ord)(hashWriteQuery, &buf)(b.GetQuery())

	buf.Write(bodyAsBytes(b.GetBody()))

	return buf.Bytes()
}

// MakeHash converts a [Builder] into a hash string, convenient to use as a cache key
func MakeHash(b *Builder) string {
	return fmt.Sprintf("%x", sha256.Sum256(makeBytes(b)))
}
