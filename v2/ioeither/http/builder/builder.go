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

package builder

import (
	"bytes"
	"net/http"
	"strconv"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/http/builder"
	H "github.com/IBM/fp-go/v2/http/headers"
	"github.com/IBM/fp-go/v2/ioeither"
	IOEH "github.com/IBM/fp-go/v2/ioeither/http"
	LZ "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

func Requester(builder *R.Builder) IOEH.Requester {

	withBody := F.Curry3(func(data []byte, url string, method string) IOEither[*http.Request] {
		return ioeither.TryCatchError(func() (*http.Request, error) {
			req, err := http.NewRequest(method, url, bytes.NewReader(data))
			if err == nil {
				req.Header.Set(H.ContentLength, strconv.Itoa(len(data)))
				H.Monoid.Concat(req.Header, builder.GetHeaders())
			}
			return req, err
		})
	})

	withoutBody := F.Curry2(func(url string, method string) IOEither[*http.Request] {
		return ioeither.TryCatchError(func() (*http.Request, error) {
			req, err := http.NewRequest(method, url, http.NoBody)
			if err == nil {
				H.Monoid.Concat(req.Header, builder.GetHeaders())
			}
			return req, err
		})
	})

	return F.Pipe5(
		builder.GetBody(),
		O.Fold(LZ.Of(E.Of[error](withoutBody)), E.Map[error](withBody)),
		E.Ap[func(string) IOEither[*http.Request]](builder.GetTargetURL()),
		E.Flap[error, IOEither[*http.Request]](builder.GetMethod()),
		E.GetOrElse(ioeither.Left[*http.Request, error]),
		ioeither.Map[error](func(req *http.Request) *http.Request {
			req.Header = H.Monoid.Concat(req.Header, builder.GetHeaders())
			return req
		}),
	)
}
