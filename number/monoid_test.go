package number

import (
	"testing"

	M "github.com/IBM/fp-go/monoid/testing"
)

func TestMonoidSum(t *testing.T) {
	M.AssertLaws(t, MonoidSum[int]())([]int{0, 1, 1000, -1})
}
