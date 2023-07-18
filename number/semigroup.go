package number

import (
	S "github.com/ibm/fp-go/semigroup"
)

func SemigroupSum[A int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128]() S.Semigroup[A] {
	return S.MakeSemigroup(func(first A, second A) A {
		return first + second
	})
}
