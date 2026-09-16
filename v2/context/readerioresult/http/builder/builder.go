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

// Package builder provides utilities for building HTTP requests in a functional way
// using the ReaderIOResult monad. It integrates with the http/builder package to
// create composable, type-safe HTTP request builders with proper error handling
// and context support.
//
// The main function, Requester, converts a Builder from the http/builder package
// into a ReaderIOResult that produces HTTP requests. This allows for:
//   - Immutable request building with method chaining
//   - Automatic header management including Content-Length
//   - Support for requests with and without bodies
//   - Proper error handling wrapped in Either
//   - Context propagation for cancellation and timeouts
//
// Example usage:
//
//	import (
//	    "context"
//	    B "github.com/IBM/fp-go/v2/http/builder"
//	    RB "github.com/IBM/fp-go/v2/context/readerioresult/http/builder"
//	)
//
//	builder := F.Pipe3(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithMethod("POST"),
//	    B.WithJSONBody(userData),
//	)
//
//	requester := RB.Requester(builder)
//	result := requester(t.Context())()
package builder

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	RIOEH "github.com/IBM/fp-go/v2/context/readerioresult/http"
	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	HTTP "github.com/IBM/fp-go/v2/http"
	R "github.com/IBM/fp-go/v2/http/builder"
	H "github.com/IBM/fp-go/v2/http/headers"
	LZ "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

// Requester converts an http/builder.Builder into a ReaderIOResult that produces HTTP requests.
// It handles both requests with and without bodies, automatically managing headers including
// Content-Length for requests with bodies.
//
// The function performs the following operations:
//  1. Extracts the request body (if present) from the builder
//  2. Creates appropriate request constructor (with or without body)
//  3. Applies the target URL from the builder
//  4. Applies the HTTP method from the builder
//  5. Merges headers from the builder into the request
//  6. Handles any errors that occur during request construction
//
// For requests with a body:
//   - Sets the Content-Length header automatically
//   - Uses bytes.NewReader to create the request body
//   - Merges builder headers into the request
//
// For requests without a body:
//   - Creates a request with nil body
//   - Merges builder headers into the request
//
// Parameters:
//   - builder: A pointer to an http/builder.Builder containing request configuration
//
// Returns:
//   - A Requester (ReaderIOResult[*http.Request]) that, when executed with a context,
//     produces either an error or a configured *http.Request
//
// Example with body:
//
//	import (
//	    B "github.com/IBM/fp-go/v2/http/builder"
//	    RB "github.com/IBM/fp-go/v2/context/readerioresult/http/builder"
//	)
//
//	builder := F.Pipe3(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithMethod("POST"),
//	    B.WithJSONBody(map[string]string{"name": "John"}),
//	)
//	requester := RB.Requester(builder)
//	result := requester(t.Context())()
//
// Example without body:
//
//	builder := F.Pipe2(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithMethod("GET"),
//	)
//	requester := RB.Requester(builder)
//	result := requester(t.Context())()
func Requester(builder *R.Builder) RIOEH.Requester {
	return F.Pipe4(
		builder.GetBody(),
		O.Fold(LZ.Of(noBody), result.Map(withBody)),
		result.Ap[RIOE.ReaderIOResult[*http.Request]](F.Pipe1(
			builder.GetTargetURL(),
			result.Map(F.Curry2(F.Bind1of3(RIOEH.MakeRequest)(builder.GetMethod()))),
		)),
		result.GetOrElse(RIOE.Left[*http.Request]),
		RIOE.Map(F.Bind1of2(mergeHeaders)(builder.GetHeaders())),
	)
}

var (
	// toReader converts a byte slice into a fresh [io.Reader]
	toReader = F.Flow2(
		bytes.NewReader,
		FL.ToReader[*bytes.Reader],
	)

	// noBody feeds an empty body into the request constructor
	noBody = result.Of(withReader(RIOE.Of[io.Reader](nil)))
)

// withReader feeds the body into a request constructor that is still waiting for its body.
func withReader(body RIOE.ReaderIOResult[io.Reader]) func(RIOE.Kleisli[io.Reader, *http.Request]) RIOE.ReaderIOResult[*http.Request] {
	return F.Bind1st(RIOE.MonadChain[io.Reader, *http.Request], body)
}

// withBody feeds data as the body into a request constructor and sets the Content-Length header.
// The [io.Reader] over data is created on each execution, so the resulting requester can be
// executed repeatedly (e.g. when retrying).
func withBody(data []byte) func(RIOE.Kleisli[io.Reader, *http.Request]) RIOE.ReaderIOResult[*http.Request] {
	return F.Flow2(
		withReader(F.Pipe1(
			RIOE.Of(data),
			RIOE.Map(toReader),
		)),
		RIOE.Map(HTTP.WithHeader(H.ContentLength)(strconv.Itoa(len(data)))),
	)
}

// mergeHeaders merges a copy of headers into the headers of the request and returns the request.
// Copying ensures that the request never shares header maps or value slices with the builder.
func mergeHeaders(headers http.Header, req *http.Request) *http.Request {
	req.Header = H.Monoid.Concat(req.Header, headers.Clone())
	return req
}
