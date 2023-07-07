package option

import (
	"fmt"
	"testing"

	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {

	assert.Equal(t, 2, F.Pipe1(None[int](), Reduce(utils.Sum, 2)))
	assert.Equal(t, 5, F.Pipe1(Some(3), Reduce(utils.Sum, 2)))
}

func TestIsNone(t *testing.T) {
	assert.True(t, IsNone(None[int]()))
	assert.False(t, IsNone(Of(1)))
}

func TestIsSome(t *testing.T) {
	assert.True(t, IsSome(Of(1)))
	assert.False(t, IsSome(None[int]()))
}

func TestMapOption(t *testing.T) {

	assert.Equal(t, F.Pipe1(Some(2), Map(utils.Double)), Some(4))

	assert.Equal(t, F.Pipe1(None[int](), Map(utils.Double)), None[int]())
}

func TestTryCachOption(t *testing.T) {

	res := TryCatch(utils.Error)

	assert.Equal(t, None[int](), res)
}

func TestAp(t *testing.T) {
	assert.Equal(t, Some(4), F.Pipe1(
		Some(utils.Double),
		Ap[int, int](Some(2)),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		Some(utils.Double),
		Ap[int, int](None[int]()),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[func(int) int](),
		Ap[int, int](Some(2)),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[func(int) int](),
		Ap[int, int](None[int]()),
	))
}

func TestChain(t *testing.T) {
	f := func(n int) Option[int] { return Some(n * 2) }
	g := func(_ int) Option[int] { return None[int]() }

	assert.Equal(t, Some(2), F.Pipe1(
		Some(1),
		Chain(f),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[int](),
		Chain(f),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		Some(1),
		Chain(g),
	))

	assert.Equal(t, None[int](), F.Pipe1(
		None[int](),
		Chain(g),
	))
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, Of(1), F.Pipe1(Of(Of(1)), Flatten[int]))
}

func TestFold(t *testing.T) {
	f := F.Constant("none")
	g := func(s string) string { return fmt.Sprintf("some%d", len(s)) }

	fold := Fold(f, g)

	assert.Equal(t, "none", fold(None[string]()))
	assert.Equal(t, "some3", fold(Some("abc")))
}

func TestFromPredicate(t *testing.T) {
	p := func(n int) bool { return n > 2 }
	f := FromPredicate(p)

	assert.Equal(t, None[int](), f(1))
	assert.Equal(t, Some(3), f(3))
}

func TestAlt(t *testing.T) {
	assert.Equal(t, Some(1), F.Pipe1(Some(1), Alt(F.Constant(Some(2)))))
	assert.Equal(t, Some(2), F.Pipe1(Some(2), Alt(F.Constant(None[int]()))))
	assert.Equal(t, Some(1), F.Pipe1(None[int](), Alt(F.Constant(Some(1)))))
	assert.Equal(t, None[int](), F.Pipe1(None[int](), Alt(F.Constant(None[int]()))))
}
