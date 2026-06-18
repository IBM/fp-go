package io

import (
	"math/rand/v2"
)

func Int() IO[int] {
	return rand.Int
}

func IntN(n int) IO[int] {
	return func() int {
		return rand.IntN(n)
	}
}
