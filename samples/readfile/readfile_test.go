package readfile

import (
	"context"
	"testing"

	R "github.com/ibm/fp-go/context/readerioeither"
	"github.com/ibm/fp-go/context/readerioeither/file"
	E "github.com/ibm/fp-go/either"
	F "github.com/ibm/fp-go/function"
	IO "github.com/ibm/fp-go/io"
	J "github.com/ibm/fp-go/json"
	"github.com/stretchr/testify/assert"
)

type RecordType struct {
	Data string `json:"data"`
}

func TestReadSingleFile(t *testing.T) {

	data := F.Pipe2(
		file.ReadFile("./file.json"),
		R.ChainEitherK(J.Unmarshal[RecordType]),
		R.ChainFirstIOK(IO.Logf[RecordType]("Log: %v")),
	)

	result := data(context.Background())

	assert.Equal(t, E.Of[error](RecordType{"Carsten"}), result())
}
