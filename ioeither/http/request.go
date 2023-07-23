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
	"io"
	"net/http"

	B "github.com/IBM/fp-go/bytes"
	ER "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	IOEF "github.com/IBM/fp-go/ioeither/file"
	J "github.com/IBM/fp-go/json"
)

type Client interface {
	Do(req *http.Request) IOE.IOEither[error, *http.Response]
}

type client struct {
	delegate *http.Client
}

func (client client) Do(req *http.Request) IOE.IOEither[error, *http.Response] {
	return IOE.TryCatch(func() (*http.Response, error) {
		return client.delegate.Do(req)
	}, ER.IdentityError)
}

func MakeClient(httpClient *http.Client) Client {
	return client{delegate: httpClient}
}

func ReadAll(client Client) func(*http.Request) IOE.IOEither[error, []byte] {
	return func(req *http.Request) IOE.IOEither[error, []byte] {
		return IOEF.ReadAll(F.Pipe2(
			req,
			client.Do,
			IOE.Map[error](func(resp *http.Response) io.ReadCloser {
				return resp.Body
			}),
		),
		)
	}
}

func ReadText(client Client) func(*http.Request) IOE.IOEither[error, string] {
	return F.Flow2(
		ReadAll(client),
		IOE.Map[error](B.ToString),
	)
}

func ReadJson[A any](client Client) func(*http.Request) IOE.IOEither[error, A] {
	return F.Flow2(
		ReadAll(client),
		IOE.ChainEitherK(J.Unmarshal[A]),
	)
}
