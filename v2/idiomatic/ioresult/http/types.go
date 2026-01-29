package http

import (
	"net/http"

	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
)

type (
	IOResult[A any]    = ioresult.IOResult[A]
	Kleisli[A, B any]  = ioresult.Kleisli[A, B]
	Operator[A, B any] = ioresult.Operator[A, B]
	Requester          = IOResult[*http.Request]

	Client interface {
		Do(Requester) IOResult[*http.Response]
	}
)
