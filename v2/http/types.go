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

// Package http provides functional programming utilities for working with HTTP
// requests and responses. It offers type-safe abstractions, validation functions,
// and utilities for handling HTTP operations in a functional style.
//
// The package includes:
//   - Type definitions for HTTP responses with bodies
//   - Validation functions for HTTP responses
//   - JSON content type validation
//   - Error handling with detailed HTTP error information
//   - Functional utilities for accessing response components
//
// Types:
//
// FullResponse represents a complete HTTP response including both the response
// object and the body as a byte array. It's implemented as a Pair for functional
// composition:
//
//	type FullResponse = Pair[*http.Response, []byte]
//
// The Response and Body functions provide lens-like access to the components:
//
//	resp := Response(fullResponse)  // Get *http.Response
//	body := Body(fullResponse)      // Get []byte
//
// Validation:
//
// ValidateResponse checks if an HTTP response has a successful status code (2xx):
//
//	result := ValidateResponse(response)
//	// Returns Either[error, *http.Response]
//
// ValidateJSONResponse validates both the status code and Content-Type header:
//
//	result := ValidateJSONResponse(response)
//	// Returns Either[error, *http.Response]
//
// Error Handling:
//
// HttpError provides detailed information about HTTP failures:
//
//	err := StatusCodeError(response)
//	if httpErr, ok := err.(*HttpError); ok {
//	    code := httpErr.StatusCode()
//	    headers := httpErr.Headers()
//	    body := httpErr.Body()
//	    url := httpErr.URL()
//	}
package http

import (
	H "net/http"

	P "github.com/IBM/fp-go/v2/pair"
)

type (
	// FullResponse represents a complete HTTP response including both the
	// *http.Response object and the response body as a byte slice.
	//
	// It's implemented as a Pair to enable functional composition and
	// transformation of HTTP responses. This allows you to work with both
	// the response metadata (status, headers) and body content together.
	//
	// Example:
	//   fullResp := MakePair(response, bodyBytes)
	//   resp := Response(fullResp)  // Extract *http.Response
	//   body := Body(fullResp)      // Extract []byte
	FullResponse = P.Pair[*H.Response, []byte]
)

var (
	// Response is a lens-like accessor that extracts the *http.Response
	// from a FullResponse. It provides functional access to the response
	// metadata including status code, headers, and other HTTP response fields.
	//
	// Example:
	//   fullResp := MakePair(response, bodyBytes)
	//   resp := Response(fullResp)
	//   statusCode := resp.StatusCode
	Response = P.Head[*H.Response, []byte]

	// Body is a lens-like accessor that extracts the response body bytes
	// from a FullResponse. It provides functional access to the raw body
	// content without needing to read from an io.Reader.
	//
	// Example:
	//   fullResp := MakePair(response, bodyBytes)
	//   body := Body(fullResp)
	//   content := string(body)
	Body = P.Tail[*H.Response, []byte]

	// FromResponse creates a function that constructs a FullResponse from
	// a given *http.Response. It returns a function that takes a body byte
	// slice and combines it with the response to create a FullResponse.
	//
	// This is useful for functional composition where you want to partially
	// apply the response and later provide the body.
	//
	// Example:
	//   makeFullResp := FromResponse(response)
	//   fullResp := makeFullResp(bodyBytes)
	FromResponse = P.FromHead[[]byte, *H.Response]

	// FromBody creates a function that constructs a FullResponse from
	// a given body byte slice. It returns a function that takes an
	// *http.Response and combines it with the body to create a FullResponse.
	//
	// This is useful for functional composition where you want to partially
	// apply the body and later provide the response.
	//
	// Example:
	//   makeFullResp := FromBody(bodyBytes)
	//   fullResp := makeFullResp(response)
	FromBody = P.FromTail[*H.Response, []byte]
)
