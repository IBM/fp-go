package readereither

import (
	"context"
	"testing"

	A "github.com/ibm/fp-go/array"
	ET "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestSequenceArray(t *testing.T) {

	n := 10

	readers := A.MakeBy(n, Of[context.Context, error, int])
	exp := ET.Of[error](A.MakeBy(n, F.Identity[int]))

	g := F.Pipe1(
		readers,
		SequenceArray[context.Context, error, int],
	)

	assert.Equal(t, exp, g(context.Background()))
}
