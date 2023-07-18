package either

import (
	"testing"

	F "github.com/ibm/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	l := Logger[error, string]()

	r := Right[error]("test")

	res := F.Pipe1(
		r,
		l("out"),
	)

	assert.Equal(t, r, res)
}
