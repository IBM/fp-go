package alt

import (
	F "github.com/IBM/fp-go/v2/function"
)

func AltAllArray[HKTA any](
	falt AltType[HKTA],
) func(HKTA) func([]HKTA) HKTA {
	return func(startWith HKTA) func([]HKTA) HKTA {
		return func(as []HKTA) HKTA {
			current := startWith
			for _, next := range as {
				current = falt(F.Constant(next))(current)
			}
			return current
		}
	}
}

func AltAllSeq[HKTA any](
	falt AltType[HKTA],
) func(HKTA) func(Seq[HKTA]) HKTA {
	return func(startWith HKTA) func(Seq[HKTA]) HKTA {
		return func(as Seq[HKTA]) HKTA {
			current := startWith
			for next := range as {
				current = falt(F.Constant(next))(current)
			}
			return current
		}
	}
}
