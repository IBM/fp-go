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

package readfile

import (
	"context"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	R "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerioresult/file"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	J "github.com/IBM/fp-go/v2/json"
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
