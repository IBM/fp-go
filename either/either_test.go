package either

import (
	"errors"
	"testing"

	F "github.com/ibm/fp-go/function"
	"github.com/ibm/fp-go/internal/utils"
	O "github.com/ibm/fp-go/option"
	S "github.com/ibm/fp-go/string"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	var e Either[error, string]

	assert.Equal(t, Of[error](""), e)
}

func TestIsLeft(t *testing.T) {
	err := errors.New("Some error")
	withError := Left[error, string](err)

	assert.True(t, IsLeft(withError))
	assert.False(t, IsRight(withError))
}

func TestIsRight(t *testing.T) {
	noError := Right[error]("Carsten")

	assert.True(t, IsRight(noError))
	assert.False(t, IsLeft(noError))
}

func TestMapEither(t *testing.T) {

	assert.Equal(t, F.Pipe1(Right[error]("abc"), Map[error](utils.StringLen)), Right[error](3))

	val2 := F.Pipe1(Left[string, string]("s"), Map[string](utils.StringLen))
	exp2 := Left[string, int]("s")

	assert.Equal(t, val2, exp2)
}

func TestUnwrapError(t *testing.T) {
	a := ""
	err := errors.New("Some error")
	withError := Left[error, string](err)

	res, extracted := UnwrapError(withError)
	assert.Equal(t, a, res)
	assert.Equal(t, extracted, err)

}

func TestReduce(t *testing.T) {

	s := S.Semigroup()

	assert.Equal(t, "foobar", F.Pipe1(Right[string]("bar"), Reduce[string](s.Concat, "foo")))
	assert.Equal(t, "foo", F.Pipe1(Left[string, string]("bar"), Reduce[string](s.Concat, "foo")))

}

func TestAp(t *testing.T) {
	f := S.Size

	assert.Equal(t, Right[string](3), F.Pipe1(Right[string](f), Ap[string, string, int](Right[string]("abc"))))
	assert.Equal(t, Left[string, int]("maError"), F.Pipe1(Right[string](f), Ap[string, string, int](Left[string, string]("maError"))))
	assert.Equal(t, Left[string, int]("mabError"), F.Pipe1(Left[string, func(string) int]("mabError"), Ap[string, string, int](Left[string, string]("maError"))))
}

func TestAlt(t *testing.T) {
	assert.Equal(t, Right[string](1), F.Pipe1(Right[string](1), Alt(F.Constant(Right[string](2)))))
	assert.Equal(t, Right[string](1), F.Pipe1(Right[string](1), Alt(F.Constant(Left[string, int]("a")))))
	assert.Equal(t, Right[string](2), F.Pipe1(Left[string, int]("b"), Alt(F.Constant(Right[string](2)))))
	assert.Equal(t, Left[string, int]("b"), F.Pipe1(Left[string, int]("a"), Alt(F.Constant(Left[string, int]("b")))))
}

func TestChainFirst(t *testing.T) {
	f := F.Flow2(S.Size, Right[string, int])

	assert.Equal(t, Right[string]("abc"), F.Pipe1(Right[string]("abc"), ChainFirst(f)))
	assert.Equal(t, Left[string, string]("maError"), F.Pipe1(Left[string, string]("maError"), ChainFirst(f)))
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[string, int, int](F.Constant("a"))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})
	assert.Equal(t, Right[string](1), f(Right[string](1)))
	assert.Equal(t, Left[string, int]("a"), f(Right[string](-1)))
	assert.Equal(t, Left[string, int]("b"), f(Left[string, int]("b")))
}

func TestFromOption(t *testing.T) {
	assert.Equal(t, Left[string, int]("none"), FromOption[string, int](F.Constant("none"))(O.None[int]()))
	assert.Equal(t, Right[string](1), FromOption[string, int](F.Constant("none"))(O.Some(1)))
}
