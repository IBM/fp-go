package json

import (
	"fmt"
	"testing"

	E "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

type Json map[string]any

func TestJsonMarshal(t *testing.T) {

	resRight := Unmarshal[Json]([]byte("{\"a\": \"b\"}"))
	assert.True(t, E.IsRight(resRight))

	resLeft := Unmarshal[Json]([]byte("{\"a\""))
	assert.True(t, E.IsLeft(resLeft))

	res1 := F.Pipe1(
		resRight,
		E.Chain(Marshal[Json]),
	)
	fmt.Println(res1)
}
