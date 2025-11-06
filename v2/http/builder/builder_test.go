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
	"fmt"
	"net/http"
	"net/url"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
	FD "github.com/IBM/fp-go/v2/http/form"
	H "github.com/IBM/fp-go/v2/http/headers"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {

	name := H.ContentType
	withContentType := WithHeader(name)
	withoutContentType := WithoutHeader(name)

	b1 := F.Pipe1(
		Default,
		withContentType(C.JSON),
	)

	b2 := F.Pipe1(
		b1,
		withContentType(C.TextPlain),
	)

	b3 := F.Pipe1(
		b2,
		withoutContentType,
	)

	assert.Equal(t, O.None[string](), Default.GetHeader(name))
	assert.Equal(t, O.Of(C.JSON), b1.GetHeader(name))
	assert.Equal(t, O.Of(C.TextPlain), b2.GetHeader(name))
	assert.Equal(t, O.None[string](), b3.GetHeader(name))
}

func TestWithFormData(t *testing.T) {
	data := F.Pipe1(
		FD.Default,
		FD.WithValue("a")("b"),
	)

	res := F.Pipe1(
		Default,
		WithFormData(data),
	)

	assert.Equal(t, C.FormEncoded, Headers.Get(res).Get(H.ContentType))
}

func TestHash(t *testing.T) {

	b1 := F.Pipe4(
		Default,
		WithContentType(C.JSON),
		WithHeader(H.Accept)(C.JSON),
		WithURL("http://www.example.com"),
		WithJSON(map[string]string{"a": "b"}),
	)

	b2 := F.Pipe4(
		Default,
		WithURL("http://www.example.com"),
		WithHeader(H.Accept)(C.JSON),
		WithContentType(C.JSON),
		WithJSON(map[string]string{"a": "b"}),
	)

	assert.Equal(t, MakeHash(b1), MakeHash(b2))
	assert.NotEqual(t, MakeHash(Default), MakeHash(b2))

	fmt.Println(MakeHash(b1))
}

// TestGetTargetURL tests URL construction with query parameters
func TestGetTargetURL(t *testing.T) {
	builder := F.Pipe3(
		Default,
		WithURL("http://www.example.com?existing=param"),
		WithQueryArg("limit")("10"),
		WithQueryArg("offset")("20"),
	)

	result := builder.GetTargetURL()
	assert.True(t, E.IsRight(result), "Expected Right result")

	url := E.GetOrElse(func(error) string { return "" })(result)
	assert.Contains(t, url, "limit=10")
	assert.Contains(t, url, "offset=20")
	assert.Contains(t, url, "existing=param")
}

// TestGetTargetURLWithInvalidURL tests error handling for invalid URLs
func TestGetTargetURLWithInvalidURL(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithURL("://invalid-url"),
	)

	result := builder.GetTargetURL()
	assert.True(t, E.IsLeft(result), "Expected Left result for invalid URL")
}

// TestGetTargetUrl tests the deprecated GetTargetUrl function
func TestGetTargetUrl(t *testing.T) {
	builder := F.Pipe2(
		Default,
		WithURL("http://www.example.com"),
		WithQueryArg("test")("value"),
	)

	result := builder.GetTargetUrl()
	assert.True(t, E.IsRight(result), "Expected Right result")

	url := E.GetOrElse(func(error) string { return "" })(result)
	assert.Contains(t, url, "test=value")
}

// TestSetMethod tests the SetMethod function
func TestSetMethod(t *testing.T) {
	builder := Default.SetMethod("POST")

	assert.Equal(t, "POST", builder.GetMethod())
}

// TestSetQuery tests the SetQuery function
func TestSetQuery(t *testing.T) {
	query := make(url.Values)
	query.Set("key1", "value1")
	query.Set("key2", "value2")

	builder := Default.SetQuery(query)

	assert.Equal(t, "value1", builder.GetQuery().Get("key1"))
	assert.Equal(t, "value2", builder.GetQuery().Get("key2"))
}

// TestSetHeaders tests the SetHeaders function
func TestSetHeaders(t *testing.T) {
	headers := make(http.Header)
	headers.Set("X-Custom-Header", "custom-value")
	headers.Set("Authorization", "Bearer token")

	builder := Default.SetHeaders(headers)

	assert.Equal(t, "custom-value", builder.GetHeaders().Get("X-Custom-Header"))
	assert.Equal(t, "Bearer token", builder.GetHeaders().Get("Authorization"))
}

// TestGetHeaderValues tests the GetHeaderValues function
func TestGetHeaderValues(t *testing.T) {
	builder := F.Pipe2(
		Default,
		WithHeader("Accept")("application/json"),
		WithHeader("Accept")("text/html"),
	)

	values := builder.GetHeaderValues("Accept")
	assert.Contains(t, values, "text/html")
}

// TestGetUrl tests the deprecated GetUrl function
func TestGetUrl(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithURL("http://www.example.com"),
	)

	assert.Equal(t, "http://www.example.com", builder.GetUrl())
}

// TestSetUrl tests the deprecated SetUrl function
func TestSetUrl(t *testing.T) {
	builder := Default.SetUrl("http://www.example.com")

	assert.Equal(t, "http://www.example.com", builder.GetURL())
}

// TestWithJson tests the deprecated WithJson function
func TestWithJson(t *testing.T) {
	data := map[string]string{"key": "value"}

	builder := F.Pipe1(
		Default,
		WithJson(data),
	)

	contentType := O.GetOrElse(F.Constant(""))(builder.GetHeader(H.ContentType))
	assert.Equal(t, C.JSON, contentType)
	assert.True(t, O.IsSome(builder.GetBody()))
}

// TestQueryArg tests the QueryArg lens
func TestQueryArg(t *testing.T) {
	lens := QueryArg("test")

	builder := F.Pipe1(
		Default,
		lens.Set(O.Some("value")),
	)

	assert.Equal(t, O.Some("value"), lens.Get(builder))
	assert.Equal(t, "value", builder.GetQuery().Get("test"))
}

// TestWithQueryArg tests the WithQueryArg function
func TestWithQueryArg(t *testing.T) {
	builder := F.Pipe2(
		Default,
		WithQueryArg("param1")("value1"),
		WithQueryArg("param2")("value2"),
	)

	assert.Equal(t, "value1", builder.GetQuery().Get("param1"))
	assert.Equal(t, "value2", builder.GetQuery().Get("param2"))
}

// TestWithoutQueryArg tests the WithoutQueryArg function
func TestWithoutQueryArg(t *testing.T) {
	builder := F.Pipe3(
		Default,
		WithQueryArg("param1")("value1"),
		WithQueryArg("param2")("value2"),
		WithoutQueryArg("param1"),
	)

	assert.Equal(t, "", builder.GetQuery().Get("param1"))
	assert.Equal(t, "value2", builder.GetQuery().Get("param2"))
}

// TestGetHash tests the GetHash method
func TestGetHash(t *testing.T) {
	builder := F.Pipe2(
		Default,
		WithURL("http://www.example.com"),
		WithMethod("POST"),
	)

	hash := builder.GetHash()
	assert.NotEmpty(t, hash)
	assert.Equal(t, MakeHash(builder), hash)
}

// TestWithBytes tests the WithBytes function
func TestWithBytes(t *testing.T) {
	data := []byte("test data")

	builder := F.Pipe1(
		Default,
		WithBytes(data),
	)

	body := builder.GetBody()
	assert.True(t, O.IsSome(body))
}

// TestWithoutBody tests the WithoutBody function
func TestWithoutBody(t *testing.T) {
	builder := F.Pipe2(
		Default,
		WithBytes([]byte("data")),
		WithoutBody,
	)

	assert.True(t, O.IsNone(builder.GetBody()))
}

// TestWithGet tests the WithGet convenience function
func TestWithGet(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithGet,
	)

	assert.Equal(t, "GET", builder.GetMethod())
}

// TestWithPost tests the WithPost convenience function
func TestWithPost(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithPost,
	)

	assert.Equal(t, "POST", builder.GetMethod())
}

// TestWithPut tests the WithPut convenience function
func TestWithPut(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithPut,
	)

	assert.Equal(t, "PUT", builder.GetMethod())
}

// TestWithDelete tests the WithDelete convenience function
func TestWithDelete(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithDelete,
	)

	assert.Equal(t, "DELETE", builder.GetMethod())
}

// TestWithBearer tests the WithBearer function
func TestWithBearer(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithBearer("my-token"),
	)

	auth := O.GetOrElse(F.Constant(""))(builder.GetHeader(H.Authorization))
	assert.Equal(t, "Bearer my-token", auth)
}

// TestWithContentType tests the WithContentType function
func TestWithContentType(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithContentType(C.TextPlain),
	)

	contentType := O.GetOrElse(F.Constant(""))(builder.GetHeader(H.ContentType))
	assert.Equal(t, C.TextPlain, contentType)
}

// TestWithAuthorization tests the WithAuthorization function
func TestWithAuthorization(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithAuthorization("Basic abc123"),
	)

	auth := O.GetOrElse(F.Constant(""))(builder.GetHeader(H.Authorization))
	assert.Equal(t, "Basic abc123", auth)
}

// TestBuilderChaining tests that builder operations can be chained
func TestBuilderChaining(t *testing.T) {
	builder := F.Pipe3(
		Default,
		WithURL("http://www.example.com"),
		WithMethod("POST"),
		WithHeader("X-Test")("test-value"),
	)

	// Verify all operations were applied
	assert.Equal(t, "http://www.example.com", builder.GetURL())
	assert.Equal(t, "POST", builder.GetMethod())

	testHeader := O.GetOrElse(F.Constant(""))(builder.GetHeader("X-Test"))
	assert.Equal(t, "test-value", testHeader)
}

// TestWithQuery tests the WithQuery function
func TestWithQuery(t *testing.T) {
	query := make(url.Values)
	query.Set("key1", "value1")
	query.Set("key2", "value2")

	builder := F.Pipe1(
		Default,
		WithQuery(query),
	)

	assert.Equal(t, "value1", builder.GetQuery().Get("key1"))
	assert.Equal(t, "value2", builder.GetQuery().Get("key2"))
}

// TestWithHeaders tests the WithHeaders function
func TestWithHeaders(t *testing.T) {
	headers := make(http.Header)
	headers.Set("X-Test", "test-value")

	builder := F.Pipe1(
		Default,
		WithHeaders(headers),
	)

	assert.Equal(t, "test-value", builder.GetHeaders().Get("X-Test"))
}

// TestWithUrl tests the deprecated WithUrl function
func TestWithUrl(t *testing.T) {
	builder := F.Pipe1(
		Default,
		WithUrl("http://www.example.com"),
	)

	assert.Equal(t, "http://www.example.com", builder.GetURL())
}

// TestComplexBuilderComposition tests a complex builder composition
func TestComplexBuilderComposition(t *testing.T) {
	builder := F.Pipe5(
		Default,
		WithURL("http://api.example.com/users"),
		WithPost,
		WithJSON(map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		}),
		WithBearer("secret-token"),
		WithQueryArg("notify")("true"),
	)

	assert.Equal(t, "http://api.example.com/users", builder.GetURL())
	assert.Equal(t, "POST", builder.GetMethod())

	contentType := O.GetOrElse(F.Constant(""))(builder.GetHeader(H.ContentType))
	assert.Equal(t, C.JSON, contentType)

	auth := O.GetOrElse(F.Constant(""))(builder.GetHeader(H.Authorization))
	assert.Equal(t, "Bearer secret-token", auth)

	assert.Equal(t, "true", builder.GetQuery().Get("notify"))
	assert.True(t, O.IsSome(builder.GetBody()))
}
