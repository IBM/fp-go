package http

import (
	H "net/http"

	T "github.com/IBM/fp-go/tuple"
)

type (
	// FullResponse represents a full http response, including headers and body
	FullResponse = T.Tuple2[*H.Response, []byte]
)

var (
	Response = T.First[*H.Response, []byte]
	Body     = T.Second[*H.Response, []byte]
)
