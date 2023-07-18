package eq

import (
	EQ "github.com/ibm/fp-go/eq"
	F "github.com/ibm/fp-go/function"
)

// Eq implements an equals predicate on the basis of `map` and `ap`
func Eq[HKTA, HKTABOOL, HKTBOOL, A any](
	fmap func(HKTA, func(A) func(A) bool) HKTABOOL,
	fap func(HKTABOOL, HKTA) HKTBOOL,

	e EQ.Eq[A],
) func(l, r HKTA) HKTBOOL {
	c := F.Curry2(e.Equals)
	return func(fl, fr HKTA) HKTBOOL {
		return fap(fmap(fl, c), fr)
	}
}
