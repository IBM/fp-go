package identity

import (
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, Of(2), F.Pipe1(1, Map(utils.Double)))
}

func TestChain(t *testing.T) {
	assert.Equal(t, Of(2), F.Pipe1(1, Chain(utils.Double)))
}

func TestAp(t *testing.T) {
	assert.Equal(t, Of(2), F.Pipe1(utils.Double, Ap[int, int](1)))
}
