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
	"bytes"
	"io"
	"net/http"

	B "github.com/IBM/fp-go/bytes"
	FL "github.com/IBM/fp-go/file"
	F "github.com/IBM/fp-go/function"
	H "github.com/IBM/fp-go/http"
	IOE "github.com/IBM/fp-go/ioeither"
	IOEF "github.com/IBM/fp-go/ioeither/file"
	J "github.com/IBM/fp-go/json"
	P "github.com/IBM/fp-go/pair"
)

type (
	// Requester is a reader that constructs a request
	Requester = IOE.IOEither[error, *http.Request]

	Client interface {
		Do(Requester) IOE.IOEither[error, *http.Response]
	}

	client struct {
		delegate *http.Client
		doIOE    func(*http.Request) IOE.IOEither[error, *http.Response]
	}
)

var (
	// MakeRequest is an eitherized version of [http.NewRequest]
	MakeRequest = IOE.Eitherize3(http.NewRequest)
	makeRequest = F.Bind13of3(MakeRequest)

	// specialize
	MakeGetRequest = makeRequest("GET", nil)
)

// MakeBodyRequest creates a request that carries a body
func MakeBodyRequest(method string, body IOE.IOEither[error, []byte]) func(url string) IOE.IOEither[error, *http.Request] {
	onBody := F.Pipe1(
		body,
		IOE.Map[error](F.Flow2(
			bytes.NewReader,
			FL.ToReader[*bytes.Reader],
		)),
	)
	onRelease := IOE.Of[error, io.Reader]
	withMethod := F.Bind1of3(MakeRequest)(method)

	return F.Flow2(
		F.Bind1of2(withMethod),
		IOE.WithResource[*http.Request](onBody, onRelease),
	)
}

func (client client) Do(req Requester) IOE.IOEither[error, *http.Response] {
	return F.Pipe1(
		req,
		IOE.Chain(client.doIOE),
	)
}

func MakeClient(httpClient *http.Client) Client {
	return client{delegate: httpClient, doIOE: IOE.Eitherize1(httpClient.Do)}
}

// ReadFullResponse sends a request,  reads the response as a byte array and represents the result as a tuple
func ReadFullResponse(client Client) func(Requester) IOE.IOEither[error, H.FullResponse] {
	return F.Flow3(
		client.Do,
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

// ReadAll sends a request and reads the response as bytes
func ReadAll(client Client) func(Requester) IOE.IOEither[error, []byte] {
	return F.Flow2(
		ReadFullResponse(client),
		IOE.Map[error](H.Body),
	)
}

// ReadText sends a request, reads the response and represents the response as a text string
func ReadText(client Client) func(Requester) IOE.IOEither[error, string] {
	return F.Flow2(
		ReadAll(client),
		IOE.Map[error](B.ToString),
	)
}

// ReadJson sends a request, reads the response and parses the response as JSON
//
// Deprecated: use [ReadJSON] instead
func ReadJson[A any](client Client) func(Requester) IOE.IOEither[error, A] {
	return ReadJSON[A](client)
}

// readJSON sends a request, reads the response and parses the response as a []byte
func readJSON(client Client) func(Requester) IOE.IOEither[error, []byte] {
	return F.Flow3(
		ReadFullResponse(client),
		IOE.ChainFirstEitherK(F.Flow2(
			H.Response,
			H.ValidateJSONResponse,
		)),
		IOE.Map[error](H.Body),
	)
}

// ReadJSON sends a request, reads the response and parses the response as JSON
func ReadJSON[A any](client Client) func(Requester) IOE.IOEither[error, A] {
	return F.Flow2(
		readJSON(client),
		IOE.ChainEitherK[error](J.Unmarshal[A]),
	)
}
