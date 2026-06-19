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

package http_test

import (
	"context"
	"fmt"
	H "net/http"

	RIOH "github.com/IBM/fp-go/v2/context/readerioresult/http"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	HT "github.com/IBM/fp-go/v2/http"
)

// ExampleReadFullResponse demonstrates basic usage of ReadFullResponse.
func ExampleReadFullResponse() {
	client := RIOH.MakeClient(H.DefaultClient)
	request := RIOH.MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	fullResp := RIOH.ReadFullResponse(client)(request)
	result := fullResp(context.Background())()

	// Extract response and body from the FullResponse pair
	statusCode := F.Pipe1(
		result,
		E.Map[error](F.Flow2(
			HT.Response,
			func(r *H.Response) int { return r.StatusCode },
		)),
	)

	fmt.Println(E.IsRight(statusCode))
	// Output: true
}

// ExampleReadFullResponse_accessingComponents demonstrates accessing response components.
func ExampleReadFullResponse_accessingComponents() {
	client := RIOH.MakeClient(H.DefaultClient)
	request := RIOH.MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	fullResp := RIOH.ReadFullResponse(client)(request)
	result := fullResp(context.Background())()

	// Access the response metadata
	hasContentType := F.Pipe1(
		result,
		E.Map[error](F.Flow2(
			HT.Response,
			func(r *H.Response) bool {
				return r.Header.Get("Content-Type") != ""
			},
		)),
	)

	// Access the body bytes
	hasBody := F.Pipe1(
		result,
		E.Map[error](F.Flow2(
			HT.Body,
			func(b []byte) bool { return len(b) > 0 },
		)),
	)

	fmt.Println(E.IsRight(hasContentType) && E.IsRight(hasBody))
	// Output: true
}
