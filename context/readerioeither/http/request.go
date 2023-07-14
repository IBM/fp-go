package http

import (
	"context"
	"io"
	"net/http"

	B "github.com/ibm/fp-go/bytes"
	RIOE "github.com/ibm/fp-go/context/readerioeither"
	F "github.com/ibm/fp-go/function"
	IOE "github.com/ibm/fp-go/ioeither"
	IOEF "github.com/ibm/fp-go/ioeither/file"
	J "github.com/ibm/fp-go/json"
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
	NewRequest = RIOE.Eitherize3(http.NewRequestWithContext)
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

// ReadAll sends a request and reads the response as a byte array
func ReadAll(client Client) func(Requester) RIOE.ReaderIOEither[[]byte] {
	return func(req Requester) RIOE.ReaderIOEither[[]byte] {
		doReq := client.Do(req)
		return func(ctx context.Context) IOE.IOEither[error, []byte] {
			return IOEF.ReadAll(F.Pipe2(
				ctx,
				doReq,
				IOE.Map[error](func(resp *http.Response) io.ReadCloser {
					return resp.Body
				}),
			),
			)
		}
	}
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
	return F.Flow2(
		ReadAll(client),
		RIOE.ChainEitherK(J.Unmarshal[A]),
	)
}
