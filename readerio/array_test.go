package readerio

import (
	"context"
	"testing"

	A "github.com/ibm/fp-go/array"
	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {
	f := TraverseArray(func(a string) ReaderIO[context.Context, string] {
		return Of[context.Context](a + a)
	})
	ctx := context.Background()
	assert.Equal(t, A.Empty[string](), F.Pipe1(A.Empty[string](), f)(ctx)())
	assert.Equal(t, []string{"aa", "bb"}, F.Pipe1([]string{"a", "b"}, f)(ctx)())
}
