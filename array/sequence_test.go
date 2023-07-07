package array

import (
	"testing"

	"github.com/stretchr/testify/assert"

	O "github.com/ibm/fp-go/option"
)

func TestSequenceOption(t *testing.T) {
	seq := ArrayOption[int]()

	assert.Equal(t, O.Of([]int{1, 3}), seq([]O.Option[int]{O.Of(1), O.Of(3)}))
	assert.Equal(t, O.None[[]int](), seq([]O.Option[int]{O.Of(1), O.None[int]()}))
}
