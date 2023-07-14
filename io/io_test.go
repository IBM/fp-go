package io

import (
	"math/rand"
	"testing"

	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of(1), Map(utils.Double))())
}

func TestChain(t *testing.T) {
	f := func(n int) IO[int] {
		return Of(n * 2)
	}
	assert.Equal(t, 2, F.Pipe1(Of(1), Chain(f))())
}

func TestAp(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of(utils.Double), Ap[int, int](Of(1)))())
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, 1, F.Pipe1(Of(Of(1)), Flatten[int])())
}

func TestMemoize(t *testing.T) {
	data := Memoize(MakeIO(rand.Int))

	value1 := data()
	value2 := data()

	assert.Equal(t, value1, value2)
}

func TestApFirst(t *testing.T) {

	x := F.Pipe1(
		Of("a"),
		ApFirst[string](Of("b")),
	)

	assert.Equal(t, "a", x())
}

func TestApSecond(t *testing.T) {

	x := F.Pipe1(
		Of("a"),
		ApSecond[string](Of("b")),
	)

	assert.Equal(t, "b", x())
}
