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
// The package defines the HTTP header names commonly used by applications as
// lower case constants (as required by HTTP/2 and HTTP/3), grouped into:
//   - Representation headers: ContentType, ContentLength, ContentEncoding, ...
//   - Request headers: Accept, Authorization, Cookie, UserAgent, Origin, ...
//   - Response headers: Location, SetCookie, WWWAuthenticate, RetryAfter, ...
//   - Caching and conditional requests: CacheControl, ETag, IfNoneMatch, ...
//   - CORS: AccessControlAllowOrigin, AccessControlAllowMethods, ...
//   - Security: StrictTransportSecurity, ContentSecurityPolicy, ...
//   - Proxy and tracing: XForwardedFor, XRequestID, Traceparent, ...
//
// Monoid:
//
// The Monoid provides a way to combine multiple http.Header maps:
//
//	headers1 := make(http.Header)
//	headers1.Set("X-Custom", "value1")
//
//	headers2 := make(http.Header)
//	headers2.Set(Authorization, "Bearer token")
//
//	combined := Monoid.Concat(headers1, headers2)
//	// combined now contains both headers
//
// Lenses:
//
// AtValues and AtValue provide lens-based access to header values:
//
//	// AtValues focuses on all values of a header (Option[[]string])
//	contentTypeLens := AtValues(ContentType)
//	values := contentTypeLens.Get(headers)
//
//	// AtValue focuses on the first value of a header (Option[string])
//	authLens := AtValue(Authorization)
//	token := authLens.Get(headers) // Returns Option[string]
//
// The lenses support functional updates:
//
//	// Set a header value
//	newHeaders := AtValue(ContentType).Set(O.Some(content.JSON))(headers)
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

// HTTP header name constants.
//
// These constants provide type-safe access to the HTTP header names that are
// commonly used by applications, grouped by purpose: representation headers,
// request headers, response headers, caching and conditional requests, CORS,
// security, proxy/forwarding and tracing.
//
// Header names are lower case as required by HTTP/2 (RFC 9113, Section 8.2.2)
// and HTTP/3 (RFC 9114, Section 4.2), which mandate that field names be
// transmitted in lower case. Go's http.Header canonicalizes keys on Get/Set,
// so these constants work transparently with HTTP/1.1 as well.
//
// The semantics of the standard headers are defined in RFC 9110 (HTTP
// Semantics) and RFC 9111 (HTTP Caching) unless noted otherwise.

// Representation and content headers, used in both requests and responses.
const (
	// ContentType indicates the media type of the message body, optionally
	// with parameters such as the character set.
	// Example: "content-type: application/json; charset=utf-8"
	ContentType = "content-type"

	// ContentLength indicates the size of the message body in bytes.
	// Example: "content-length: 348"
	ContentLength = "content-length"

	// ContentEncoding lists the encodings (typically compression) that have
	// been applied to the body, in the order they were applied.
	// Example: "content-encoding: gzip"
	ContentEncoding = "content-encoding"

	// ContentLanguage describes the natural language(s) of the intended
	// audience of the body.
	// Example: "content-language: en-US"
	ContentLanguage = "content-language"

	// ContentDisposition indicates whether the body should be displayed inline
	// or treated as an attachment, and may carry a suggested file name
	// (RFC 6266). It is also used for the parts of multipart/form-data bodies.
	// Example: "content-disposition: attachment; filename=\"report.pdf\""
	ContentDisposition = "content-disposition"

	// ContentRange indicates where a partial body belongs within the full
	// representation, used with 206 Partial Content responses.
	// Example: "content-range: bytes 200-1000/67589"
	ContentRange = "content-range"

	// ContentLocation indicates a URI that identifies the specific resource
	// corresponding to the representation in the message body.
	// Example: "content-location: /documents/foo.json"
	ContentLocation = "content-location"

	// TransferEncoding lists the transfer codings applied to the message for
	// safe transport. It is HTTP/1.1 only and not allowed in HTTP/2 and HTTP/3.
	// Example: "transfer-encoding: chunked"
	TransferEncoding = "transfer-encoding"

	// Date contains the date and time at which the message was originated.
	// Example: "date: Wed, 21 Oct 2015 07:28:00 GMT"
	Date = "date"

	// Link conveys one or more typed links to related resources (RFC 8288),
	// commonly used for pagination and resource hints.
	// Example: "link: <https://api.example.com/items?page=2>; rel=\"next\""
	Link = "link"
)

// Request headers, sent by the client to describe the request and its preferences.
const (
	// Accept specifies the media types that are acceptable for the response.
	// Example: "accept: application/json"
	Accept = "accept"

	// AcceptCharset specifies the character sets that are acceptable for the
	// response.
	// Example: "accept-charset: utf-8"
	AcceptCharset = "accept-charset"

	// AcceptEncoding specifies the content encodings (usually compression
	// algorithms) that the client can understand.
	// Example: "accept-encoding: gzip, deflate, br"
	AcceptEncoding = "accept-encoding"

	// AcceptLanguage specifies the natural languages preferred for the response.
	// Example: "accept-language: en-US,en;q=0.9,de;q=0.8"
	AcceptLanguage = "accept-language"

	// Authorization contains credentials for authenticating the client with
	// the server.
	// Example: "authorization: Bearer token123"
	Authorization = "authorization"

	// ProxyAuthorization contains credentials for authenticating the client
	// with an intermediate proxy.
	// Example: "proxy-authorization: Basic YWxhZGRpbjpvcGVuc2VzYW1l"
	ProxyAuthorization = "proxy-authorization"

	// Cookie contains cookies previously sent by the server via SetCookie
	// (RFC 6265).
	// Example: "cookie: session=abc123; theme=dark"
	Cookie = "cookie"

	// Host specifies the host and port of the target server. In HTTP/2 and
	// HTTP/3 it is replaced by the ":authority" pseudo-header.
	// Example: "host: api.example.com"
	Host = "host"

	// UserAgent identifies the client software originating the request.
	// Example: "user-agent: fp-go/2.0"
	UserAgent = "user-agent"

	// Referer contains the address of the resource from which the request was
	// initiated. The misspelling is part of the standard.
	// Example: "referer: https://example.com/page"
	Referer = "referer"

	// Origin indicates the scheme, host and port that caused the request. It is
	// used by CORS and for CSRF protection.
	// Example: "origin: https://example.com"
	Origin = "origin"

	// Range requests only part of a representation, typically a byte range.
	// Example: "range: bytes=200-1000"
	Range = "range"

	// Expect indicates behaviour the client expects from the server before the
	// body is sent.
	// Example: "expect: 100-continue"
	Expect = "expect"

	// Forwarded discloses information about the client and the proxies
	// involved in the request path (RFC 7239).
	// Example: "forwarded: for=192.0.2.60;proto=https;by=203.0.113.43"
	Forwarded = "forwarded"
)

// Response headers, sent by the server to describe the response.
const (
	// Location indicates the URL to redirect to (3xx responses) or the URL of
	// a newly created resource (201 Created).
	// Example: "location: /users/42"
	Location = "location"

	// Server describes the software used by the origin server.
	// Example: "server: nginx/1.25.3"
	Server = "server"

	// SetCookie sends a cookie from the server to the client (RFC 6265). It
	// may appear multiple times in one response.
	// Example: "set-cookie: session=abc123; Path=/; HttpOnly; Secure"
	SetCookie = "set-cookie"

	// WWWAuthenticate defines the authentication scheme(s) that should be used
	// to access the resource. It is sent with 401 Unauthorized responses.
	// Example: "www-authenticate: Bearer realm=\"api\""
	WWWAuthenticate = "www-authenticate"

	// ProxyAuthenticate defines the authentication scheme that should be used
	// to access a resource behind a proxy. It is sent with 407 responses.
	// Example: "proxy-authenticate: Basic realm=\"proxy\""
	ProxyAuthenticate = "proxy-authenticate"

	// RetryAfter indicates how long the client should wait before making a
	// follow-up request. It is used with 429, 503 and 3xx responses and holds
	// either a number of seconds or an HTTP date.
	// Example: "retry-after: 120"
	RetryAfter = "retry-after"

	// Allow lists the HTTP methods supported by the target resource. It is sent
	// with 405 Method Not Allowed and OPTIONS responses.
	// Example: "allow: GET, POST, HEAD"
	Allow = "allow"

	// AcceptRanges indicates whether the server supports range requests and in
	// which unit.
	// Example: "accept-ranges: bytes"
	AcceptRanges = "accept-ranges"
)

// Caching and conditional request headers (RFC 9111 and RFC 9110, Section 13).
const (
	// CacheControl holds directives that control caching in both requests and
	// responses.
	// Example: "cache-control: no-cache, max-age=0"
	CacheControl = "cache-control"

	// ETag is an opaque identifier for a specific version of a resource. It is
	// used for cache validation and optimistic concurrency control.
	// Example: "etag: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
	ETag = "etag"

	// LastModified contains the date and time at which the origin server
	// believes the resource was last modified.
	// Example: "last-modified: Wed, 21 Oct 2015 07:28:00 GMT"
	LastModified = "last-modified"

	// Expires contains the date and time after which the response is
	// considered stale. It is ignored if CacheControl contains max-age.
	// Example: "expires: Wed, 21 Oct 2015 07:28:00 GMT"
	Expires = "expires"

	// Age is the time in seconds the response has been stored in a cache.
	// Example: "age: 24"
	Age = "age"

	// Vary lists the request headers that were used to select the
	// representation, so caches can key responses correctly.
	// Example: "vary: accept-encoding, origin"
	Vary = "vary"

	// Pragma is a deprecated HTTP/1.0 caching header, kept for compatibility
	// with old caches. Use CacheControl instead.
	// Example: "pragma: no-cache"
	Pragma = "pragma"

	// IfMatch makes the request conditional: it succeeds only if the
	// resource's ETag matches one of the listed values. It is used to prevent
	// lost updates.
	// Example: "if-match: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
	IfMatch = "if-match"

	// IfNoneMatch makes the request conditional: it succeeds only if the
	// resource's ETag matches none of the listed values. It is used for cache
	// revalidation and yields 304 Not Modified when the resource is unchanged.
	// Example: "if-none-match: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
	IfNoneMatch = "if-none-match"

	// IfModifiedSince makes the request conditional: the resource is returned
	// only if it was modified after the given date.
	// Example: "if-modified-since: Wed, 21 Oct 2015 07:28:00 GMT"
	IfModifiedSince = "if-modified-since"

	// IfUnmodifiedSince makes the request conditional: it succeeds only if the
	// resource was not modified after the given date.
	// Example: "if-unmodified-since: Wed, 21 Oct 2015 07:28:00 GMT"
	IfUnmodifiedSince = "if-unmodified-since"

	// IfRange makes a range request conditional: the range is returned only if
	// the given ETag or date still matches, otherwise the full resource is sent.
	// Example: "if-range: \"33a64df551425fcc55e4d42a148795d9f25f89d4\""
	IfRange = "if-range"
)

// Cross-Origin Resource Sharing (CORS) headers, defined by the WHATWG Fetch standard.
const (
	// AccessControlAllowOrigin indicates which origin may access the response.
	// Example: "access-control-allow-origin: https://example.com"
	AccessControlAllowOrigin = "access-control-allow-origin"

	// AccessControlAllowMethods lists the methods allowed for the actual
	// request. It is sent in response to a preflight request.
	// Example: "access-control-allow-methods: GET, POST, PUT"
	AccessControlAllowMethods = "access-control-allow-methods"

	// AccessControlAllowHeaders lists the request headers allowed for the
	// actual request. It is sent in response to a preflight request.
	// Example: "access-control-allow-headers: content-type, authorization"
	AccessControlAllowHeaders = "access-control-allow-headers"

	// AccessControlAllowCredentials indicates whether the response may be
	// exposed when the request includes credentials such as cookies.
	// Example: "access-control-allow-credentials: true"
	AccessControlAllowCredentials = "access-control-allow-credentials"

	// AccessControlExposeHeaders lists the response headers that scripts in the
	// browser are allowed to read.
	// Example: "access-control-expose-headers: etag, link"
	AccessControlExposeHeaders = "access-control-expose-headers"

	// AccessControlMaxAge indicates how long, in seconds, the result of a
	// preflight request may be cached.
	// Example: "access-control-max-age: 600"
	AccessControlMaxAge = "access-control-max-age"

	// AccessControlRequestMethod is sent by the browser in a preflight request
	// to announce the method of the actual request.
	// Example: "access-control-request-method: PUT"
	AccessControlRequestMethod = "access-control-request-method"

	// AccessControlRequestHeaders is sent by the browser in a preflight request
	// to announce the headers of the actual request.
	// Example: "access-control-request-headers: content-type"
	AccessControlRequestHeaders = "access-control-request-headers"
)

// Security headers, sent by the server to protect clients against common attacks.
const (
	// StrictTransportSecurity instructs browsers to access the site only over
	// HTTPS for the given duration (HSTS, RFC 6797).
	// Example: "strict-transport-security: max-age=31536000; includeSubDomains"
	StrictTransportSecurity = "strict-transport-security"

	// ContentSecurityPolicy restricts the sources from which the browser may
	// load resources, mitigating cross-site scripting and injection attacks.
	// Example: "content-security-policy: default-src 'self'"
	ContentSecurityPolicy = "content-security-policy"

	// XContentTypeOptions with the value "nosniff" prevents browsers from
	// guessing a media type other than the declared ContentType.
	// Example: "x-content-type-options: nosniff"
	XContentTypeOptions = "x-content-type-options"

	// XFrameOptions controls whether the page may be rendered in a frame,
	// protecting against clickjacking. It is superseded by the CSP
	// frame-ancestors directive.
	// Example: "x-frame-options: DENY"
	XFrameOptions = "x-frame-options"

	// ReferrerPolicy controls how much referrer information is included in the
	// Referer header of subsequent requests.
	// Example: "referrer-policy: strict-origin-when-cross-origin"
	ReferrerPolicy = "referrer-policy"
)

// Proxy, forwarding and tracing headers. The x- headers are de-facto
// standards; Traceparent and Tracestate are defined by W3C Trace Context.
const (
	// XForwardedFor identifies the originating client IP address and the
	// proxies a request passed through. Prefer Forwarded for new systems.
	// Example: "x-forwarded-for: 203.0.113.195, 70.41.3.18"
	XForwardedFor = "x-forwarded-for"

	// XForwardedHost identifies the original host requested by the client
	// before a proxy rewrote it.
	// Example: "x-forwarded-host: example.com"
	XForwardedHost = "x-forwarded-host"

	// XForwardedProto identifies the protocol (http or https) the client used
	// to connect to a proxy or load balancer.
	// Example: "x-forwarded-proto: https"
	XForwardedProto = "x-forwarded-proto"

	// XRequestID carries a unique identifier for a request, used to correlate
	// log entries across services.
	// Example: "x-request-id: f058ebd6-02f7-4d3f-942e-904344e8cde5"
	XRequestID = "x-request-id"

	// XCorrelationID carries an identifier shared by all requests that belong
	// to the same logical operation across services.
	// Example: "x-correlation-id: 3f2c1a7e-9b8d-4c6f-a1e2-5d4b3c2a1f0e"
	XCorrelationID = "x-correlation-id"

	// Traceparent carries the distributed trace context (trace ID, parent span
	// ID and flags) as defined by W3C Trace Context.
	// Example: "traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	Traceparent = "traceparent"

	// Tracestate carries vendor-specific trace information alongside
	// Traceparent, as defined by W3C Trace Context.
	// Example: "tracestate: congo=t61rcWkgMzE"
	Tracestate = "tracestate"
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
	//   h2.Set(Authorization, "Bearer token")
	//
	//   combined := Monoid.Concat(h1, h2)
	//   // combined contains both X-Custom and Authorization headers
	Monoid = RG.UnionMonoid[http.Header](A.Semigroup[string]())

	// AtValues is a Lens that focuses on all values of a specific header.
	// It returns a lens that accesses the optional []string slice of header values,
	// None if the header is absent.
	// The header name is automatically canonicalized using MIME header key rules.
	//
	// Parameters:
	//   - name: The header name (will be canonicalized)
	//
	// Returns:
	//   - A Lens[http.Header, Option[[]string]] focusing on the header's values
	//
	// Example:
	//   lens := AtValues(ContentType)
	//   values := lens.Get(headers) // Returns Option[[]string]
	//   newHeaders := lens.Set(O.Some([]string{content.JSON}))(headers)
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
	//   lens := AtValue(Authorization)
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
