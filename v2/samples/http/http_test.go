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

	HTTP "net/http"

	A "github.com/IBM/fp-go/v2/array"
	R "github.com/IBM/fp-go/v2/context/readerioresult"
	H "github.com/IBM/fp-go/v2/context/readerioresult/http"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/result"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

type PostItem struct {
	UserID uint   `json:"userId"`
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type CatFact struct {
	Fact string `json:"fact"`
}

func idxToURL(idx int) string {
	return fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", idx+1)
}

// TestMultipleHttpRequests shows how to execute multiple HTTP requests in parallel assuming
// that the response structure of all requests is identical, which is why we can use [R.TraverseArray]
func TestMultipleHttpRequests(t *testing.T) {
	// prepare the http client
	client := H.MakeClient(HTTP.DefaultClient)
	// readSinglePost sends a GET request and parses the response as [PostItem]
	readSinglePost := H.ReadJSON[PostItem](client)

	// total number of http requests
	count := 10

	data := F.Pipe3(
		A.MakeBy(count, idxToURL),
		R.TraverseArray(F.Flow3(
			H.MakeGetRequest,
			readSinglePost,
			R.ChainFirstIOK(IO.Logf[PostItem]("Log Single: %v")),
		)),
		R.ChainFirstIOK(IO.Logf[[]PostItem]("Log Result: %v")),
		R.Map(A.Size[PostItem]),
	)

	res := data(context.Background())

	assert.Equal(t, result.Of(count), res())
}

func heterogeneousHTTPRequests() ReaderIOResult[T.Tuple2[PostItem, CatFact]] {
	// prepare the http client
	client := H.MakeClient(HTTP.DefaultClient)
	// readSinglePost sends a GET request and parses the response as [PostItem]
	readSinglePost := H.ReadJSON[PostItem](client)
	// readSingleCatFact sends a GET request and parses the response as [CatFact]
	readSingleCatFact := H.ReadJSON[CatFact](client)

	return F.Pipe3(
		T.MakeTuple2("https://jsonplaceholder.typicode.com/posts/1", "https://catfact.ninja/fact"),
		T.Map2(H.MakeGetRequest, H.MakeGetRequest),
		R.TraverseTuple2(
			readSinglePost,
			readSingleCatFact,
		),
		R.ChainFirstIOK(IO.Logf[T.Tuple2[PostItem, CatFact]]("Log Result: %v")),
	)

}

// TestHeterogeneousHttpRequests shows how to execute multiple HTTP requests in parallel when
// the response structure of these requests is different. We use [R.TraverseTuple2] to account for the different types
func TestHeterogeneousHttpRequests(t *testing.T) {
	data := heterogeneousHTTPRequests()

	result := data(context.Background())

	fmt.Println(result())
}

// BenchmarkHeterogeneousHttpRequests shows how to execute multiple HTTP requests in parallel when
// the response structure of these requests is different. We use [R.TraverseTuple2] to account for the different types
func BenchmarkHeterogeneousHttpRequests(b *testing.B) {
	for b.Loop() {
		heterogeneousHTTPRequests()(b.Context())()
	}
}
