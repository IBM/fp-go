package readfile

import (
	"context"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/array"
	R "github.com/IBM/fp-go/context/readerioeither"
	"github.com/IBM/fp-go/context/readerioeither/file"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	J "github.com/IBM/fp-go/json"
	"github.com/stretchr/testify/assert"
)

type RecordType struct {
	Data string `json:"data"`
}

// TestReadSingleFile reads the content of a file from disk and parses it into
// a struct
func TestReadSingleFile(t *testing.T) {

	data := F.Pipe2(
		file.ReadFile("./data/file.json"),
		R.ChainEitherK(J.Unmarshal[RecordType]),
		R.ChainFirstIOK(IO.Logf[RecordType]("Log: %v")),
	)

	result := data(context.Background())

	assert.Equal(t, E.Of[error](RecordType{"Carsten"}), result())
}

func idxToFilename(idx int) string {
	return fmt.Sprintf("./data/file%d.json", idx+1)
}

// TestReadMultipleFiles reads the content of a multiple from disk and parses them into
// structs
func TestReadMultipleFiles(t *testing.T) {

	data := F.Pipe2(
		A.MakeBy(3, idxToFilename),
		R.TraverseArray(F.Flow3(
			file.ReadFile,
			R.ChainEitherK(J.Unmarshal[RecordType]),
			R.ChainFirstIOK(IO.Logf[RecordType]("Log Single: %v")),
		)),
		R.ChainFirstIOK(IO.Logf[[]RecordType]("Log Result: %v")),
	)

	result := data(context.Background())

	assert.Equal(t, E.Of[error](A.From(RecordType{"file1"}, RecordType{"file2"}, RecordType{"file3"})), result())
}
