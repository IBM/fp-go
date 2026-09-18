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
	"io"
	"net/http"
	"net/url"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/http/builder"
	C "github.com/IBM/fp-go/v2/http/content"
	HD "github.com/IBM/fp-go/v2/http/headers"
	IO "github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
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
		ioeither.Map[error](func(r *http.Request) *url.URL {
			return r.URL
		}),
		ioeither.ChainFirstIOK[error](func(u *url.URL) IO.IO[Void] {
			return IO.FromImpure(func() {
				q := u.Query()
				assert.Equal(t, "10", q.Get("limit"))
				assert.Equal(t, "b", q.Get("a"))
			})
		}),
	)

	assert.True(t, E.IsRight(req()))
}

// getRequest executes the requester and returns the request, failing the test on error
func getRequest(t *testing.T, requester IOEither[*http.Request]) *http.Request {
	req, err := E.Unwrap(requester())
	assert.NoError(t, err)
	assert.NotNil(t, req)
	return req
}

// TestBuilderWithoutBody tests creating a request without a body
func TestBuilderWithoutBody(t *testing.T) {
	req := getRequest(t, Requester(F.Pipe2(
		R.Default,
		R.WithURL("https://api.example.com/users"),
		R.WithMethod("GET"),
	)))

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://api.example.com/users", req.URL.String())
	assert.Equal(t, http.NoBody, req.Body)
	assert.Empty(t, req.Header.Get(HD.ContentLength))
}

// TestBuilderWithBody tests creating a request with a body
func TestBuilderWithBody(t *testing.T) {
	bodyData := []byte(`{"name":"John","age":30}`)

	req := getRequest(t, Requester(F.Pipe3(
		R.Default,
		R.WithURL("https://api.example.com/users"),
		R.WithMethod("POST"),
		R.WithBytes(bodyData),
	)))

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "24", req.Header.Get(HD.ContentLength))

	data, err := io.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, bodyData, data)
}

// TestBuilderWithBodyRepeatable tests that executing a requester repeatedly
// produces a fresh, fully readable body each time (e.g. for retries)
func TestBuilderWithBodyRepeatable(t *testing.T) {
	bodyData := []byte(`{"name":"John","age":30}`)

	requester := Requester(F.Pipe3(
		R.Default,
		R.WithURL("https://api.example.com/users"),
		R.WithMethod("POST"),
		R.WithBytes(bodyData),
	))

	for range 2 {
		data, err := io.ReadAll(getRequest(t, requester).Body)
		assert.NoError(t, err)
		assert.Equal(t, bodyData, data)
	}
}

// TestBuilderWithHeaders tests that headers of the builder are set on the request
func TestBuilderWithHeaders(t *testing.T) {
	req := getRequest(t, Requester(F.Pipe3(
		R.Default,
		R.WithURL("https://api.example.com/data"),
		R.WithHeader(HD.Authorization)("Bearer token123"),
		R.WithHeader(HD.Accept)(C.JSON),
	)))

	assert.Equal(t, "Bearer token123", req.Header.Get(HD.Authorization))
	assert.Equal(t, C.JSON, req.Header.Get(HD.Accept))
}

// TestBuilderHeadersAreIsolated tests that modifying the headers of a request
// does not modify the headers of the builder
func TestBuilderHeadersAreIsolated(t *testing.T) {
	builder := F.Pipe2(
		R.Default,
		R.WithURL("https://api.example.com/data"),
		R.WithHeader(HD.Accept)(C.JSON),
	)

	req := getRequest(t, Requester(builder))
	req.Header.Set(HD.XRequestID, "12345")
	req.Header.Add(HD.Accept, C.TextPlain)

	assert.Empty(t, builder.GetHeaders().Get(HD.XRequestID))
	assert.Equal(t, []string{C.JSON}, builder.GetHeaderValues(HD.Accept))
}

// TestBuilderWithInvalidURL tests error handling for invalid URLs
func TestBuilderWithInvalidURL(t *testing.T) {
	requester := Requester(F.Pipe1(
		R.Default,
		R.WithURL("://invalid-url"),
	))

	assert.True(t, E.IsLeft(requester()))
}
