package either

import (
	"testing"

	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func TestTraverse(t *testing.T) {
	f := func(n int) O.Option[int] {
		if n >= 2 {
			return O.Of(n)
		}
		return O.None[int]()
	}
	trav := Traverse[string, int, int, O.Option[Either[string, int]]](
		O.Of[Either[string, int]],
		O.MonadMap[int, Either[string, int]],
	)(f)

	assert.Equal(t, O.Of(Left[int]("a")), F.Pipe1(Left[int]("a"), trav))
	assert.Equal(t, O.None[Either[string, int]](), F.Pipe1(Right[string](1), trav))
	assert.Equal(t, O.Of(Right[string](3)), F.Pipe1(Right[string](3), trav))
}

func TestSequence(t *testing.T) {

	seq := Sequence(
		O.Of[Either[string, int]],
		O.MonadMap[int, Either[string, int]],
	)

	assert.Equal(t, O.Of(Right[string](1)), seq(Right[string](O.Of(1))))
	assert.Equal(t, O.Of(Left[int]("a")), seq(Left[O.Option[int]]("a")))
	assert.Equal(t, O.None[Either[string, int]](), seq(Right[string](O.None[int]())))
}
