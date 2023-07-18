package Http

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
