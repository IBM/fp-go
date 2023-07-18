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
