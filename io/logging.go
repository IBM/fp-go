package io

import (
	"log"

	G "github.com/IBM/fp-go/io/generic"
)

// Logger constructs a logger function that can be used with ChainXXXIOK
func Logger[A any](loggers ...*log.Logger) func(string) func(A) IO[any] {
	return G.Logger[IO[any], A](loggers...)
}

// Logf constructs a logger function that can be used with ChainXXXIOK
// the string prefix contains the format string for the log value
func Logf[A any](prefix string) func(A) IO[any] {
	return G.Logf[IO[any], A](prefix)
}
