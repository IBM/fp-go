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

package examples

import (
	"encoding/json"
	"fmt"
	"os"

	B "github.com/IBM/fp-go/v2/bytes"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioeither/file"
	J "github.com/IBM/fp-go/v2/json"
	T "github.com/IBM/fp-go/v2/tuple"
)

type Sample struct {
	Value int `json:"a"`
}

func (s Sample) getValue() int {
	return s.Value
}

func Example_io_flow() {

	// IOE.IOEither[error, string]
	text := F.Pipe2(
		"data/file1.txt",
		file.ReadFile,
		IOE.Map[error](B.ToString),
	)

	// IOE.IOEither[error, int]
	value := F.Pipe3(
		"data/file2.json",
		file.ReadFile,
		IOE.ChainEitherK(J.Unmarshal[Sample]),
		IOE.Map[error](Sample.getValue),
	)

	// IOE.IOEither[error, string]
	result := F.Pipe1(
		IOE.SequenceT2(text, value),
		IOE.Map[error](func(res T.Tuple2[string, int]) string {
			return fmt.Sprintf("Text: %s, Number: %d", res.F1, res.F2)
		}),
	)

	fmt.Println(result())

	// Output:
	// Right[string](Text: Some data, Number: 10)

}

func io_flow_idiomatic() error {

	// []byte
	file1AsBytes, err := os.ReadFile("data/file1.txt")
	if err != nil {
		return err
	}
	// string
	text := string(file1AsBytes)

	// []byte
	file2AsBytes, err := os.ReadFile("data/file2.json")
	if err != nil {
		return err
	}
	var value Sample
	if err := json.Unmarshal(file2AsBytes, &value); err != nil {
		return err
	}
	// string
	result := fmt.Sprintf("Text: %s, Number: %d", text, value.Value)

	fmt.Println(result)

	return nil
}

func Example_io_flow_idiomatic() {
	if err := io_flow_idiomatic(); err != nil {
		panic(err)
	}

	// Output:
	// Text: Some data, Number: 10
}
