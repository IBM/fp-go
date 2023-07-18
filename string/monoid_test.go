package string

import (
	"testing"

	M "github.com/ibm/fp-go/monoid/testing"
)

func TestMonoid(t *testing.T) {
	M.AssertLaws(t, Monoid)([]string{"", "a", "some value"})
}
