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

// Package http provides functional HTTP client utilities built on top of ReaderIOResult monad.
// It offers a composable way to make HTTP requests with context support, error handling,
// and response parsing capabilities. The package follows functional programming principles
// to ensure type-safe, testable, and maintainable HTTP operations.
//
// The main abstractions include:
//   - Requester: A reader that constructs HTTP requests with context
//   - Client: An interface for executing HTTP requests
//   - Response readers: Functions to parse responses as bytes, text, or JSON
//
// Example usage:
//
//	client := MakeClient(http.DefaultClient)
//	request := MakeGetRequest("https://api.example.com/data")
//	result := ReadJSON[MyType](client)(request)
//	response := result(t.Context())()
package http

import (
	"io"
	"net/http"

	B "github.com/IBM/fp-go/v2/bytes"
	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	H "github.com/IBM/fp-go/v2/http"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	J "github.com/IBM/fp-go/v2/json"
	P "github.com/IBM/fp-go/v2/pair"
)

type (
	// Requester is a reader that constructs an HTTP request with context support.
	// It represents a computation that, given a context, produces either an error
	// or an HTTP request. This allows for composable request building with proper
	// error handling and context propagation.
	Requester = RIOE.ReaderIOResult[*http.Request]

	// Client is an interface for executing HTTP requests in a functional way.
	// It wraps the standard http.Client and provides a Do method that works
	// with the ReaderIOResult monad for composable, type-safe HTTP operations.
	Client interface {
		// Do executes an HTTP request and returns the response wrapped in a ReaderIOResult.
		// It takes a Requester (which builds the request) and returns a computation that,
		// when executed with a context, performs the HTTP request and returns either
		// an error or the HTTP response.
		//
		// Parameters:
		//   - req: A Requester that builds the HTTP request
		//
		// Returns:
		//   - A ReaderIOResult that produces either an error or an *http.Response
		Do(Requester) RIOE.ReaderIOResult[*http.Response]
	}

	// client is the internal implementation of the Client interface.
	// It wraps a standard http.Client and provides functional HTTP operations.
	client struct {
		delegate *http.Client
		doIOE    IOE.Kleisli[error, *http.Request, *http.Response]
	}
)

var (
	// MakeRequest is an eitherized version of http.NewRequestWithContext.
	// It creates a Requester that builds an HTTP request with the given method, URL, and body.
	// This function properly handles errors and wraps them in the Either monad.
	//
	// Parameters:
	//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
	//   - url: The target URL for the request
	//   - body: Optional request body (can be nil)
	//
	// Returns:
	//   - A Requester that produces either an error or an *http.Request
	MakeRequest = RIOE.Eitherize3(http.NewRequestWithContext)

	// makeRequest is a partially applied version of MakeRequest with the context parameter bound.
	makeRequest = F.Bind13of3(MakeRequest)

	// MakeGetRequest creates a GET request for the specified URL.
	// It's a convenience function that specializes MakeRequest for GET requests with no body.
	//
	// Parameters:
	//   - url: The target URL for the GET request
	//
	// Returns:
	//   - A Requester that produces either an error or an *http.Request
	//
	// Example:
	//   req := MakeGetRequest("https://api.example.com/users")
	MakeGetRequest = makeRequest("GET", nil)
)

func (client client) Do(req Requester) RIOE.ReaderIOResult[*http.Response] {
	return F.Pipe1(
		req,
		RIOE.ChainIOEitherK(client.doIOE),
	)
}

// MakeClient creates a functional HTTP client wrapper around a standard http.Client.
// The returned Client provides methods for executing HTTP requests in a functional,
// composable way using the ReaderIOResult monad.
//
// Parameters:
//   - httpClient: A standard *http.Client to wrap (e.g., http.DefaultClient)
//
// Returns:
//   - A Client that can execute HTTP requests functionally
//
// Example:
//
//	client := MakeClient(http.DefaultClient)
//	// or with custom client
//	customClient := &http.Client{Timeout: 10 * time.Second}
//	client := MakeClient(customClient)
func MakeClient(httpClient *http.Client) Client {
	return client{delegate: httpClient, doIOE: IOE.Eitherize1(httpClient.Do)}
}

// ReadFullResponse sends an HTTP request, reads the complete response body as a byte array,
// and returns both the response and body as a tuple (FullResponse).
// It validates the HTTP status code and handles errors appropriately.
//
// The function performs the following steps:
//  1. Executes the HTTP request using the provided client
//  2. Validates the response status code (checks for HTTP errors)
//  3. Reads the entire response body into a byte array
//  4. Returns a tuple containing the response and body
//
// Parameters:
//   - client: The HTTP client to use for executing the request
//
// Returns:
//   - A function that takes a Requester and returns a ReaderIOResult[FullResponse]
//     where FullResponse is a tuple of (*http.Response, []byte)
//
// Example:
//
//	client := MakeClient(http.DefaultClient)
//	request := MakeGetRequest("https://api.example.com/data")
//	fullResp := ReadFullResponse(client)(request)
//	result := fullResp(t.Context())()
func ReadFullResponse(client Client) RIOE.Operator[*http.Request, H.FullResponse] {
	return func(req Requester) RIOE.ReaderIOResult[H.FullResponse] {
		return F.Flow3(
			client.Do(req),
			IOE.ChainEitherK(H.ValidateResponse),
			IOE.Chain(func(resp *http.Response) IOE.IOEither[error, H.FullResponse] {
				return F.Pipe1(
					F.Pipe3(
						resp,
						H.GetBody,
						IOE.Of[error, io.ReadCloser],
						IOEF.ReadAll[io.ReadCloser],
					),
					IOE.Map[error](F.Bind1st(P.MakePair[*http.Response, []byte], resp)),
				)
			}),
		)
	}
}

// ReadAll sends an HTTP request and reads the complete response body as a byte array.
// It validates the HTTP status code and returns the raw response body bytes.
// This is useful when you need to process the response body in a custom way.
//
// Parameters:
//   - client: The HTTP client to use for executing the request
//
// Returns:
//   - A function that takes a Requester and returns a ReaderIOResult[[]byte]
//     containing the response body as bytes
//
// Example:
//
//	client := MakeClient(http.DefaultClient)
//	request := MakeGetRequest("https://api.example.com/data")
//	readBytes := ReadAll(client)
//	result := readBytes(request)(t.Context())()
func ReadAll(client Client) RIOE.Operator[*http.Request, []byte] {
	return F.Flow2(
		ReadFullResponse(client),
		RIOE.Map(H.Body),
	)
}

// ReadText sends an HTTP request, reads the response body, and converts it to a string.
// It validates the HTTP status code and returns the response body as a UTF-8 string.
// This is convenient for APIs that return plain text responses.
//
// Parameters:
//   - client: The HTTP client to use for executing the request
//
// Returns:
//   - A function that takes a Requester and returns a ReaderIOResult[string]
//     containing the response body as a string
//
// Example:
//
//	client := MakeClient(http.DefaultClient)
//	request := MakeGetRequest("https://api.example.com/text")
//	readText := ReadText(client)
//	result := readText(request)(t.Context())()
func ReadText(client Client) RIOE.Operator[*http.Request, string] {
	return F.Flow2(
		ReadAll(client),
		RIOE.Map(B.ToString),
	)
}

// ReadJson sends an HTTP request, reads the response, and parses it as JSON.
//
// Deprecated: Use [ReadJSON] instead. This function is kept for backward compatibility
// but will be removed in a future version. The capitalized version follows Go naming
// conventions for acronyms.
func ReadJson[A any](client Client) RIOE.Operator[*http.Request, A] {
	return ReadJSON[A](client)
}

// readJSON is an internal helper that reads the response body and validates JSON content type.
// It performs the following validations:
//  1. Validates HTTP status code
//  2. Validates that the response Content-Type is application/json
//  3. Reads the response body as bytes
//
// This function is used internally by ReadJSON to ensure proper JSON response handling.
func readJSON(client Client) RIOE.Operator[*http.Request, []byte] {
	return F.Flow3(
		ReadFullResponse(client),
		RIOE.ChainFirstEitherK(F.Flow2(
			H.Response,
			H.ValidateJSONResponse,
		)),
		RIOE.Map(H.Body),
	)
}

// ReadJSON sends an HTTP request, reads the response, and parses it as JSON into type A.
// It validates both the HTTP status code and the Content-Type header to ensure the
// response is valid JSON before attempting to unmarshal.
//
// Type Parameters:
//   - A: The target type to unmarshal the JSON response into
//
// Parameters:
//   - client: The HTTP client to use for executing the request
//
// Returns:
//   - A function that takes a Requester and returns a ReaderIOResult[A]
//     containing the parsed JSON data
//
// Example:
//
//	type User struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//
//	client := MakeClient(http.DefaultClient)
//	request := MakeGetRequest("https://api.example.com/user/1")
//	readUser := ReadJSON[User](client)
//	result := readUser(request)(t.Context())()
func ReadJSON[A any](client Client) RIOE.Operator[*http.Request, A] {
	return F.Flow2(
		readJSON(client),
		RIOE.ChainEitherK(J.Unmarshal[A]),
	)
}
