package semigroup

import (
	F "github.com/IBM/fp-go/function"
)

/*
*
HKTA = HKT<A>
HKTFA = HKT<func(A)A>
*/
func ApplySemigroup[A, HKTA, HKTFA any](
	_map func(HKTA, func(A) func(A) A) HKTFA,
	_ap func(HKTFA, HKTA) HKTA,

	s Semigroup[A],
) Semigroup[HKTA] {

	cb := F.Curry2(s.Concat)
	return MakeSemigroup(func(first HKTA, second HKTA) HKTA {
		return _ap(_map(first, cb), second)
	})
}
