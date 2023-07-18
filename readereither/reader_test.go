package readereither

import (
	"testing"

	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MyContext string

const defaultContext MyContext = "default"

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext, error](1),
		Map[MyContext, error](utils.Double),
	)

	assert.Equal(t, ET.Of[error](2), g(defaultContext))

}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext, error](utils.Double),
		Ap[MyContext, error, int, int](Of[MyContext, error](1)),
	)
	assert.Equal(t, ET.Of[error](2), g(defaultContext))

}

func TestFlatten(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext, string](Of[MyContext, string]("a")),
		Flatten[MyContext, string, string],
	)

	assert.Equal(t, ET.Of[string]("a"), g(defaultContext))
}
