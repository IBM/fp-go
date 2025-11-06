// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package example1

import (
	"os"

	B "github.com/IBM/fp-go/v2/bytes"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	J "github.com/IBM/fp-go/v2/json"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

type (
	WriterType = func([]byte) IOE.IOEither[error, []byte]

	Dependencies struct {
		Writer WriterType
	}
)

func getWriter(deps *Dependencies) WriterType {
	return deps.Writer
}

// SerializeToWriter marshals the input to JSON and persists the result via the [Writer] passed in via the [*Dependencies]
func SerializeToWriter[A any](data A) RIOE.ReaderIOEither[*Dependencies, error, []byte] {
	return F.Pipe1(
		RIOE.Ask[*Dependencies, error](),
		RIOE.ChainIOEitherK[*Dependencies](F.Flow2(
			getWriter,
			F.Pipe2(
				data,
				J.Marshal[A],
				E.Fold(F.Flow2(
					IOE.Left[[]byte, error],
					F.Constant1[WriterType, IOE.IOEither[error, []byte]],
				), I.Ap[IOE.IOEither[error, []byte], []byte]),
			),
		)),
	)
}

func ExampleReaderIOEither() {

	// writeToStdOut implements a writer to stdout
	writeToStdOut := func(data []byte) IOE.IOEither[error, []byte] {
		return IOE.TryCatchError(func() ([]byte, error) {
			_, err := os.Stdout.Write(data)
			return data, err
		})
	}

	deps := Dependencies{
		Writer: writeToStdOut,
	}

	data := map[string]string{
		"a": "b",
		"c": "d",
	}

	// writeData will persist to a configurable target
	writeData := F.Pipe1(
		SerializeToWriter(data),
		RIOE.Map[*Dependencies, error](B.ToString),
	)

	_ = writeData(&deps)()

	// Output:
	// {"a":"b","c":"d"}
}
