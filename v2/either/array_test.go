package either

import (
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	TST "github.com/IBM/fp-go/v2/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestCompactArray(t *testing.T) {
	ar := A.From(
		Of[string]("ok"),
		Left[string]("err"),
		Of[string]("ok"),
	)

	res := CompactArray(ar)
	assert.Equal(t, 2, len(res))
}

func TestSequenceArray(t *testing.T) {

	s := TST.SequenceArrayTest(
		FromStrictEquals[error, bool](),
		Pointed[error, string](),
		Pointed[error, bool](),
		Functor[error, []string, bool](),
		SequenceArray[error, string],
	)

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}

func TestSequenceArrayError(t *testing.T) {

	s := TST.SequenceArrayErrorTest(
		FromStrictEquals[error, bool](),
		Left[string, error],
		Left[bool, error],
		Pointed[error, string](),
		Pointed[error, bool](),
		Functor[error, []string, bool](),
		SequenceArray[error, string],
	)
	// run across four bits
	s(4)(t)
}
