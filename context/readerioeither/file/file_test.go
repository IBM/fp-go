package file

import (
	"context"
	"fmt"

	RIO "github.com/IBM/fp-go/context/readerio"
	R "github.com/IBM/fp-go/context/readerioeither"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	J "github.com/IBM/fp-go/json"
)

type RecordType struct {
	Data string `json:"data"`
}

func getData(r RecordType) string {
	return r.Data
}

func ExampleReadFile() {

	data := F.Pipe4(
		ReadFile("./data/file.json"),
		R.ChainEitherK(J.Unmarshal[RecordType]),
		R.ChainFirstIOK(IO.Logf[RecordType]("Log: %v")),
		R.Map(getData),
		R.GetOrElse(F.Flow2(
			errors.ToString,
			RIO.Of[string],
		)),
	)

	result := data(context.Background())

	fmt.Println(result())

	// Output: Carsten
}
