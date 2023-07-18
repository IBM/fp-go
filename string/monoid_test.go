package string

import (
	"testing"

	M "github.com/IBM/fp-go/monoid/testing"
)

func TestMonoid(t *testing.T) {
	M.AssertLaws(t, Monoid)([]string{"", "a", "some value"})
}
