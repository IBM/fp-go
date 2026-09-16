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
	"io"
	"net/http"
	"strconv"

	E "github.com/IBM/fp-go/v2/either"
	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	HTTP "github.com/IBM/fp-go/v2/http"
	R "github.com/IBM/fp-go/v2/http/builder"
	H "github.com/IBM/fp-go/v2/http/headers"
	"github.com/IBM/fp-go/v2/ioeither"
	IOEH "github.com/IBM/fp-go/v2/ioeither/http"
	LZ "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

func Requester(builder *R.Builder) IOEH.Requester {
	return F.Pipe4(
		builder.GetBody(),
		O.Fold(LZ.Of(noBody), E.Map[error](withBody)),
		E.Ap[IOEither[*http.Request]](F.Pipe1(
			builder.GetTargetURL(),
			E.Map[error](F.Curry2(F.Bind1of3(IOEH.MakeRequest)(builder.GetMethod()))),
		)),
		E.GetOrElse(ioeither.Left[*http.Request, error]),
		ioeither.Map[error](F.Bind1of2(mergeHeaders)(builder.GetHeaders())),
	)
}

var (
	// toReader converts a byte slice into a fresh [io.Reader]
	toReader = F.Flow2(
		bytes.NewReader,
		FL.ToReader[*bytes.Reader],
	)

	// noBody feeds an empty body into the request constructor
	noBody = E.Of[error](withReader(ioeither.Of[error, io.Reader](http.NoBody)))
)

// withReader feeds the body into a request constructor that is still waiting for its body.
func withReader(body IOEither[io.Reader]) func(ioeither.Kleisli[error, io.Reader, *http.Request]) IOEither[*http.Request] {
	return F.Bind1st(ioeither.MonadChain[error, io.Reader, *http.Request], body)
}

// withBody feeds data as the body into a request constructor and sets the Content-Length header.
// The [io.Reader] over data is created on each execution, so the resulting requester can be
// executed repeatedly (e.g. when retrying).
func withBody(data []byte) func(ioeither.Kleisli[error, io.Reader, *http.Request]) IOEither[*http.Request] {
	return F.Flow2(
		withReader(F.Pipe1(
			ioeither.Of[error](data),
			ioeither.Map[error](toReader),
		)),
		ioeither.Map[error](HTTP.WithHeader(H.ContentLength)(strconv.Itoa(len(data)))),
	)
}

// mergeHeaders merges a copy of headers into the headers of the request and returns the request.
// Copying ensures that the request never shares header maps or value slices with the builder.
func mergeHeaders(headers http.Header, req *http.Request) *http.Request {
	req.Header = H.Monoid.Concat(req.Header, headers.Clone())
	return req
}
