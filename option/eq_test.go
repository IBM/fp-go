package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {

	r1 := Of(1)
	r2 := Of(1)
	r3 := Of(2)

	n1 := None[int]()

	eq := FromStrictEquals[int]()

	assert.True(t, eq.Equals(r1, r1))
	assert.True(t, eq.Equals(r1, r2))
	assert.False(t, eq.Equals(r1, r3))
	assert.False(t, eq.Equals(r1, n1))

	assert.True(t, eq.Equals(n1, n1))
	assert.False(t, eq.Equals(n1, r2))
}
