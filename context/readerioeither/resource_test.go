// Copyright (c) 2023 IBM Corp.
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

package readerioeither

import (
	"context"
	"io"
	"os"

	B "github.com/IBM/fp-go/bytes"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
)

var (
	openFile = F.Flow3(
		IOE.Eitherize1(os.Open),
		FromIOEither[*os.File],
		ChainFirstIOK(F.Flow2(
			(*os.File).Name,
			IO.Logf[string]("Opened file [%s]"),
		)),
	)
)

func closeFile(f *os.File) ReaderIOEither[string] {
	return F.Pipe1(
		TryCatch(func(_ context.Context) func() (string, error) {
			return func() (string, error) {
				return f.Name(), f.Close()
			}
		}),
		ChainFirstIOK(IO.Logf[string]("Closed file [%s]")),
	)
}

func ExampleWithResource() {

	stringReader := WithResource[string](openFile("data/file.txt"), closeFile)

	rdr := stringReader(func(f *os.File) ReaderIOEither[string] {
		return F.Pipe2(
			TryCatch(func(_ context.Context) func() ([]byte, error) {
				return func() ([]byte, error) {
					return io.ReadAll(f)
				}
			}),
			ChainFirstIOK(F.Flow2(
				B.Size,
				IO.Logf[int]("Read content of length [%d]"),
			)),
			Map(B.ToString),
		)
	})

	contentIOE := F.Pipe2(
		context.Background(),
		rdr,
		IOE.ChainFirstIOK[error](IO.Printf[string]("Content: %s")),
	)

	contentIOE()

	// Output: Content: Carsten
}
