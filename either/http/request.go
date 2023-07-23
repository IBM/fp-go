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
	"net/http"

	E "github.com/IBM/fp-go/either"
)

var (
	PostRequest = bodyRequest("POST")
	PutRequest  = bodyRequest("PUT")

	GetRequest     = noBodyRequest("GET")
	DeleteRequest  = noBodyRequest("DELETE")
	OptionsRequest = noBodyRequest("OPTIONS")
	HeadRequest    = noBodyRequest("HEAD")
)

func bodyRequest(method string) func(string) func([]byte) E.Either[error, *http.Request] {
	return func(url string) func([]byte) E.Either[error, *http.Request] {
		return func(body []byte) E.Either[error, *http.Request] {
			return E.TryCatchError(func() (*http.Request, error) {
				return http.NewRequest(method, url, bytes.NewReader(body))
			})
		}
	}
}

func noBodyRequest(method string) func(string) E.Either[error, *http.Request] {
	return func(url string) E.Either[error, *http.Request] {
		return E.TryCatchError(func() (*http.Request, error) {
			return http.NewRequest(method, url, nil)
		})
	}
}
