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

package benchmarks

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	HTTP "net/http"

	A "github.com/IBM/fp-go/v2/array"
	R "github.com/IBM/fp-go/v2/context/readerioresult"
	H "github.com/IBM/fp-go/v2/context/readerioresult/http"
	F "github.com/IBM/fp-go/v2/function"
	T "github.com/IBM/fp-go/v2/tuple"
)

type PostItem struct {
	UserID uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type CatFact struct {
	Fact string `json:"fact"`
}

func heterogeneousHTTPRequests(count int) R.ReaderIOResult[[]T.Tuple2[PostItem, CatFact]] {
	// prepare the http client
	client := H.MakeClient(HTTP.DefaultClient)
	// readSinglePost sends a GET request and parses the response as [PostItem]
	readSinglePost := H.ReadJSON[PostItem](client)
	// readSingleCatFact sends a GET request and parses the response as [CatFact]
	readSingleCatFact := H.ReadJSON[CatFact](client)

	single := F.Pipe2(
		T.MakeTuple2("https://jsonplaceholder.typicode.com/posts/1", "https://catfact.ninja/fact"),
		T.Map2(H.MakeGetRequest, H.MakeGetRequest),
		R.TraverseTuple2(
			readSinglePost,
			readSingleCatFact,
		),
	)

	return F.Pipe1(
		A.Replicate(count, single),
		R.SequenceArray[T.Tuple2[PostItem, CatFact]],
	)
}

func heterogeneousHTTPRequestsIdiomatic(count int) ([]T.Tuple2[PostItem, CatFact], error) {
	// prepare the http client
	var result []T.Tuple2[PostItem, CatFact]

	for i := 0; i < count; i++ {
		resp, err := HTTP.Get("https://jsonplaceholder.typicode.com/posts/1")
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var item PostItem
		err = json.Unmarshal(body, &item)
		if err != nil {
			return nil, err
		}
		resp, err = HTTP.Get("https://catfact.ninja/fact")
		if err != nil {
			return nil, err
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var fact CatFact
		err = json.Unmarshal(body, &item)
		if err != nil {
			return nil, err
		}
		result = append(result, T.MakeTuple2(item, fact))
	}
	return result, nil
}

// BenchmarkHeterogeneousHttpRequests shows how to execute multiple HTTP requests in parallel when
// the response structure of these requests is different. We use [R.TraverseTuple2] to account for the different types
func BenchmarkHeterogeneousHttpRequests(b *testing.B) {

	count := 100
	var benchResults any

	b.Run("functional", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResults = heterogeneousHTTPRequests(count)(context.Background())()
		}
	})

	b.Run("idiomatic", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			benchResults, _ = heterogeneousHTTPRequestsIdiomatic(count)
		}
	})

	globalResult = benchResults
}
