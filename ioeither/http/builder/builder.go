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

package builder

import (
	"bytes"
	"io"
	"net/http"

	FL "github.com/IBM/fp-go/file"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	IOEH "github.com/IBM/fp-go/ioeither/http"
	LZ "github.com/IBM/fp-go/lazy"
	L "github.com/IBM/fp-go/optics/lens"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
)

type (
	Builder struct {
		method  O.Option[string]
		url     string
		headers http.Header
		body    O.Option[IOE.IOEither[error, []byte]]
	}

	BuilderBuilder = func(*Builder) *Builder
)

var (
	// Default is the default builder
	Default = &Builder{method: O.Some(defaultMethod()), headers: make(http.Header), body: noBody}

	defaultMethod = F.Constant(http.MethodGet)

	// Url is a [L.Lens] for the URL
	Url = L.MakeLensRef((*Builder).GetUrl, (*Builder).SetUrl)
	// Method is a [L.Lens] for the HTTP method
	Method = L.MakeLensRef((*Builder).GetMethod, (*Builder).SetMethod)
	// Body is a [L.Lens] for the request body
	Body = L.MakeLensRef((*Builder).GetBody, (*Builder).SetBody)
	// Headers is a [L.Lens] for the complete set of request headers
	Headers = L.MakeLensRef((*Builder).GetHeaders, (*Builder).SetHeaders)

	getHeader = F.Bind2of2((*Builder).GetHeader)
	delHeader = F.Bind2of2((*Builder).DelHeader)
	setHeader = F.Bind2of3((*Builder).SetHeader)

	noHeader = O.None[string]()
	noBody   = O.None[IOE.IOEither[error, []byte]]()

	// WithMethod creates a [BuilderBuilder] for a certain method
	WithMethod = Method.Set
	// WithUrl creates a [BuilderBuilder] for a certain method
	WithUrl = Url.Set
	// WithHeaders creates a [BuilderBuilder] for a set of headers
	WithHeaders = Headers.Set
	// WithBody creates a [BuilderBuilder] for a request body
	WithBody = F.Flow2(
		O.Of[IOE.IOEither[error, []byte]],
		Body.Set,
	)
	// WithContentType adds the content type header
	WithContentType = WithHeader("Content-Type")

	// WithGet adds the [http.MethodGet] method
	WithGet = WithMethod(http.MethodGet)
	// WithPost adds the [http.MethodPost] method
	WithPost = WithMethod(http.MethodPost)
	// WithPut adds the [http.MethodPut] method
	WithPut = WithMethod(http.MethodPut)
	// WithDelete adds the [http.MethodDelete] method
	WithDelete = WithMethod(http.MethodDelete)

	// Requester creates a requester from a builder
	Requester = (*Builder).Requester

	// WithoutBody creates a [BuilderBuilder] to remove the body
	WithoutBody = Body.Set(noBody)
)

func (builder *Builder) clone() *Builder {
	cpy := *builder
	cpy.headers = cpy.headers.Clone()
	return &cpy
}

func (builder *Builder) GetUrl() string {
	return builder.url
}

func (builder *Builder) GetMethod() string {
	return F.Pipe1(
		builder.method,
		O.GetOrElse(defaultMethod),
	)
}

func (builder *Builder) GetHeaders() http.Header {
	return builder.headers
}

func (builder *Builder) GetBody() O.Option[IOE.IOEither[error, []byte]] {
	return builder.body
}

func (builder *Builder) SetMethod(method string) *Builder {
	builder.method = O.Some(method)
	return builder
}

func (builder *Builder) SetUrl(url string) *Builder {
	builder.url = url
	return builder
}

func (builder *Builder) SetHeaders(headers http.Header) *Builder {
	builder.headers = headers
	return builder
}

func (builder *Builder) SetBody(body O.Option[IOE.IOEither[error, []byte]]) *Builder {
	builder.body = body
	return builder
}

func (builder *Builder) SetHeader(name, value string) *Builder {
	builder.headers.Set(name, value)
	return builder
}

func (builder *Builder) DelHeader(name string) *Builder {
	builder.headers.Del(name)
	return builder
}

func (builder *Builder) GetHeader(name string) O.Option[string] {
	return F.Pipe2(
		name,
		builder.headers.Get,
		O.FromPredicate(S.IsNonEmpty),
	)
}

func (builder *Builder) GetHeaderValues(name string) []string {
	return builder.headers.Values(name)
}

func (builder *Builder) AddHeaderHeader(name, value string) *Builder {
	builder.headers.Add(name, value)
	return builder
}

func (builder *Builder) Requester() IOEH.Requester {
	return F.Pipe3(
		builder.GetBody(),
		O.Map(IOE.Map[error](F.Flow2(
			bytes.NewReader,
			FL.ToReader[*bytes.Reader],
		))),
		O.GetOrElse(F.Constant(IOE.Of[error, io.Reader](nil))),
		IOE.Chain(func(rdr io.Reader) IOE.IOEither[error, *http.Request] {
			return IOE.TryCatchError(func() (*http.Request, error) {
				req, err := http.NewRequest(builder.GetMethod(), builder.GetUrl(), rdr)
				if err == nil {
					for name, value := range builder.GetHeaders() {
						req.Header[name] = value
					}
				}
				return req, err
			})
		}),
	)
}

// Header returns a [L.Lens] for a single header
func Header(name string) L.Lens[*Builder, O.Option[string]] {
	get := getHeader(name)
	set := F.Bind1of2(setHeader(name))
	del := F.Flow2(
		LZ.Of[*Builder],
		LZ.Map(delHeader(name)),
	)

	return L.MakeLens[*Builder, O.Option[string]](get, func(b *Builder, value O.Option[string]) *Builder {
		cpy := b.clone()
		return F.Pipe1(
			value,
			O.Fold(del(cpy), set(cpy)),
		)
	})
}

// WithHeader creates a [BuilderBuilder] for a certain header
func WithHeader(name string) func(value string) BuilderBuilder {
	return F.Flow2(
		O.Of[string],
		Header(name).Set,
	)
}

// WithoutHeader creates a [BuilderBuilder] to remove a certain header
func WithoutHeader(name string) BuilderBuilder {
	return Header(name).Set(noHeader)
}
