package alt

import "iter"

type (
	AltType[HKTA any] = func(func() HKTA) func(HKTA) HKTA

	Seq[T any] = iter.Seq[T]
)
