package readerioeither

import (
	"context"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	R "github.com/IBM/fp-go/reader"
	RIO "github.com/IBM/fp-go/readerio"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		Map[context.Context, error](utils.Double),
	)

	assert.Equal(t, E.Of[error](2), g(context.Background())())
}

func TestOrLeft(t *testing.T) {
	f := OrLeft[string, context.Context, string, int](func(s string) RIO.ReaderIO[context.Context, string] {
		return RIO.Of[context.Context](s + "!")
	})

	g1 := F.Pipe1(
		Right[context.Context, string](1),
		f,
	)

	g2 := F.Pipe1(
		Left[context.Context, int]("a"),
		f,
	)

	assert.Equal(t, E.Of[string](1), g1(context.Background())())
	assert.Equal(t, E.Left[int]("a!"), g2(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[context.Context, error](utils.Double),
		Ap[int](Right[context.Context, error](1)),
	)

	assert.Equal(t, E.Right[error](2), g(context.Background())())
}

func TestChainReaderK(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		ChainReaderK[context.Context, error](func(v int) R.Reader[context.Context, string] {
			return R.Of[context.Context](fmt.Sprintf("%d", v))
		}),
	)

	assert.Equal(t, E.Right[error]("1"), g(context.Background())())
}
