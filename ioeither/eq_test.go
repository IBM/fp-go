package ioeither

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEq(t *testing.T) {

	r1 := Of[string](1)
	r2 := Of[string](1)
	r3 := Of[string](2)

	e1 := Left[int]("a")
	e2 := Left[int]("a")
	e3 := Left[int]("b")

	eq := FromStrictEquals[string, int]()

	assert.True(t, eq.Equals(r1, r1))
	assert.True(t, eq.Equals(r1, r2))
	assert.False(t, eq.Equals(r1, r3))
	assert.False(t, eq.Equals(r1, e1))

	assert.True(t, eq.Equals(e1, e1))
	assert.True(t, eq.Equals(e1, e2))
	assert.False(t, eq.Equals(e1, e3))
	assert.False(t, eq.Equals(e2, r2))
}
