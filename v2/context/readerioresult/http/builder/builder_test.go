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

package builder

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/http/builder"
	IO "github.com/IBM/fp-go/v2/io"
	"github.com/stretchr/testify/assert"
)

func TestBuilderWithQuery(t *testing.T) {
	// add some query
	withLimit := R.WithQueryArg("limit")("10")
	withURL := R.WithURL("http://www.example.org?a=b")

	b := F.Pipe2(
		R.Default,
		withLimit,
		withURL,
	)

	req := F.Pipe3(
		b,
		Requester,
		RIOE.Map(func(r *http.Request) *url.URL {
			return r.URL
		}),
		RIOE.ChainFirstIOK(func(u *url.URL) IO.IO[Void] {
			return IO.FromImpure(func() {
				q := u.Query()
				assert.Equal(t, "10", q.Get("limit"))
				assert.Equal(t, "b", q.Get("a"))
			})
		}),
	)

	assert.True(t, E.IsRight(req(t.Context())()))
}

// TestBuilderWithoutBody tests creating a request without a body
func TestBuilderWithoutBody(t *testing.T) {
	builder := F.Pipe2(
		R.Default,
		R.WithURL("https://api.example.com/users"),
		R.WithMethod("GET"),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://api.example.com/users", req.URL.String())
	assert.Nil(t, req.Body, "Expected nil body for GET request")
}

// TestBuilderWithBody tests creating a request with a body
func TestBuilderWithBody(t *testing.T) {
	bodyData := []byte(`{"name":"John","age":30}`)

	builder := F.Pipe3(
		R.Default,
		R.WithURL("https://api.example.com/users"),
		R.WithMethod("POST"),
		R.WithBytes(bodyData),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https://api.example.com/users", req.URL.String())
	assert.NotNil(t, req.Body, "Expected non-nil body for POST request")
	assert.Equal(t, "24", req.Header.Get("Content-Length"))
}

// TestBuilderWithHeaders tests that headers are properly set
func TestBuilderWithHeaders(t *testing.T) {
	builder := F.Pipe3(
		R.Default,
		R.WithURL("https://api.example.com/data"),
		R.WithHeader("Authorization")("Bearer token123"),
		R.WithHeader("Accept")("application/json"),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "Bearer token123", req.Header.Get("Authorization"))
	assert.Equal(t, "application/json", req.Header.Get("Accept"))
}

// TestBuilderWithInvalidURL tests error handling for invalid URLs
func TestBuilderWithInvalidURL(t *testing.T) {
	builder := F.Pipe1(
		R.Default,
		R.WithURL("://invalid-url"),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsLeft(result), "Expected Left result for invalid URL")
}

// TestBuilderWithEmptyMethod tests creating a request with empty method
func TestBuilderWithEmptyMethod(t *testing.T) {
	builder := F.Pipe2(
		R.Default,
		R.WithURL("https://api.example.com/users"),
		R.WithMethod(""),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	// Empty method should still work (defaults to GET in http.NewRequest)
	assert.True(t, E.IsRight(result), "Expected Right result")
}

// TestBuilderWithMultipleHeaders tests setting multiple headers
func TestBuilderWithMultipleHeaders(t *testing.T) {
	builder := F.Pipe4(
		R.Default,
		R.WithURL("https://api.example.com/data"),
		R.WithHeader("X-Custom-Header-1")("value1"),
		R.WithHeader("X-Custom-Header-2")("value2"),
		R.WithHeader("X-Custom-Header-3")("value3"),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "value1", req.Header.Get("X-Custom-Header-1"))
	assert.Equal(t, "value2", req.Header.Get("X-Custom-Header-2"))
	assert.Equal(t, "value3", req.Header.Get("X-Custom-Header-3"))
}

// TestBuilderWithBodyAndHeaders tests combining body and headers
func TestBuilderWithBodyAndHeaders(t *testing.T) {
	bodyData := []byte(`{"test":"data"}`)

	builder := F.Pipe4(
		R.Default,
		R.WithURL("https://api.example.com/submit"),
		R.WithMethod("PUT"),
		R.WithBytes(bodyData),
		R.WithHeader("X-Request-ID")("12345"),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "PUT", req.Method)
	assert.NotNil(t, req.Body, "Expected non-nil body")
	assert.Equal(t, "12345", req.Header.Get("X-Request-ID"))
	assert.Equal(t, "15", req.Header.Get("Content-Length"))
}

// TestBuilderContextCancellation tests that context cancellation is respected
func TestBuilderContextCancellation(t *testing.T) {
	builder := F.Pipe1(
		R.Default,
		R.WithURL("https://api.example.com/users"),
	)

	requester := Requester(builder)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // Cancel immediately

	result := requester(ctx)()

	// The request should still be created (cancellation affects execution, not creation)
	// But we verify the context is properly passed
	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	if req != nil {
		assert.Equal(t, ctx, req.Context(), "Expected context to be set in request")
	}
}

// TestBuilderWithDifferentMethods tests various HTTP methods
func TestBuilderWithDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			builder := F.Pipe2(
				R.Default,
				R.WithURL("https://api.example.com/resource"),
				R.WithMethod(method),
			)

			requester := Requester(builder)
			result := requester(t.Context())()

			assert.True(t, E.IsRight(result), "Expected Right result for method %s", method)

			req := E.GetOrElse(func(error) *http.Request { return nil })(result)
			assert.NotNil(t, req, "Expected non-nil request for method %s", method)
			assert.Equal(t, method, req.Method)
		})
	}
}

// TestBuilderWithJSON tests creating a request with JSON body
func TestBuilderWithJSON(t *testing.T) {
	data := map[string]string{"username": "testuser", "email": "test@example.com"}

	builder := F.Pipe3(
		R.Default,
		R.WithURL("https://api.example.com/v1/users"),
		R.WithMethod("POST"),
		R.WithJSON(data),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https://api.example.com/v1/users", req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.NotNil(t, req.Body)
}

// TestBuilderWithBearer tests adding Bearer token
func TestBuilderWithBearer(t *testing.T) {
	builder := F.Pipe2(
		R.Default,
		R.WithURL("https://api.example.com/protected"),
		R.WithBearer("my-secret-token"),
	)

	requester := Requester(builder)
	result := requester(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	req := E.GetOrElse(func(error) *http.Request { return nil })(result)
	assert.NotNil(t, req, "Expected non-nil request")
	assert.Equal(t, "Bearer my-secret-token", req.Header.Get("Authorization"))
}
