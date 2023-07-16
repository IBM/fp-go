package http

import (
	"testing"

	E "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func NoError[A any](t *testing.T) func(E.Either[error, A]) bool {
	return E.Fold(func(err error) bool {
		return assert.NoError(t, err)
	}, F.Constant1[A](true))
}

func Error[A any](t *testing.T) func(E.Either[error, A]) bool {
	return E.Fold(F.Constant1[error](true), func(A) bool {
		return assert.Error(t, nil)
	})
}

func TestValidateJsonContentTypeString(t *testing.T) {

	res := F.Pipe1(
		validateJsonContentTypeString("application/json"),
		NoError[ParsedMediaType](t),
	)

	assert.True(t, res)
}

func TestValidateInvalidJsonContentTypeString(t *testing.T) {

	res := F.Pipe1(
		validateJsonContentTypeString("application/xml"),
		Error[ParsedMediaType](t),
	)

	assert.True(t, res)
}
