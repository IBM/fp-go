package readerio

import (
	"context"
	"testing"

	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context](1),
		Map[context.Context](utils.Double),
	)

	assert.Equal(t, 2, g(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[context.Context](utils.Double),
		Ap[int](Of[context.Context](1)),
	)

	assert.Equal(t, 2, g(context.Background())())
}
