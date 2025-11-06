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
	ParsedMediaType = P.Pair[string, map[string]string]

	HttpError struct {
		statusCode int
		headers    H.Header
		body       []byte
		url        *url.URL
	}
)

var (
	// mime type to check if a media type matches
	isJSONMimeType = regexp.MustCompile(`application/(?:\w+\+)?json`).MatchString
	// ValidateResponse validates an HTTP response and returns an [E.Either] if the response is not a success
	ValidateResponse = E.FromPredicate(isValidStatus, StatusCodeError)
	// alidateJsonContentTypeString parses a content type a validates that it is valid JSON
	validateJSONContentTypeString = F.Flow2(
		ParseMediaType,
		E.ChainFirst(F.Flow2(
			P.Head[string, map[string]string],
			E.FromPredicate(isJSONMimeType, errors.OnSome[string]("mimetype [%s] is not a valid JSON content type")),
		)),
	)
	// ValidateJSONResponse checks if an HTTP response is a valid JSON response
	ValidateJSONResponse = F.Flow2(
		E.Of[error, *H.Response],
		E.ChainFirst(F.Flow5(
			GetHeader,
			R.Lookup[H.Header](HeaderContentType),
			O.Chain(A.First[string]),
			E.FromOption[string](errors.OnNone("unable to access the [%s] header", HeaderContentType)),
			E.ChainFirst(validateJSONContentTypeString),
		)))
	// ValidateJsonResponse checks if an HTTP response is a valid JSON response
	//
	// Deprecated: use [ValidateJSONResponse] instead
	ValidateJsonResponse = ValidateJSONResponse
)

const (
	HeaderContentType = "Content-Type"
)

// ParseMediaType parses a media type into a tuple
func ParseMediaType(mediaType string) E.Either[error, ParsedMediaType] {
	m, p, err := mime.ParseMediaType(mediaType)
	return E.TryCatchError(P.MakePair(m, p), err)
}

// Error fulfills the error interface
func (r *HttpError) Error() string {
	return fmt.Sprintf("invalid status code [%d] when accessing URL [%s]", r.statusCode, r.url)
}

func (r *HttpError) String() string {
	return r.Error()
}

func (r *HttpError) StatusCode() int {
	return r.statusCode
}

func (r *HttpError) Headers() H.Header {
	return r.headers
}

func (r *HttpError) URL() *url.URL {
	return r.url
}

func (r *HttpError) Body() []byte {
	return r.body
}

func GetHeader(resp *H.Response) H.Header {
	return resp.Header
}

func GetBody(resp *H.Response) io.ReadCloser {
	return resp.Body
}

func isValidStatus(resp *H.Response) bool {
	return resp.StatusCode >= H.StatusOK && resp.StatusCode < H.StatusMultipleChoices
}

// StatusCodeError creates an instance of [HttpError] filled with information from the response
func StatusCodeError(resp *H.Response) error {
	// read the body
	bodyRdr := GetBody(resp)
	defer bodyRdr.Close()
	// try to access body content
	body, _ := io.ReadAll(bodyRdr)
	// return an error with comprehensive information
	return &HttpError{statusCode: resp.StatusCode, headers: GetHeader(resp).Clone(), body: body, url: resp.Request.URL}
}
