package readerioeither

import (
	"context"
	"testing"

	A "github.com/ibm/fp-go/array"
	ET "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {
	f := TraverseArray(func(a string) ReaderIOEither[context.Context, string, string] {
		if len(a) > 0 {
			return Right[context.Context, string](a + a)
		}
		return Left[context.Context, string, string]("e")
	})
	ctx := context.Background()
	assert.Equal(t, ET.Right[string](A.Empty[string]()), F.Pipe1(A.Empty[string](), f)(ctx)())
	assert.Equal(t, ET.Right[string]([]string{"aa", "bb"}), F.Pipe1([]string{"a", "b"}, f)(ctx)())
	assert.Equal(t, ET.Left[[]string]("e"), F.Pipe1([]string{"a", ""}, f)(ctx)())
}
