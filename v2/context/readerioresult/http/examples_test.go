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

	A "github.com/IBM/fp-go/v2/array"
	RIOH "github.com/IBM/fp-go/v2/context/readerioresult/http"
	F "github.com/IBM/fp-go/v2/function"
	HT "github.com/IBM/fp-go/v2/http"
	"github.com/IBM/fp-go/v2/result"
)

// ExampleReadFullResponse demonstrates basic usage of ReadFullResponse.
func ExampleReadFullResponse() {
	client := RIOH.MakeClient(H.DefaultClient)
	request := RIOH.MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	fullResp := RIOH.ReadFullResponse(client)(request)
	res := fullResp(context.Background())()

	// Extract response and body from the FullResponse pair
	statusCode := F.Pipe1(
		res,
		result.Map(F.Flow2(
			HT.Response,
			HT.GetStatusCode,
		)),
	)

	fmt.Println(result.IsRight(statusCode))
	// Output: true
}

// ExampleReadFullResponse_accessingComponents demonstrates accessing response components.
func ExampleReadFullResponse_accessingComponents() {
	client := RIOH.MakeClient(H.DefaultClient)
	request := RIOH.MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	fullResp := RIOH.ReadFullResponse(client)(request)
	res := fullResp(context.Background())()

	// Access the response metadata
	hasContentType := F.Pipe1(
		res,
		result.Map(F.Flow2(
			HT.Response,
			func(r *H.Response) bool {
				return r.Header.Get("Content-Type") != ""
			},
		)),
	)

	// Access the body bytes
	hasBody := F.Pipe1(
		res,
		result.Map(F.Flow2(
			HT.Body,
			A.IsNonEmpty,
		)),
	)

	fmt.Println(result.IsRight(hasContentType) && result.IsRight(hasBody))
	// Output: true
}
