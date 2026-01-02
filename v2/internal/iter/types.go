package iter

import (
	I "iter"
)

type (
	// Seq represents Go's standard library iterator type for single values.
	// It's an alias for iter.Seq[A] and provides interoperability with Go 1.23+ range-over-func.
	Seq[A any] = I.Seq[A]
)
