package array

import (
	"testing"

	O "github.com/IBM/fp-go/ord"
	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {

	ordInt := O.FromStrictCompare[int]()

	input := []int{2, 1, 3}

	res := Sort(ordInt)(input)

	assert.Equal(t, []int{1, 2, 3}, res)
	assert.Equal(t, []int{2, 1, 3}, input)

}
