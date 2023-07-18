package number

import (
	"testing"

	"github.com/stretchr/testify/assert"

	M "github.com/IBM/fp-go/magma"
)

func TestSemigroupIsMagma(t *testing.T) {
	sum := SemigroupSum[int]()

	var magma M.Magma[int] = sum

	assert.Equal(t, sum.Concat(1, 2), magma.Concat(1, 2))
	assert.Equal(t, sum.Concat(1, 2), sum.Concat(2, 1))
}
