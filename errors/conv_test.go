package errors

import (
	"fmt"
	"testing"

	F "github.com/ibm/fp-go/function"
	O "github.com/ibm/fp-go/option"
	"github.com/stretchr/testify/assert"
)

type MyError struct{}

func (m *MyError) Error() string {
	return "boom"
}

func TestAs(t *testing.T) {
	root := &MyError{}
	err := fmt.Errorf("This is my custom error, %w", root)

	errO := F.Pipe1(
		err,
		As[*MyError](),
	)

	assert.Equal(t, O.Of(root), errO)
}

func TestNotAs(t *testing.T) {
	err := fmt.Errorf("This is my custom error")

	errO := F.Pipe1(
		err,
		As[*MyError](),
	)

	assert.Equal(t, O.None[*MyError](), errO)
}
