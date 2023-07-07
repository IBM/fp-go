package number

import (
	M "github.com/ibm/fp-go/magma"
)

func MagmaSub[A int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128]() M.Magma[A] {
	return M.MakeMagma(func(first A, second A) A {
		return first - second
	})
}
