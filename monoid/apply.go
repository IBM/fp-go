package monoid

import (
	S "github.com/IBM/fp-go/semigroup"
)

func ApplicativeMonoid[A, HKTA, HKTFA any](
	_of func(A) HKTA,
	_map func(HKTA, func(A) func(A) A) HKTFA,
	_ap func(HKTFA, HKTA) HKTA,

	m Monoid[A],
) Monoid[HKTA] {

	return MakeMonoid(
		S.ApplySemigroup[A](_map, _ap, m).Concat,
		_of(m.Empty()),
	)
}
