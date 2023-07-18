package number

import (
	M "github.com/IBM/fp-go/monoid"
)

func MonoidSum[A int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128]() M.Monoid[A] {
	s := SemigroupSum[A]()
	return M.MakeMonoid(
		s.Concat,
		0,
	)
}
