package reader

import (
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"

	"github.com/IBM/fp-go/internal/utils"
)

func TestMap(t *testing.T) {

	assert.Equal(t, 2, F.Pipe1(Of[string](1), Map[string](utils.Double))(""))
}

func TestAp(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of[int](utils.Double), Ap[int, int, int](Of[int](1)))(0))
}
