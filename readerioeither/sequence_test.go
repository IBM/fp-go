package readerioeither

import (
	"context"
	"testing"

	ET "github.com/IBM/fp-go/either"
	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

func TestSequence2(t *testing.T) {
	// two readers of heterogeneous types
	first := Of[context.Context, error]("a")
	second := Of[context.Context, error](1)

	// compose
	s2 := SequenceT2[context.Context, error, string, int]
	res := s2(first, second)

	ctx := context.Background()
	assert.Equal(t, ET.Right[error](T.MakeTuple2("a", 1)), res(ctx)())
}
