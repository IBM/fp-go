package record

import (
	"testing"

	S "github.com/ibm/fp-go/string"
	"github.com/stretchr/testify/assert"
)

func TestUnionMonoid(t *testing.T) {
	m := UnionMonoid[string](S.Semigroup())

	e := Empty[string, string]()

	x := map[string]string{
		"a": "a1",
		"b": "b1",
		"c": "c1",
	}

	y := map[string]string{
		"b": "b2",
		"c": "c2",
		"d": "d2",
	}

	res := map[string]string{
		"a": "a1",
		"b": "b1b2",
		"c": "c1c2",
		"d": "d2",
	}

	assert.Equal(t, x, m.Concat(x, m.Empty()))
	assert.Equal(t, x, m.Concat(m.Empty(), x))

	assert.Equal(t, x, m.Concat(x, e))
	assert.Equal(t, x, m.Concat(e, x))

	assert.Equal(t, res, m.Concat(x, y))
}
