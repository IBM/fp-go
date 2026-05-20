package witherable

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

type (
	Option[A any]       = option.Option[A]
	Separated[A, B any] = pair.Pair[A, B]
)

// HKTWA = W[A]
// HKTFOB = F[Option[B]]
// HKTFWB = F[W[B]]
type WitherType[A, HKTWA, HKTFOB, HKTFWB any] = func(func(A) HKTFOB) func(HKTWA) HKTFWB

// HKTWA = W[A]
// HKTFEBC = F[Either[B, C]]
// HKTSWBC = F[Either[W[B], W[C]]]
type WiltType[A, HKTWA, HKTFEBC, HKTSWBC any] = func(func(A) HKTFEBC) func(HKTWA) HKTSWBC
