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
	"bytes"
	"io"
	H "net/http"
	"net/url"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NoError[A any](t *testing.T) func(E.Either[error, A]) bool {
	return E.Fold(func(err error) bool {
		return assert.NoError(t, err)
	}, F.Constant1[A](true))
}

func Error[A any](t *testing.T) func(E.Either[error, A]) bool {
	return E.Fold(F.Constant1[error](true), func(A) bool {
		return assert.Error(t, nil)
	})
}

func TestValidateJsonContentTypeString(t *testing.T) {
	res := F.Pipe1(
		validateJSONContentTypeString(C.JSON),
		NoError[ParsedMediaType](t),
	)
	assert.True(t, res)
}

func TestValidateInvalidJsonContentTypeString(t *testing.T) {
	res := F.Pipe1(
		validateJSONContentTypeString("application/xml"),
		Error[ParsedMediaType](t),
	)
	assert.True(t, res)
}

// TestParseMediaType tests parsing valid media types
func TestParseMediaType(t *testing.T) {
	tests := []struct {
		name      string
		mediaType string
		wantType  string
		wantParam map[string]string
	}{
		{
			name:      "simple JSON",
			mediaType: "application/json",
			wantType:  "application/json",
			wantParam: map[string]string{},
		},
		{
			name:      "JSON with charset",
			mediaType: "application/json; charset=utf-8",
			wantType:  "application/json",
			wantParam: map[string]string{"charset": "utf-8"},
		},
		{
			name:      "HTML with charset",
			mediaType: "text/html; charset=iso-8859-1",
			wantType:  "text/html",
			wantParam: map[string]string{"charset": "iso-8859-1"},
		},
		{
			name:      "multipart with boundary",
			mediaType: "multipart/form-data; boundary=----WebKitFormBoundary",
			wantType:  "multipart/form-data",
			wantParam: map[string]string{"boundary": "----WebKitFormBoundary"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseMediaType(tt.mediaType)
			require.True(t, E.IsRight(result), "ParseMediaType should succeed")

			parsed := E.GetOrElse(func(error) ParsedMediaType {
				return P.MakePair("", map[string]string{})
			})(result)
			mediaType := P.Head(parsed)
			params := P.Tail(parsed)

			assert.Equal(t, tt.wantType, mediaType)
			assert.Equal(t, tt.wantParam, params)
		})
	}
}

// TestParseMediaTypeInvalid tests parsing invalid media types
func TestParseMediaTypeInvalid(t *testing.T) {
	result := ParseMediaType("invalid media type")
	assert.True(t, E.IsLeft(result), "ParseMediaType should fail for invalid input")
}

// TestHttpErrorMethods tests all HttpError methods
func TestHttpErrorMethods(t *testing.T) {
	testURL, _ := url.Parse("https://example.com/api/test")
	headers := make(H.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Custom", "value")
	body := []byte(`{"error": "not found"}`)

	httpErr := &HttpError{
		statusCode: 404,
		headers:    headers,
		body:       body,
		url:        testURL,
	}

	// Test StatusCode
	assert.Equal(t, 404, httpErr.StatusCode())

	// Test Headers
	returnedHeaders := httpErr.Headers()
	assert.Equal(t, "application/json", returnedHeaders.Get("Content-Type"))
	assert.Equal(t, "value", returnedHeaders.Get("X-Custom"))

	// Test Body
	assert.Equal(t, body, httpErr.Body())
	assert.Equal(t, `{"error": "not found"}`, string(httpErr.Body()))

	// Test URL
	assert.Equal(t, testURL, httpErr.URL())
	assert.Equal(t, "https://example.com/api/test", httpErr.URL().String())

	// Test Error
	errMsg := httpErr.Error()
	assert.Contains(t, errMsg, "404")
	assert.Contains(t, errMsg, "https://example.com/api/test")

	// Test String
	assert.Equal(t, errMsg, httpErr.String())
}

// TestGetHeader tests the GetHeader function
func TestGetHeader(t *testing.T) {
	resp := &H.Response{
		Header: make(H.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	resp.Header.Set("Authorization", "Bearer token")

	headers := GetHeader(resp)
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, "Bearer token", headers.Get("Authorization"))
}

// TestGetBody tests the GetBody function
func TestGetBody(t *testing.T) {
	bodyContent := []byte("test body content")
	resp := &H.Response{
		Body: io.NopCloser(bytes.NewReader(bodyContent)),
	}

	body := GetBody(resp)
	defer body.Close()

	data, err := io.ReadAll(body)
	require.NoError(t, err)
	assert.Equal(t, bodyContent, data)
}

// TestIsValidStatus tests the isValidStatus function
func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"200 OK", H.StatusOK, true},
		{"201 Created", H.StatusCreated, true},
		{"204 No Content", H.StatusNoContent, true},
		{"299 (edge of 2xx)", 299, true},
		{"300 Multiple Choices", H.StatusMultipleChoices, false},
		{"301 Moved Permanently", H.StatusMovedPermanently, false},
		{"400 Bad Request", H.StatusBadRequest, false},
		{"404 Not Found", H.StatusNotFound, false},
		{"500 Internal Server Error", H.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &H.Response{StatusCode: tt.statusCode}
			assert.Equal(t, tt.want, isValidStatus(resp))
		})
	}
}

// TestValidateResponse tests the ValidateResponse function
func TestValidateResponse(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		resp := &H.Response{
			StatusCode: H.StatusOK,
			Header:     make(H.Header),
		}

		result := ValidateResponse(resp)
		assert.True(t, E.IsRight(result))

		validResp := E.GetOrElse(func(error) *H.Response { return nil })(result)
		assert.Equal(t, resp, validResp)
	})

	t.Run("error response", func(t *testing.T) {
		testURL, _ := url.Parse("https://example.com/test")
		resp := &H.Response{
			StatusCode: H.StatusNotFound,
			Header:     make(H.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte("not found"))),
			Request:    &H.Request{URL: testURL},
		}

		result := ValidateResponse(resp)
		assert.True(t, E.IsLeft(result))

		// Extract error using Fold
		var httpErr *HttpError
		E.Fold(
			func(err error) *H.Response {
				var ok bool
				httpErr, ok = err.(*HttpError)
				require.True(t, ok, "error should be *HttpError")
				return nil
			},
			func(r *H.Response) *H.Response { return r },
		)(result)
		assert.Equal(t, 404, httpErr.StatusCode())
	})
}

// TestStatusCodeError tests the StatusCodeError function
func TestStatusCodeError(t *testing.T) {
	testURL, _ := url.Parse("https://api.example.com/users/123")
	bodyContent := []byte(`{"error": "user not found"}`)

	headers := make(H.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Request-ID", "abc123")

	resp := &H.Response{
		StatusCode: H.StatusNotFound,
		Header:     headers,
		Body:       io.NopCloser(bytes.NewReader(bodyContent)),
		Request:    &H.Request{URL: testURL},
	}

	err := StatusCodeError(resp)
	require.Error(t, err)

	httpErr, ok := err.(*HttpError)
	require.True(t, ok, "error should be *HttpError")

	// Verify all fields
	assert.Equal(t, 404, httpErr.StatusCode())
	assert.Equal(t, testURL, httpErr.URL())
	assert.Equal(t, bodyContent, httpErr.Body())

	// Verify headers are cloned
	returnedHeaders := httpErr.Headers()
	assert.Equal(t, "application/json", returnedHeaders.Get("Content-Type"))
	assert.Equal(t, "abc123", returnedHeaders.Get("X-Request-ID"))

	// Verify error message
	errMsg := httpErr.Error()
	assert.Contains(t, errMsg, "404")
	assert.Contains(t, errMsg, "https://api.example.com/users/123")
}

// TestValidateJSONResponse tests the ValidateJSONResponse function
func TestValidateJSONResponse(t *testing.T) {
	t.Run("valid JSON response", func(t *testing.T) {
		resp := &H.Response{
			StatusCode: H.StatusOK,
			Header:     make(H.Header),
		}
		resp.Header.Set("Content-Type", "application/json")

		result := ValidateJSONResponse(resp)
		assert.True(t, E.IsRight(result), "should accept valid JSON response")
	})

	t.Run("JSON with charset", func(t *testing.T) {
		resp := &H.Response{
			StatusCode: H.StatusOK,
			Header:     make(H.Header),
		}
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")

		result := ValidateJSONResponse(resp)
		assert.True(t, E.IsRight(result), "should accept JSON with charset")
	})

	t.Run("JSON variant (hal+json)", func(t *testing.T) {
		resp := &H.Response{
			StatusCode: H.StatusOK,
			Header:     make(H.Header),
		}
		resp.Header.Set("Content-Type", "application/hal+json")

		result := ValidateJSONResponse(resp)
		assert.True(t, E.IsRight(result), "should accept JSON variants")
	})

	t.Run("non-JSON content type", func(t *testing.T) {
		resp := &H.Response{
			StatusCode: H.StatusOK,
			Header:     make(H.Header),
		}
		resp.Header.Set("Content-Type", "text/html")

		result := ValidateJSONResponse(resp)
		assert.True(t, E.IsLeft(result), "should reject non-JSON content type")
	})

	t.Run("missing Content-Type header", func(t *testing.T) {
		resp := &H.Response{
			StatusCode: H.StatusOK,
			Header:     make(H.Header),
		}

		result := ValidateJSONResponse(resp)
		assert.True(t, E.IsLeft(result), "should reject missing Content-Type")
	})

	t.Run("valid JSON with error status code", func(t *testing.T) {
		// Note: ValidateJSONResponse only validates Content-Type, not status code
		// It wraps the response in Right(response) first, then validates headers
		resp := &H.Response{
			StatusCode: H.StatusInternalServerError,
			Header:     make(H.Header),
		}
		resp.Header.Set("Content-Type", "application/json")

		result := ValidateJSONResponse(resp)
		// This actually succeeds because ValidateJSONResponse doesn't check status
		assert.True(t, E.IsRight(result), "ValidateJSONResponse only checks Content-Type, not status")
	})
}

// TestFullResponseAccessors tests Response and Body accessors
func TestFullResponseAccessors(t *testing.T) {
	resp := &H.Response{
		StatusCode: H.StatusOK,
		Header:     make(H.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	bodyContent := []byte(`{"message": "success"}`)
	fullResp := P.MakePair(resp, bodyContent)

	// Test Response accessor
	extractedResp := Response(fullResp)
	assert.Equal(t, resp, extractedResp)
	assert.Equal(t, H.StatusOK, extractedResp.StatusCode)

	// Test Body accessor
	extractedBody := Body(fullResp)
	assert.Equal(t, bodyContent, extractedBody)
	assert.Equal(t, `{"message": "success"}`, string(extractedBody))
}

// TestHeaderContentTypeConstant tests the HeaderContentType constant
func TestHeaderContentTypeConstant(t *testing.T) {
	assert.Equal(t, "Content-Type", HeaderContentType)

	// Test usage with http.Header
	headers := make(H.Header)
	headers.Set(HeaderContentType, "application/json")
	assert.Equal(t, "application/json", headers.Get(HeaderContentType))
}
