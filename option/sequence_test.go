package option

import (
	"testing"

	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

func TestSequenceT(t *testing.T) {
	// one argumemt
	s1 := SequenceT1[int]
	assert.Equal(t, Of(T.MakeTuple1(1)), s1(Of(1)))

	// two arguments
	s2 := SequenceT2[int, string]
	assert.Equal(t, Of(T.MakeTuple2(1, "a")), s2(Of(1), Of("a")))

	// three arguments
	s3 := SequenceT3[int, string, bool]
	assert.Equal(t, Of(T.MakeTuple3(1, "a", true)), s3(Of(1), Of("a"), Of(true)))

	// four arguments
	s4 := SequenceT4[int, string, bool, int]
	assert.Equal(t, Of(T.MakeTuple4(1, "a", true, 2)), s4(Of(1), Of("a"), Of(true), Of(2)))

	// three with one none
	assert.Equal(t, None[T.Tuple3[int, string, bool]](), s3(Of(1), Of("a"), None[bool]()))
}
