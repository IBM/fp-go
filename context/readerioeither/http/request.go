package http

import (
	"io"
	"net/http"

	B "github.com/ibm/fp-go/bytes"
	RIOE "github.com/ibm/fp-go/context/readerioeither"
	F "github.com/ibm/fp-go/function"
	H "github.com/ibm/fp-go/http"
	IOE "github.com/ibm/fp-go/ioeither"
	IOEF "github.com/ibm/fp-go/ioeither/file"
	J "github.com/ibm/fp-go/json"
	T "github.com/ibm/fp-go/tuple"
)

type (
	// Requester is a reader that constructs a request
	Requester = RIOE.ReaderIOEither[*http.Request]

	Client interface {
		// Do can send an HTTP request considering a context
		Do(Requester) RIOE.ReaderIOEither[*http.Response]
	}

	client struct {
		delegate *http.Client
		doIOE    func(*http.Request) IOE.IOEither[error, *http.Response]
	}
)

var (
	// MakeRequest is an eitherized version of [http.NewRequestWithContext]
	MakeRequest = RIOE.Eitherize3(http.NewRequestWithContext)
	makeRequest = F.Bind13of3(MakeRequest)

	// specialize
	MakeGetRequest = makeRequest("GET", nil)
)

func (client client) Do(req Requester) RIOE.ReaderIOEither[*http.Response] {
	return F.Pipe1(
		req,
		RIOE.ChainIOEitherK(client.doIOE),
	)
}

// MakeClient creates an HTTP client proxy
func MakeClient(httpClient *http.Client) Client {
	return client{delegate: httpClient, doIOE: IOE.Eitherize1(httpClient.Do)}
}

// ReadFullResponse sends a request,  reads the response as a byte array and represents the result as a tuple
func ReadFullResponse(client Client) func(Requester) RIOE.ReaderIOEither[H.FullResponse] {
	return func(req Requester) RIOE.ReaderIOEither[H.FullResponse] {
		return F.Flow3(
			client.Do(req),
			IOE.ChainEitherK(H.ValidateResponse),
			IOE.Chain(func(resp *http.Response) IOE.IOEither[error, H.FullResponse] {
				return F.Pipe1(
					F.Pipe3(
						resp,
						H.GetBody,
						IOE.Of[error, io.ReadCloser],
						IOEF.ReadAll[io.ReadCloser],
					),
					IOE.Map[error](F.Bind1st(T.MakeTuple2[*http.Response, []byte], resp)),
				)
			}),
		)
	}
}

// func test(resp *http.Response) IOE.IOEither[error, H.FullResponse] {
// 	x := F.Pipe3(
// 		resp,
// 		T.Replicate2[*http.Response],
// 		T.Map2(
// 			IOE.Of[error, *http.Response],
// 			F.Flow3(
// 				H.GetBody,
// 				IOE.Of[error, io.ReadCloser],
// 				IOEF.ReadAll[io.ReadCloser],
// 			),
// 		),
// 	)
// }

// ReadAll sends a request and reads the response as bytes
func ReadAll(client Client) func(Requester) RIOE.ReaderIOEither[[]byte] {
	return F.Flow2(
		ReadFullResponse(client),
		RIOE.Map(H.Body),
	)
}

// ReadText sends a request, reads the response and represents the response as a text string
func ReadText(client Client) func(Requester) RIOE.ReaderIOEither[string] {
	return F.Flow2(
		ReadAll(client),
		RIOE.Map(B.ToString),
	)
}

// ReadJson sends a request, reads the response and parses the response as JSON
func ReadJson[A any](client Client) func(Requester) RIOE.ReaderIOEither[A] {
	return F.Flow3(
		ReadFullResponse(client),
		RIOE.ChainFirstEitherK(F.Flow2(
			H.Response,
			H.ValidateJsonResponse,
		)),
		RIOE.ChainEitherK(F.Flow2(
			H.Body,
			J.Unmarshal[A],
		)),
	)
}
