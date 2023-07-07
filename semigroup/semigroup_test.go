package semigroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {

	first := First[int]()

	assert.Equal(t, 1, first.Concat(1, 2))
}

func TestLast(t *testing.T) {

	last := Last[int]()

	assert.Equal(t, 2, last.Concat(1, 2))
}
