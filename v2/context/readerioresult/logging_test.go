package readerioresult

import (
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestLoggingContext(t *testing.T) {

	data := F.Pipe2(
		Of("Sample"),
		LogEntryExit[string]("TestLoggingContext1"),
		LogEntryExit[string]("TestLoggingContext2"),
	)

	assert.Equal(t, result.Of("Sample"), data(t.Context())())
}

func TestLoggingContextWithLogger(t *testing.T) {

	data := F.Pipe4(
		Of("Sample"),
		LogEntryExit[string]("TestLoggingContext1"),
		Map(strings.ToUpper),
		ChainFirst(F.Flow2(
			WithLoggingID[string],
			ChainIOK(io.Logf[pair.Pair[LoggingID, string]]("Prefix: %s")),
		)),
		LogEntryExit[string]("TestLoggingContext2"),
	)

	assert.Equal(t, result.Of("SAMPLE"), data(t.Context())())
}
