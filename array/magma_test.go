package array

import (
	"testing"

	"github.com/stretchr/testify/assert"

	M "github.com/IBM/fp-go/magma"
)

var subInt = M.MakeMagma(func(first int, second int) int {
	return first - second
})

func TestConcatAll(t *testing.T) {

	var subAll = ConcatAll(subInt)(0)

	assert.Equal(t, subAll([]int{1, 2, 3}), -6)

}
