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
	"net/http"
	"net/url"
	"strconv"

	E "github.com/IBM/fp-go/either"
	ENDO "github.com/IBM/fp-go/endomorphism"
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/http/content"
	FM "github.com/IBM/fp-go/http/form"
	H "github.com/IBM/fp-go/http/headers"
	IOG "github.com/IBM/fp-go/io/generic"
	IOE "github.com/IBM/fp-go/ioeither"
	IOEH "github.com/IBM/fp-go/ioeither/http"
	J "github.com/IBM/fp-go/json"
	LZ "github.com/IBM/fp-go/lazy"
	L "github.com/IBM/fp-go/optics/lens"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	T "github.com/IBM/fp-go/tuple"
)

type (
	Builder struct {
		method  O.Option[string]
		url     string
		headers http.Header
		body    O.Option[IOE.IOEither[error, []byte]]
		query   url.Values
	}

	// Endomorphism returns an [ENDO.Endomorphism] that transforms a builder
	Endomorphism = ENDO.Endomorphism[*Builder]
)

var (
	// Default is the default builder
	Default = &Builder{method: O.Some(defaultMethod()), headers: make(http.Header), body: noBody}

	defaultMethod = F.Constant(http.MethodGet)

	// Monoid is the [M.Monoid] for the [Endomorphism]
	Monoid = ENDO.Monoid[*Builder]()

	// Url is a [L.Lens] for the URL
	Url = L.MakeLensRef((*Builder).GetUrl, (*Builder).SetUrl)
	// Method is a [L.Lens] for the HTTP method
	Method = L.MakeLensRef((*Builder).GetMethod, (*Builder).SetMethod)
	// Body is a [L.Lens] for the request body
	Body = L.MakeLensRef((*Builder).GetBody, (*Builder).SetBody)
	// Headers is a [L.Lens] for the complete set of request headers
	Headers = L.MakeLensRef((*Builder).GetHeaders, (*Builder).SetHeaders)
	// Query is a [L.Lens] for the set of query parameters
	Query = L.MakeLensRef((*Builder).GetQuery, (*Builder).SetQuery)

	rawQuery = L.MakeLensRef(getRawQuery, setRawQuery)

	getHeader = F.Bind2of2((*Builder).GetHeader)
	delHeader = F.Bind2of2((*Builder).DelHeader)
	setHeader = F.Bind2of3((*Builder).SetHeader)

	noHeader   = O.None[string]()
	noBody     = O.None[IOE.IOEither[error, []byte]]()
	noQueryArg = O.None[string]()

	parseUrl   = E.Eitherize1(url.Parse)
	parseQuery = E.Eitherize1(url.ParseQuery)

	// WithQuery creates a [Endomorphism] for a complete set of query parameters
	WithQuery = Query.Set
	// WithMethod creates a [Endomorphism] for a certain method
	WithMethod = Method.Set
	// WithUrl creates a [Endomorphism] for a certain method
	WithUrl = Url.Set
	// WithHeaders creates a [Endomorphism] for a set of headers
	WithHeaders = Headers.Set
	// WithBody creates a [Endomorphism] for a request body
	WithBody = F.Flow2(
		O.Of[IOE.IOEither[error, []byte]],
		Body.Set,
	)
	// WithBytes creates a [Endomorphism] for a request body using bytes
	WithBytes = F.Flow2(
		IOE.Of[error, []byte],
		WithBody,
	)
	// WithContentType adds the [H.ContentType] header
	WithContentType = WithHeader(H.ContentType)
	// WithAuthorization adds the [H.Authorization] header
	WithAuthorization = WithHeader(H.Authorization)

	// WithGet adds the [http.MethodGet] method
	WithGet = WithMethod(http.MethodGet)
	// WithPost adds the [http.MethodPost] method
	WithPost = WithMethod(http.MethodPost)
	// WithPut adds the [http.MethodPut] method
	WithPut = WithMethod(http.MethodPut)
	// WithDelete adds the [http.MethodDelete] method
	WithDelete = WithMethod(http.MethodDelete)

	// WithBearer creates a [Endomorphism] to add a Bearer [H.Authorization] header
	WithBearer = F.Flow2(
		S.Format[string]("Bearer %s"),
		WithAuthorization,
	)

	// Requester creates a requester from a builder
	Requester = (*Builder).Requester

	// WithoutBody creates a [Endomorphism] to remove the body
	WithoutBody = F.Pipe1(
		noBody,
		Body.Set,
	)
)

func setRawQuery(u *url.URL, raw string) *url.URL {
	u.RawQuery = raw
	return u
}

func getRawQuery(u *url.URL) string {
	return u.RawQuery
}

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

func (builder *Builder) GetQuery() url.Values {
	return builder.query
}

func (builder *Builder) SetQuery(query url.Values) *Builder {
	builder.query = query
	return builder
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

func (builder *Builder) Requester() IOEH.Requester {

	withBody := F.Curry3(func(data []byte, url string, method string) IOE.IOEither[error, *http.Request] {
		return IOE.TryCatchError(func() (*http.Request, error) {
			req, err := http.NewRequest(method, url, bytes.NewReader(data))
			if err == nil {
				req.Header.Set(H.ContentLength, strconv.Itoa(len(data)))
				H.Monoid.Concat(req.Header, builder.headers)
			}
			return req, err
		})
	})

	withoutBody := F.Curry2(func(url string, method string) IOE.IOEither[error, *http.Request] {
		return IOE.TryCatchError(func() (*http.Request, error) {
			req, err := http.NewRequest(method, url, nil)
			if err == nil {
				H.Monoid.Concat(req.Header, builder.headers)
			}
			return req, err
		})
	})

	// construct the final URL
	targetUrl := F.Pipe3(
		builder,
		Url.Get,
		parseUrl,
		E.Chain(F.Flow4(
			T.Replicate2[*url.URL],
			T.Map2(
				F.Flow2(
					F.Curry2(setRawQuery),
					E.Of[error, func(string) *url.URL],
				),
				F.Flow3(
					rawQuery.Get,
					parseQuery,
					E.Map[error](F.Flow2(
						F.Curry2(FM.ValuesMonoid.Concat)(builder.GetQuery()),
						(url.Values).Encode,
					)),
				),
			),
			T.Tupled2(E.MonadAp[*url.URL, error, string]),
			E.Map[error]((*url.URL).String),
		)),
	)

	return F.Pipe5(
		builder,
		Body.Get,
		O.Fold(LZ.Of(IOE.Of[error](withoutBody)), IOE.Map[error](withBody)),
		IOG.Map[IOE.IOEither[error, func(string) func(string) IOE.IOEither[error, *http.Request]], IOE.IOEither[error, func(string) IOE.IOEither[error, *http.Request]]](E.Ap[func(string) IOE.IOEither[error, *http.Request]](targetUrl)),
		IOE.Flap[error, IOE.IOEither[error, *http.Request]](builder.GetMethod()),
		IOE.Flatten[error, *http.Request],
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

	return L.MakeLens(get, func(b *Builder, value O.Option[string]) *Builder {
		cpy := b.clone()
		return F.Pipe1(
			value,
			O.Fold(del(cpy), set(cpy)),
		)
	})
}

// WithHeader creates a [Endomorphism] for a certain header
func WithHeader(name string) func(value string) Endomorphism {
	return F.Flow2(
		O.Of[string],
		Header(name).Set,
	)
}

// WithoutHeader creates a [Endomorphism] to remove a certain header
func WithoutHeader(name string) Endomorphism {
	return Header(name).Set(noHeader)
}

// WithFormData creates a [Endomorphism] to send form data payload
func WithFormData(value url.Values) Endomorphism {
	return F.Flow2(
		F.Pipe4(
			value,
			url.Values.Encode,
			S.ToBytes,
			IOE.Of[error, []byte],
			WithBody,
		),
		WithContentType(C.FormEncoded),
	)
}

// WithJson creates a [Endomorphism] to send JSON payload
func WithJson[T any](data T) Endomorphism {
	return F.Flow2(
		F.Pipe3(
			data,
			J.Marshal[T],
			IOE.FromEither[error, []byte],
			WithBody,
		),
		WithContentType(C.Json),
	)
}

// QueryArg is a [L.Lens] for the first value of a query argument
func QueryArg(name string) L.Lens[*Builder, O.Option[string]] {
	return F.Pipe1(
		Query,
		L.Compose[*Builder](FM.AtValue(name)),
	)
}

// WithQueryArg creates a [Endomorphism] for a certain query argument
func WithQueryArg(name string) func(value string) Endomorphism {
	return F.Flow2(
		O.Of[string],
		QueryArg(name).Set,
	)
}

// WithoutQueryArg creates a [Endomorphism] that removes a query argument
func WithoutQueryArg(name string) Endomorphism {
	return QueryArg(name).Set(noQueryArg)
}
