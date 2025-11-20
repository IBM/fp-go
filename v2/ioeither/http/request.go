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
	"bytes"
	"io"
	"net/http"

	B "github.com/IBM/fp-go/v2/bytes"
	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	H "github.com/IBM/fp-go/v2/http"
	"github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	J "github.com/IBM/fp-go/v2/json"
	R "github.com/IBM/fp-go/v2/reader"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

type (
	// Requester is a reader that constructs a request
	Requester = ioeither.IOEither[error, *http.Request]

	Client interface {
		Do(Requester) ioeither.IOEither[error, *http.Response]
	}

	client struct {
		delegate *http.Client
		doIOE    Kleisli[error, *http.Request, *http.Response]
	}
)

var (
	// MakeRequest is an eitherized version of [http.NewRequest]
	MakeRequest = ioeither.Eitherize3(http.NewRequest)
	makeRequest = F.Bind13of3(MakeRequest)

	// specialize
	MakeGetRequest = makeRequest("GET", nil)
)

// MakeBodyRequest creates a request that carries a body
func MakeBodyRequest(method string, body ioeither.IOEither[error, []byte]) Kleisli[error, string, *http.Request] {
	onBody := F.Pipe1(
		body,
		ioeither.Map[error](F.Flow2(
			bytes.NewReader,
			FL.ToReader[*bytes.Reader],
		)),
	)
	onRelease := ioeither.Of[error, io.Reader]
	withMethod := F.Bind1of3(MakeRequest)(method)

	return F.Flow2(
		F.Bind1of2(withMethod),
		ioeither.WithResource[*http.Request](onBody, onRelease),
	)
}

func (client client) Do(req Requester) ioeither.IOEither[error, *http.Response] {
	return F.Pipe1(
		req,
		ioeither.Chain(client.doIOE),
	)
}

func MakeClient(httpClient *http.Client) Client {
	return client{delegate: httpClient, doIOE: ioeither.Eitherize1(httpClient.Do)}
}

// ReadFullResponse sends a request,  reads the response as a byte array and represents the result as a tuple
func ReadFullResponse(client Client) Kleisli[error, Requester, H.FullResponse] {
	return F.Flow3(
		client.Do,
		ioeither.ChainEitherK(H.ValidateResponse),
		ioeither.Chain(F.Pipe3(
			H.GetBody,
			RIOE.FromReader[error],
			R.Map[*http.Response](IOEF.ReadAll[io.ReadCloser]),
			RIOE.ChainReaderK[error](H.FromBody),
		)),
	)
}

// ReadAll sends a request and reads the response as bytes
func ReadAll(client Client) Kleisli[error, Requester, []byte] {
	return F.Flow2(
		ReadFullResponse(client),
		ioeither.Map[error](H.Body),
	)
}

// ReadText sends a request, reads the response and represents the response as a text string
func ReadText(client Client) Kleisli[error, Requester, string] {
	return F.Flow2(
		ReadAll(client),
		ioeither.Map[error](B.ToString),
	)
}

// ReadJson sends a request, reads the response and parses the response as JSON
//
// Deprecated: use [ReadJSON] instead
func ReadJson[A any](client Client) Kleisli[error, Requester, A] {
	return ReadJSON[A](client)
}

// readJSON sends a request, reads the response and parses the response as a []byte
func readJSON(client Client) Kleisli[error, Requester, []byte] {
	return F.Flow3(
		ReadFullResponse(client),
		ioeither.ChainFirstEitherK(F.Flow2(
			H.Response,
			H.ValidateJSONResponse,
		)),
		ioeither.Map[error](H.Body),
	)
}

// ReadJSON sends a request, reads the response and parses the response as JSON
func ReadJSON[A any](client Client) Kleisli[error, Requester, A] {
	return F.Flow2(
		readJSON(client),
		ioeither.ChainEitherK(J.Unmarshal[A]),
	)
}
