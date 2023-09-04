// Copyright (c) 2023 IBM Corp.
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

	R "github.com/IBM/fp-go/context/readerioeither"
	E "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/stretchr/testify/assert"
)

type PostItem struct {
	UserId uint   `json:"userId"`
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

func (b simpleRequestBuilder) Build() R.ReaderIOEither[*H.Request] {
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

	readItem := ReadJson[PostItem](client)

	resp1 := readItem(req1)

	resE := resp1(context.TODO())()

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

	readItem := ReadJson[PostItem](client)

	resp1 := F.Pipe2(
		req1,
		readItem,
		R.Map(getTitle),
	)

	res := F.Pipe1(
		resp1(context.TODO())(),
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

	readItem := ReadJson[PostItem](client)

	response := F.Pipe2(
		request,
		readItem,
		R.Map(getTitle),
	)

	res := F.Pipe1(
		response(context.TODO())(),
		E.GetOrElse(errors.ToString),
	)

	assert.Equal(t, "sunt aut facere repellat provident occaecati excepturi optio reprehenderit", res)
}
