package io

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toUpper(s string) IO[string] {
	return Of(strings.ToUpper(s))
}

func TestTraverseArray(t *testing.T) {

	src := []string{"a", "b"}

	trv := TraverseArray(toUpper)

	res := trv(src)

	assert.Equal(t, res(), []string{"A", "B"})
}

type (
	customSlice []string
)

func TestTraverseCustomSlice(t *testing.T) {

	src := customSlice{"a", "b"}

	trv := TraverseArray(toUpper)

	res := trv(src)

	assert.Equal(t, res(), []string{"A", "B"})
}
