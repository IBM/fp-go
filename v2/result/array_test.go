package result

import (
	"errors"
	"fmt"
	"testing"

	TST "github.com/IBM/fp-go/v2/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestCompactArray(t *testing.T) {
	ar := []Result[string]{
		Of("ok"),
		Left[string](errors.New("err")),
		Of("ok"),
	}

	res := CompactArray(ar)
	assert.Equal(t, 2, len(res))
}

func TestSequenceArray(t *testing.T) {

	s := TST.SequenceArrayTest(
		FromStrictEquals[bool](),
		Pointed[string](),
		Pointed[bool](),
		Functor[[]string, bool](),
		SequenceArray[string],
	)

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}

func TestSequenceArrayError(t *testing.T) {

	s := TST.SequenceArrayErrorTest(
		FromStrictEquals[bool](),
		Left[string],
		Left[bool],
		Pointed[string](),
		Pointed[bool](),
		Functor[[]string, bool](),
		SequenceArray[string],
	)
	// run across four bits
	s(4)(t)
}
