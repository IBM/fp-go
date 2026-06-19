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
	"context"
	"fmt"
	"testing"

	H "net/http"

	R "github.com/IBM/fp-go/v2/context/readerioresult"
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	HT "github.com/IBM/fp-go/v2/http"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

type PostItem struct {
	UserID uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func getTitle(item PostItem) string {
	return item.Title
}

type simpleRequestBuilder struct {
	method  string
	url     string
	headers H.Header
}

func requestBuilder() simpleRequestBuilder {
	return simpleRequestBuilder{method: "GET"}
}

func (b simpleRequestBuilder) WithURL(url string) simpleRequestBuilder {
	b.url = url
	return b
}

func (b simpleRequestBuilder) WithHeader(key, value string) simpleRequestBuilder {
	if b.headers == nil {
		b.headers = make(H.Header)
	} else {
		b.headers = b.headers.Clone()
	}
	b.headers.Set(key, value)
	return b
}

func (b simpleRequestBuilder) Build() R.ReaderIOResult[*H.Request] {
	return func(ctx context.Context) IOE.IOEither[error, *H.Request] {
		return IOE.TryCatchError(func() (*H.Request, error) {
			req, err := H.NewRequestWithContext(ctx, b.method, b.url, nil)
			if err == nil {
				req.Header = b.headers
			}
			return req, err
		})
	}
}

func TestSendSingleRequest(t *testing.T) {

	client := MakeClient(H.DefaultClient)

	req1 := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	readItem := ReadJSON[PostItem](client)

	resp1 := readItem(req1)

	resE := resp1(t.Context())()

	fmt.Println(resE)
}

// setHeaderUnsafe updates a header value in a request object by mutating the request object
func setHeaderUnsafe(key, value string) func(*H.Request) *H.Request {
	return func(req *H.Request) *H.Request {
		req.Header.Set(key, value)
		return req
	}
}

func TestSendSingleRequestWithHeaderUnsafe(t *testing.T) {

	client := MakeClient(H.DefaultClient)

	// this is not safe from a puristic perspective, because the map call mutates the request object
	req1 := F.Pipe2(
		"https://jsonplaceholder.typicode.com/posts/1",
		MakeGetRequest,
		R.Map(setHeaderUnsafe("Content-Type", "text/html")),
	)

	readItem := ReadJSON[PostItem](client)

	resp1 := F.Pipe2(
		req1,
		readItem,
		R.Map(getTitle),
	)

	res := F.Pipe1(
		resp1(t.Context())(),
		E.GetOrElse(errors.ToString),
	)

	assert.Equal(t, "sunt aut facere repellat provident occaecati excepturi optio reprehenderit", res)
}

func TestSendSingleRequestWithHeaderSafe(t *testing.T) {

	client := MakeClient(H.DefaultClient)

	// the request builder assembles config values to construct
	// the final http request. Each `With` step creates a copy of the settings
	// so the flow is pure
	request := requestBuilder().
		WithURL("https://jsonplaceholder.typicode.com/posts/1").
		WithHeader("Content-Type", "text/html").
		Build()

	readItem := ReadJSON[PostItem](client)

	response := F.Pipe2(
		request,
		readItem,
		R.Map(getTitle),
	)

	res := F.Pipe1(
		response(t.Context())(),
		E.GetOrElse(errors.ToString),
	)

	assert.Equal(t, "sunt aut facere repellat provident occaecati excepturi optio reprehenderit", res)
}

// TestReadAll tests the ReadAll function which reads response as bytes
func TestReadAll(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")
	readBytes := ReadAll(client)

	result := readBytes(request)(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	bytes := E.GetOrElse(func(error) []byte { return nil })(result)
	assert.NotNil(t, bytes, "Expected non-nil bytes")
	assert.Greater(t, len(bytes), 0, "Expected non-empty byte array")

	// Verify it contains expected JSON content
	content := string(bytes)
	assert.Contains(t, content, "userId")
	assert.Contains(t, content, "title")
}

// TestReadText tests the ReadText function which reads response as string
func TestReadText(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")
	readText := ReadText(client)

	result := readText(request)(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	text := E.GetOrElse(func(error) string { return "" })(result)
	assert.NotEmpty(t, text, "Expected non-empty text")

	// Verify it contains expected JSON content as text
	assert.Contains(t, text, "userId")
	assert.Contains(t, text, "title")
	assert.Contains(t, text, "sunt aut facere")
}

// TestReadJson tests the deprecated ReadJson function
func TestReadJson(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")
	readItem := ReadJson[PostItem](client)

	result := readItem(request)(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")

	item := E.GetOrElse(func(error) PostItem { return PostItem{} })(result)
	assert.Equal(t, uint(1), item.UserID, "Expected UserID to be 1")
	assert.Equal(t, uint(1), item.Id, "Expected Id to be 1")
	assert.NotEmpty(t, item.Title, "Expected non-empty title")
	assert.NotEmpty(t, item.Body, "Expected non-empty body")
}

// TestReadAllWithInvalidURL tests ReadAll with an invalid URL
func TestReadAllWithInvalidURL(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("http://invalid-domain-that-does-not-exist-12345.com")
	readBytes := ReadAll(client)

	result := readBytes(request)(t.Context())()

	assert.True(t, E.IsLeft(result), "Expected Left result for invalid URL")
}

// TestReadTextWithInvalidURL tests ReadText with an invalid URL
func TestReadTextWithInvalidURL(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("http://invalid-domain-that-does-not-exist-12345.com")
	readText := ReadText(client)

	result := readText(request)(t.Context())()

	assert.True(t, E.IsLeft(result), "Expected Left result for invalid URL")
}

// TestReadJSONWithInvalidURL tests ReadJSON with an invalid URL
func TestReadJSONWithInvalidURL(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("http://invalid-domain-that-does-not-exist-12345.com")
	readItem := ReadJSON[PostItem](client)

	result := readItem(request)(t.Context())()

	assert.True(t, E.IsLeft(result), "Expected Left result for invalid URL")
}

// TestReadJSONWithInvalidJSON tests ReadJSON with non-JSON response
func TestReadJSONWithInvalidJSON(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	// This URL returns HTML, not JSON
	request := MakeGetRequest("https://www.google.com")
	readItem := ReadJSON[PostItem](client)

	result := readItem(request)(t.Context())()

	// Should fail because content-type is not application/json
	assert.True(t, E.IsLeft(result), "Expected Left result for non-JSON response")
}

// TestMakeClientWithCustomClient tests MakeClient with a custom http.Client
func TestMakeClientWithCustomClient(t *testing.T) {
	customClient := H.DefaultClient

	client := MakeClient(customClient)
	assert.NotNil(t, client, "Expected non-nil client")

	// Verify it works
	request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")
	readItem := ReadJSON[PostItem](client)
	result := readItem(request)(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")
}

// TestReadAllComposition tests composing ReadAll with other operations
func TestReadAllComposition(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	// Compose ReadAll with a map operation to get byte length
	readBytes := ReadAll(client)(request)
	readLength := R.Map(func(bytes []byte) int { return len(bytes) })(readBytes)

	result := readLength(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")
	length := E.GetOrElse(func(error) int { return 0 })(result)
	assert.Greater(t, length, 0, "Expected positive byte length")
}

// TestReadTextComposition tests composing ReadText with other operations
func TestReadTextComposition(t *testing.T) {
	client := MakeClient(H.DefaultClient)

	request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	// Compose ReadText with a map operation to get string length
	readText := ReadText(client)(request)
	readLength := R.Map(func(text string) int { return len(text) })(readText)

	result := readLength(t.Context())()

	assert.True(t, E.IsRight(result), "Expected Right result")
	length := E.GetOrElse(func(error) int { return 0 })(result)
	assert.Greater(t, length, 0, "Expected positive string length")
}

// TestReadFullResponse_Success tests successful ReadFullResponse execution
func TestReadFullResponse_Success(t *testing.T) {
	t.Run("returns response and body as pair", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		assert.True(t, E.IsRight(result), "Expected Right result")

		// Extract and verify response
		resp := F.Pipe1(result, E.Map[error](HT.Response))
		assert.True(t, E.IsRight(resp))

		httpResp := E.GetOrElse(func(error) *H.Response { return nil })(resp)
		assert.NotNil(t, httpResp)
		assert.Equal(t, H.StatusOK, httpResp.StatusCode)

		// Extract and verify body
		body := F.Pipe1(result, E.Map[error](HT.Body))
		assert.True(t, E.IsRight(body))

		bodyBytes := E.GetOrElse(func(error) []byte { return nil })(body)
		assert.NotNil(t, bodyBytes)
		assert.Greater(t, len(bodyBytes), 0)
		assert.Contains(t, string(bodyBytes), "userId")
	})
}

// TestReadFullResponse_Failure tests ReadFullResponse error handling
func TestReadFullResponse_Failure(t *testing.T) {
	t.Run("handles invalid URL", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		request := MakeGetRequest("http://invalid-domain-that-does-not-exist-12345.com")

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		assert.True(t, E.IsLeft(result), "Expected Left result for invalid URL")
	})

	t.Run("handles HTTP error status", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		// This URL returns 404
		request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/999999")

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		assert.True(t, E.IsLeft(result), "Expected Left result for 404 status")
	})

	t.Run("preserves error context", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		request := MakeGetRequest("http://invalid-domain-that-does-not-exist-12345.com")

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		err := F.Pipe1(result, E.Fold(
			F.Identity[error],
			func(HT.FullResponse) error { t.Fatal("expected Left"); return nil },
		))

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid-domain-that-does-not-exist-12345.com")
	})
}

// TestReadFullResponse_EdgeCases tests edge cases
func TestReadFullResponse_EdgeCases(t *testing.T) {
	t.Run("handles empty response body", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		// HEAD request typically has no body
		request := MakeRequest("HEAD", "https://jsonplaceholder.typicode.com/posts/1", nil)

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		assert.True(t, E.IsRight(result), "Expected Right result")

		body := F.Pipe1(result, E.Map[error](HT.Body))
		bodyBytes := E.GetOrElse(func(error) []byte { return nil })(body)
		assert.NotNil(t, bodyBytes)
		// HEAD response should have empty or minimal body
		assert.LessOrEqual(t, len(bodyBytes), 0)
	})

	t.Run("handles large response body", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		// Get all posts (larger response)
		request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts")

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		assert.True(t, E.IsRight(result), "Expected Right result")

		body := F.Pipe1(result, E.Map[error](HT.Body))
		bodyBytes := E.GetOrElse(func(error) []byte { return nil })(body)
		assert.NotNil(t, bodyBytes)
		assert.Greater(t, len(bodyBytes), 1000, "Expected large response body")
	})
}

// TestReadFullResponse_Integration tests integration with other functions
func TestReadFullResponse_Integration(t *testing.T) {
	t.Run("composes with Map to extract status code", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

		getStatusCode := F.Flow2(
			ReadFullResponse(client),
			R.Map(F.Flow2(
				HT.Response,
				HT.GetStatusCode,
			)),
		)

		result := getStatusCode(request)(t.Context())()

		assert.True(t, E.IsRight(result))
		statusCode := E.GetOrElse(func(error) int { return 0 })(result)
		assert.Equal(t, H.StatusOK, statusCode)
	})

	t.Run("composes with Chain for conditional processing", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)
		request := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

		// Chain to validate body is not empty
		validateBody := func(fr HT.FullResponse) R.ReaderIOResult[HT.FullResponse] {
			body := HT.Body(fr)
			if len(body) == 0 {
				return R.Left[HT.FullResponse](errors.OnNone("empty body")())
			}
			return R.Of[HT.FullResponse](fr)
		}

		pipeline := F.Flow2(
			ReadFullResponse(client),
			R.Chain(validateBody),
		)

		result := pipeline(request)(t.Context())()

		assert.True(t, E.IsRight(result))
	})

	t.Run("works with custom request builder", func(t *testing.T) {
		client := MakeClient(H.DefaultClient)

		request := requestBuilder().
			WithURL("https://jsonplaceholder.typicode.com/posts/1").
			WithHeader("Accept", "application/json").
			Build()

		fullResp := ReadFullResponse(client)(request)
		result := fullResp(t.Context())()

		assert.True(t, E.IsRight(result))
	})
}
