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
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	IOEF "github.com/IBM/fp-go/v2/idiomatic/ioresult/file"
	J "github.com/IBM/fp-go/v2/json"
	P "github.com/IBM/fp-go/v2/pair"
)

type (
	client struct {
		delegate *http.Client
		doIOE    Kleisli[*http.Request, *http.Response]
	}
)

var (
	// MakeRequest is an eitherized version of [http.NewRequest]
	MakeRequest = ioresult.Eitherize3(http.NewRequest)
	makeRequest = F.Bind13of3(MakeRequest)

	// specialize
	MakeGetRequest = makeRequest("GET", nil)
)

// MakeBodyRequest creates a request that carries a body
func MakeBodyRequest(method string, body IOResult[[]byte]) Kleisli[string, *http.Request] {
	onBody := F.Pipe1(
		body,
		ioresult.Map(F.Flow2(
			bytes.NewReader,
			FL.ToReader[*bytes.Reader],
		)),
	)
	onRelease := ioresult.Of[io.Reader]
	withMethod := F.Bind1of3(MakeRequest)(method)

	return F.Flow2(
		F.Bind1of2(withMethod),
		ioresult.WithResource[*http.Request](onBody, onRelease),
	)
}

func (client client) Do(req Requester) IOResult[*http.Response] {
	return F.Pipe1(
		req,
		ioresult.Chain(client.doIOE),
	)
}

func MakeClient(httpClient *http.Client) Client {
	return client{delegate: httpClient, doIOE: ioresult.Eitherize1(httpClient.Do)}
}

// ReadFullResponse sends a request,  reads the response as a byte array and represents the result as a tuple
func ReadFullResponse(client Client) Operator[*http.Request, H.FullResponse] {
	return F.Flow3(
		client.Do,
		ioresult.ChainEitherK(H.ValidateResponse),
		ioresult.Chain(func(resp *http.Response) IOResult[H.FullResponse] {
			// var x R.Reader[*http.Response, IOResult[[]byte]] = F.Flow3(
			// 	H.GetBody,
			// 	ioresult.Of,
			// 	IOEF.ReadAll,
			// )

			return F.Pipe1(
				F.Pipe3(
					resp,
					H.GetBody,
					ioresult.Of,
					IOEF.ReadAll,
				),
				ioresult.Map(F.Bind1st(P.MakePair[*http.Response, []byte], resp)),
			)
		}),
	)
}

// ReadAll sends a request and reads the response as bytes
func ReadAll(client Client) Operator[*http.Request, []byte] {
	return F.Flow2(
		ReadFullResponse(client),
		ioresult.Map(H.Body),
	)
}

// ReadText sends a request, reads the response and represents the response as a text string
func ReadText(client Client) Operator[*http.Request, string] {
	return F.Flow2(
		ReadAll(client),
		ioresult.Map(B.ToString),
	)
}

// readJSON sends a request, reads the response and parses the response as a []byte
func readJSON(client Client) Operator[*http.Request, []byte] {
	return F.Flow3(
		ReadFullResponse(client),
		ioresult.ChainFirstEitherK(F.Flow2(
			H.Response,
			H.ValidateJSONResponse,
		)),
		ioresult.Map(H.Body),
	)
}

// ReadJSON sends a request, reads the response and parses the response as JSON
func ReadJSON[A any](client Client) Operator[*http.Request, A] {
	return F.Flow2(
		readJSON(client),
		ioresult.ChainEitherK(J.Unmarshal[A]),
	)
}
