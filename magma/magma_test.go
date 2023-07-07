package magma

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	m := First[string]()

	assert.Equal(t, "a", m.Concat("a", "b"))
}

func TestSecond(t *testing.T) {
	m := Second[string]()

	assert.Equal(t, "b", m.Concat("a", "b"))
}
