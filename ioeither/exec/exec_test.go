package exec

import (
	"strings"
	"testing"

	RA "github.com/ibm/fp-go/array"
	B "github.com/ibm/fp-go/bytes"
	E "github.com/ibm/fp-go/either"
	"github.com/ibm/fp-go/exec"
	F "github.com/ibm/fp-go/function"
	IOE "github.com/ibm/fp-go/ioeither"
	"github.com/stretchr/testify/assert"
)

func TestOpenSSL(t *testing.T) {
	// execute the openSSL binary
	version := F.Pipe1(
		Command("openssl")(RA.From("version"))(B.Monoid.Empty()),
		IOE.Map[error](F.Flow3(
			exec.StdOut,
			B.ToString,
			strings.TrimSpace,
		)),
	)

	assert.True(t, E.IsRight(version()))
}
