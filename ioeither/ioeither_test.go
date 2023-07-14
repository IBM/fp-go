package ioeither

import (
	"fmt"
	"testing"

	E "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/utils"
	I "github.com/ibm/fp-go/io"
	IG "github.com/ibm/fp-go/io/generic"
	O "github.com/ibm/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, E.Of[error](2), F.Pipe1(
		Of[error](1),
		Map[error](utils.Double),
	)())

}

func TestChainEitherK(t *testing.T) {
	f := ChainEitherK(func(n int) E.Either[string, int] {
		if n > 0 {
			return E.Of[string](n)
		}
		return E.Left[int]("a")

	})
	assert.Equal(t, E.Right[string](1), f(Right[string](1))())
	assert.Equal(t, E.Left[int]("a"), f(Right[string](-1))())
	assert.Equal(t, E.Left[int]("b"), f(Left[int]("b"))())
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](F.Constant("a"))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})

	assert.Equal(t, E.Right[string](1), f(Right[string](1))())
	assert.Equal(t, E.Left[int]("a"), f(Right[string](-1))())
	assert.Equal(t, E.Left[int]("b"), f(Left[int]("b"))())
}

func TestFromOption(t *testing.T) {
	f := FromOption[int](F.Constant("a"))
	assert.Equal(t, E.Right[string](1), f(O.Some(1))())
	assert.Equal(t, E.Left[int]("a"), f(O.None[int]())())
}

func TestChainIOK(t *testing.T) {
	f := ChainIOK[string](func(n int) I.IO[string] {
		return I.MakeIO(func() string {
			return fmt.Sprintf("%d", n)
		})
	})

	assert.Equal(t, E.Right[string]("1"), f(Right[string](1))())
	assert.Equal(t, E.Left[string, string]("b"), f(Left[int]("b"))())
}

func TestChainWithIO(t *testing.T) {

	r := F.Pipe1(
		Of[error]("test"),
		// sad, we need the generics version ...
		IG.Map[IOEither[error, string], I.IO[bool]](E.IsRight[error, string]),
	)

	assert.True(t, r())
}

func TestChainFirst(t *testing.T) {
	f := func(a string) IOEither[string, int] {
		if len(a) > 2 {
			return Of[string](len(a))
		}
		return Left[int]("foo")
	}
	good := Of[string]("foo")
	bad := Of[string]("a")
	ch := ChainFirst(f)

	assert.Equal(t, E.Of[string]("foo"), F.Pipe1(good, ch)())
	assert.Equal(t, E.Left[string, string]("foo"), F.Pipe1(bad, ch)())
}

func TestChainFirstIOK(t *testing.T) {
	f := func(a string) I.IO[int] {
		return I.Of(len(a))
	}
	good := Of[string]("foo")
	ch := ChainFirstIOK[string](f)

	assert.Equal(t, E.Of[string]("foo"), F.Pipe1(good, ch)())
}

func TestApFirst(t *testing.T) {

	x := F.Pipe1(
		Of[error]("a"),
		ApFirst[error, string](Of[error]("b")),
	)

	assert.Equal(t, E.Of[error]("a"), x())
}

func TestApSecond(t *testing.T) {

	x := F.Pipe1(
		Of[error]("a"),
		ApSecond[error, string](Of[error]("b")),
	)

	assert.Equal(t, E.Of[error]("b"), x())
}
