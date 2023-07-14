package bytes

import (
	"testing"

	M "github.com/ibm/fp-go/monoid/testing"
)

func TestMonoid(t *testing.T) {
	M.AssertLaws(t, Monoid)([][]byte{[]byte(""), []byte("a"), []byte("some value")})
}
