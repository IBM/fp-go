package json

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

type TestType struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func TestToType(t *testing.T) {

	generic := map[string]any{"a": "value", "b": 1}

	assert.True(t, E.IsRight(ToTypeE[TestType](generic)))
	assert.True(t, E.IsRight(ToTypeE[TestType](&generic)))

	assert.Equal(t, E.Right[error](&TestType{A: "value", B: 1}), ToTypeE[*TestType](&generic))
	assert.Equal(t, E.Right[error](TestType{A: "value", B: 1}), F.Pipe1(ToTypeE[*TestType](&generic), E.Map[error](F.Deref[TestType])))
}
