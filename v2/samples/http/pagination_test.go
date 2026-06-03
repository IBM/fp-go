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
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"slices"
	"strconv"

	HTTP "net/http"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	H "github.com/IBM/fp-go/v2/context/readerioresult/http"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iterresult"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

// pageResponse is the JSON shape returned by the paginated REST API.
// When Next is empty the caller has reached the final page.
type pageResponse struct {
	Items []string `json:"items"`
	Next  string   `json:"next,omitempty"`
}

// getItems extracts the Items slice from a pageResponse.
func getItems(r pageResponse) []string {
	return r.Items
}

// getNext converts the Next field of a pageResponse to an Option:
// None for an empty string (last page), Some(url) otherwise.
func getNext(r pageResponse) Option[string] {
	return F.Pipe1(
		r.Next,
		O.FromNonZero[string](),
	)
}

// newPaginatedServer starts an in-process test server that serves GET /items
// with simple offset-based pagination over the supplied data slice.  Each page
// carries at most pageSize items; the last page omits the "next" field.
func newPaginatedServer(data []string, pageSize int) *httptest.Server {
	var srv *httptest.Server
	mux := HTTP.NewServeMux()
	mux.HandleFunc("/items", func(w HTTP.ResponseWriter, r *HTTP.Request) {
		page := 0
		if p := r.URL.Query().Get("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil {
				page = n
			}
		}
		start := page * pageSize
		if start >= len(data) {
			HTTP.Error(w, "page out of range", HTTP.StatusNotFound)
			return
		}
		end := start + pageSize
		if end > len(data) {
			end = len(data)
		}
		resp := pageResponse{Items: data[start:end]}
		if end < len(data) {
			resp.Next = fmt.Sprintf("%s/items?page=%d", srv.URL, page+1)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
	srv = httptest.NewServer(mux)
	return srv
}

// fetchPageStep returns the Kleisli step for [thunk.Unfold].
// The seed is Some(url), the absolute URL of the page to fetch; None signals
// end-of-sequence so that [thunk.Unfold] stops iterating.
func fetchPageStep(client H.Client) Effect[Option[string], Option[Pair[Option[string], []string]]] {
	return F.Flow2(
		O.Map(F.Flow3(
			H.MakeGetRequest,
			H.ReadJSON[pageResponse](client),
			thunk.Map(F.Pipe2(
				getItems,
				reader.ApS(P.FromHead[[]string], getNext),
				reader.Map[pageResponse](O.Of[Pair[Option[string], []string]]),
			)),
		)),
		O.GetOrElse(lazy.Of(thunk.Of(O.None[Pair[Option[string], []string]]()))),
	)
}

// ExampleUnfold_pagination demonstrates lazy HTTP pagination with [thunk.Unfold].
// The embedded test server returns a "next" URL on every page except the last;
// getNext converts an empty Next field to None, which [thunk.Unfold] treats as
// the end-of-sequence signal.
func ExampleUnfold_pagination() {
	data := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta"}
	srv := newPaginatedServer(data, 3)
	defer srv.Close()

	client := H.MakeClient(HTTP.DefaultClient)

	allPagesSeq := F.Pipe2(
		thunk.Unfold(fetchPageStep(client), O.Of(srv.URL+"/items")),
		reader.Map[context.Context](F.Pipe1(
			slices.Values[[]string],
			iterresult.ChainSeqK,
		)),
		reader.Map[context.Context](iterresult.Collect[string]),
	)

	all := allPagesSeq(context.Background())()

	fmt.Println(all)

	// Output:
	// Right[[]string]([alpha beta gamma delta epsilon zeta eta])
}
