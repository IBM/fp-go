// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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
	"net/http"

	H "github.com/IBM/fp-go/v2/http"
	IOEH "github.com/IBM/fp-go/v2/ioeither/http"
)

type (
	// Requester is a reader that constructs a request
	Requester = IOEH.Requester

	Client = IOEH.Client
)

var (
	// MakeRequest is an eitherized version of [http.NewRequest]
	MakeRequest = IOEH.MakeRequest

	// specialize
	MakeGetRequest = IOEH.MakeGetRequest
)

// MakeBodyRequest creates a request that carries a body
//
//go:inline
func MakeBodyRequest(method string, body IOResult[[]byte]) Kleisli[string, *http.Request] {
	return IOEH.MakeBodyRequest(method, body)
}

//go:inline
func MakeClient(httpClient *http.Client) Client {
	return IOEH.MakeClient(httpClient)
}

// ReadFullResponse sends a request,  reads the response as a byte array and represents the result as a tuple
//
//go:inline
func ReadFullResponse(client Client) Operator[*http.Request, H.FullResponse] {
	return IOEH.ReadFullResponse(client)
}

// ReadAll sends a request and reads the response as bytes
//
//go:inline
func ReadAll(client Client) Operator[*http.Request, []byte] {
	return IOEH.ReadAll(client)
}

// ReadText sends a request, reads the response and represents the response as a text string
//
//go:inline
func ReadText(client Client) Operator[*http.Request, string] {
	return IOEH.ReadText(client)
}

// ReadJSON sends a request, reads the response and parses the response as JSON
//
//go:inline
func ReadJSON[A any](client Client) Operator[*http.Request, A] {
	return IOEH.ReadJSON[A](client)
}
