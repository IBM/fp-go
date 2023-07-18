package array

import (
	"testing"

	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/assert"
)

type ArrayType = []int

func TestTraverse(t *testing.T) {

	traverse := Traverse(
		O.Of[ArrayType],
		O.Map[ArrayType, func(int) ArrayType],
		O.Ap[ArrayType, int],

		func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.None[int]()
			}
			return O.Of(n)
		})

	assert.Equal(t, O.None[[]int](), traverse(ArrayType{1, 2}))
	assert.Equal(t, O.Of(ArrayType{1, 3}), traverse(ArrayType{1, 3}))
}
