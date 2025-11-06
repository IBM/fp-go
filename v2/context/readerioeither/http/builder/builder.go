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
// using the ReaderIOEither monad. It integrates with the http/builder package to
// create composable, type-safe HTTP request builders with proper error handling
// and context support.
//
// The main function, Requester, converts a Builder from the http/builder package
// into a ReaderIOEither that produces HTTP requests. This allows for:
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
//	    RB "github.com/IBM/fp-go/v2/context/readerioeither/http/builder"
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
//	result := requester(context.Background())()
package builder

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	RIOE "github.com/IBM/fp-go/v2/context/readerioeither"
	RIOEH "github.com/IBM/fp-go/v2/context/readerioeither/http"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/http/builder"
	H "github.com/IBM/fp-go/v2/http/headers"
	LZ "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

// Requester converts an http/builder.Builder into a ReaderIOEither that produces HTTP requests.
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
//   - A Requester (ReaderIOEither[*http.Request]) that, when executed with a context,
//     produces either an error or a configured *http.Request
//
// Example with body:
//
//	import (
//	    B "github.com/IBM/fp-go/v2/http/builder"
//	    RB "github.com/IBM/fp-go/v2/context/readerioeither/http/builder"
//	)
//
//	builder := F.Pipe3(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithMethod("POST"),
//	    B.WithJSONBody(map[string]string{"name": "John"}),
//	)
//	requester := RB.Requester(builder)
//	result := requester(context.Background())()
//
// Example without body:
//
//	builder := F.Pipe2(
//	    B.Default,
//	    B.WithURL("https://api.example.com/users"),
//	    B.WithMethod("GET"),
//	)
//	requester := RB.Requester(builder)
//	result := requester(context.Background())()
func Requester(builder *R.Builder) RIOEH.Requester {

	withBody := F.Curry3(func(data []byte, url string, method string) RIOE.ReaderIOEither[*http.Request] {
		return RIOE.TryCatch(func(ctx context.Context) func() (*http.Request, error) {
			return func() (*http.Request, error) {
				req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(data))
				if err == nil {
					req.Header.Set(H.ContentLength, strconv.Itoa(len(data)))
					H.Monoid.Concat(req.Header, builder.GetHeaders())
				}
				return req, err
			}
		})
	})

	withoutBody := F.Curry2(func(url string, method string) RIOE.ReaderIOEither[*http.Request] {
		return RIOE.TryCatch(func(ctx context.Context) func() (*http.Request, error) {
			return func() (*http.Request, error) {
				req, err := http.NewRequestWithContext(ctx, method, url, nil)
				if err == nil {
					H.Monoid.Concat(req.Header, builder.GetHeaders())
				}
				return req, err
			}
		})
	})

	return F.Pipe5(
		builder.GetBody(),
		O.Fold(LZ.Of(E.Of[error](withoutBody)), E.Map[error](withBody)),
		E.Ap[func(string) RIOE.ReaderIOEither[*http.Request]](builder.GetTargetURL()),
		E.Flap[error, RIOE.ReaderIOEither[*http.Request]](builder.GetMethod()),
		E.GetOrElse(RIOE.Left[*http.Request]),
		RIOE.Map(func(req *http.Request) *http.Request {
			req.Header = H.Monoid.Concat(req.Header, builder.GetHeaders())
			return req
		}),
	)
}
