package option

import (
	"testing"

	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestSequenceArray(t *testing.T) {

	one := Of(1)
	two := Of(2)

	res := F.Pipe1(
		[]Option[int]{one, two},
		SequenceArray[int],
	)

	assert.Equal(t, res, Of([]int{1, 2}))
}
