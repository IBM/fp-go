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

package http

import (
	"fmt"
	"io"
	"mime"
	H "net/http"
	"net/url"
	"regexp"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/record/generic"
)

type (
	// ParsedMediaType represents a parsed MIME media type as a Pair.
	// The first element is the media type string (e.g., "application/json"),
	// and the second element is a map of parameters (e.g., {"charset": "utf-8"}).
	//
	// Example:
	//   parsed := ParseMediaType("application/json; charset=utf-8")
	//   mediaType := P.Head(parsed)      // "application/json"
	//   params := P.Tail(parsed)         // map[string]string{"charset": "utf-8"}
	ParsedMediaType = P.Pair[string, map[string]string]

	// HttpError represents an HTTP error with detailed information about
	// the failed request. It includes the status code, response headers,
	// response body, and the URL that was accessed.
	//
	// This error type is created by StatusCodeError when an HTTP response
	// has a non-successful status code (not 2xx).
	//
	// Example:
	//   if httpErr, ok := err.(*HttpError); ok {
	//       fmt.Printf("Status: %d\n", httpErr.StatusCode())
	//       fmt.Printf("URL: %s\n", httpErr.URL())
	//       fmt.Printf("Body: %s\n", string(httpErr.Body()))
	//   }
	HttpError struct {
		statusCode int
		headers    H.Header
		body       []byte
		url        *url.URL
	}
)

var (
	// isJSONMimeType is a regex matcher that checks if a media type is a valid JSON type.
	// It matches "application/json" and variants like "application/vnd.api+json".
	isJSONMimeType = regexp.MustCompile(`application/(?:\w+\+)?json`).MatchString

	// ValidateResponse validates an HTTP response and returns an Either.
	// It checks if the response has a successful status code (2xx range).
	//
	// Returns:
	//   - Right(*http.Response) if status code is 2xx
	//   - Left(error) with HttpError if status code is not 2xx
	//
	// Example:
	//   result := ValidateResponse(response)
	//   E.Fold(
	//       func(err error) { /* handle error */ },
	//       func(resp *http.Response) { /* handle success */ },
	//   )(result)
	ValidateResponse = E.FromPredicate(isValidStatus, StatusCodeError)

	// validateJSONContentTypeString parses a content type string and validates
	// that it represents a valid JSON media type. This is an internal helper
	// used by ValidateJSONResponse.
	validateJSONContentTypeString = F.Flow2(
		ParseMediaType,
		E.ChainFirst(F.Flow2(
			P.Head[string, map[string]string],
			E.FromPredicate(isJSONMimeType, errors.OnSome[string]("mimetype [%s] is not a valid JSON content type")),
		)),
	)

	// ValidateJSONResponse validates that an HTTP response is a valid JSON response.
	// It checks both the status code (must be 2xx) and the Content-Type header
	// (must be a JSON media type like "application/json").
	//
	// Returns:
	//   - Right(*http.Response) if response is valid JSON with 2xx status
	//   - Left(error) if status is not 2xx or Content-Type is not JSON
	//
	// Example:
	//   result := ValidateJSONResponse(response)
	//   E.Fold(
	//       func(err error) { /* handle non-JSON or error response */ },
	//       func(resp *http.Response) { /* handle valid JSON response */ },
	//   )(result)
	ValidateJSONResponse = F.Flow2(
		E.Of[error, *H.Response],
		E.ChainFirst(F.Flow5(
			GetHeader,
			R.Lookup[H.Header](HeaderContentType),
			O.Chain(A.First[string]),
			E.FromOption[string](errors.OnNone("unable to access the [%s] header", HeaderContentType)),
			E.ChainFirst(validateJSONContentTypeString),
		)))

	// ValidateJsonResponse checks if an HTTP response is a valid JSON response.
	//
	// Deprecated: use ValidateJSONResponse instead (note the capitalization).
	ValidateJsonResponse = ValidateJSONResponse
)

const (
	// HeaderContentType is the standard HTTP Content-Type header name.
	// It indicates the media type of the resource or data being sent.
	//
	// Example values:
	//   - "application/json"
	//   - "text/html; charset=utf-8"
	//   - "application/xml"
	HeaderContentType = "Content-Type"
)

// ParseMediaType parses a MIME media type string into its components.
// It returns a ParsedMediaType (Pair) containing the media type and its parameters.
//
// Parameters:
//   - mediaType: A media type string (e.g., "application/json; charset=utf-8")
//
// Returns:
//   - Right(ParsedMediaType) with the parsed media type and parameters
//   - Left(error) if the media type string is invalid
//
// Example:
//
//	result := ParseMediaType("application/json; charset=utf-8")
//	E.Map(func(parsed ParsedMediaType) {
//	    mediaType := P.Head(parsed)  // "application/json"
//	    params := P.Tail(parsed)     // map[string]string{"charset": "utf-8"}
//	})(result)
func ParseMediaType(mediaType string) E.Either[error, ParsedMediaType] {
	m, p, err := mime.ParseMediaType(mediaType)
	return E.TryCatchError(P.MakePair(m, p), err)
}

// Error implements the error interface for HttpError.
// It returns a formatted error message including the status code and URL.
func (r *HttpError) Error() string {
	return fmt.Sprintf("invalid status code [%d] when accessing URL [%s]", r.statusCode, r.url)
}

// String returns the string representation of the HttpError.
// It's equivalent to calling Error().
func (r *HttpError) String() string {
	return r.Error()
}

// StatusCode returns the HTTP status code from the failed response.
//
// Example:
//
//	if httpErr, ok := err.(*HttpError); ok {
//	    code := httpErr.StatusCode()  // e.g., 404, 500
//	}
func (r *HttpError) StatusCode() int {
	return r.statusCode
}

// Headers returns a clone of the HTTP headers from the failed response.
// The headers are cloned to prevent modification of the original response.
//
// Example:
//
//	if httpErr, ok := err.(*HttpError); ok {
//	    headers := httpErr.Headers()
//	    contentType := headers.Get("Content-Type")
//	}
func (r *HttpError) Headers() H.Header {
	return r.headers
}

// URL returns the URL that was accessed when the error occurred.
//
// Example:
//
//	if httpErr, ok := err.(*HttpError); ok {
//	    url := httpErr.URL()
//	    fmt.Printf("Failed to access: %s\n", url)
//	}
func (r *HttpError) URL() *url.URL {
	return r.url
}

// Body returns the response body bytes from the failed response.
// This can be useful for debugging or displaying error messages from the server.
//
// Example:
//
//	if httpErr, ok := err.(*HttpError); ok {
//	    body := httpErr.Body()
//	    fmt.Printf("Error response: %s\n", string(body))
//	}
func (r *HttpError) Body() []byte {
	return r.body
}

// GetHeader extracts the HTTP headers from an http.Response.
// This is a functional accessor for the Header field.
//
// Parameters:
//   - resp: The HTTP response
//
// Returns:
//   - The http.Header map from the response
//
// Example:
//
//	headers := GetHeader(response)
//	contentType := headers.Get("Content-Type")
func GetHeader(resp *H.Response) H.Header {
	return resp.Header
}

// GetBody extracts the response body reader from an http.Response.
// This is a functional accessor for the Body field.
//
// Parameters:
//   - resp: The HTTP response
//
// Returns:
//   - The io.ReadCloser for reading the response body
//
// Example:
//
//	body := GetBody(response)
//	defer body.Close()
//	data, err := io.ReadAll(body)
func GetBody(resp *H.Response) io.ReadCloser {
	return resp.Body
}

// isValidStatus checks if an HTTP response has a successful status code.
// A status code is considered valid if it's in the 2xx range (200-299).
//
// Parameters:
//   - resp: The HTTP response to check
//
// Returns:
//   - true if status code is 2xx, false otherwise
func isValidStatus(resp *H.Response) bool {
	return resp.StatusCode >= H.StatusOK && resp.StatusCode < H.StatusMultipleChoices
}

// StatusCodeError creates an HttpError from an http.Response with a non-successful status code.
// It reads the response body and captures all relevant information for debugging.
//
// The function:
//   - Reads and stores the response body
//   - Clones the response headers
//   - Captures the request URL
//   - Creates a comprehensive error with all this information
//
// Parameters:
//   - resp: The HTTP response with a non-successful status code
//
// Returns:
//   - An error (specifically *HttpError) with detailed information
//
// Example:
//
//	if !isValidStatus(response) {
//	    err := StatusCodeError(response)
//	    return err
//	}
func StatusCodeError(resp *H.Response) error {
	// read the body
	bodyRdr := GetBody(resp)
	defer bodyRdr.Close()
	// try to access body content
	body, _ := io.ReadAll(bodyRdr)
	// return an error with comprehensive information
	return &HttpError{statusCode: resp.StatusCode, headers: GetHeader(resp).Clone(), body: body, url: resp.Request.URL}
}
