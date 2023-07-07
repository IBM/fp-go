package array

import (
	"testing"

	M "github.com/ibm/fp-go/monoid/testing"
)

func TestMonoid(t *testing.T) {
	M.AssertLaws(t, Monoid[int]())([][]int{{}, {1}, {1, 2}})
}
