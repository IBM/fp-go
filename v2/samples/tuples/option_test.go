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

package tuples

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	IOEG "github.com/IBM/fp-go/v2/ioeither/generic"
	IOO "github.com/IBM/fp-go/v2/iooption"
)

func TestIOEitherToOption1(t *testing.T) {
	tmpDir := t.TempDir()
	content := []byte("abc")

	resIOO := F.Pipe2(
		content,
		IOEF.WriteFile(filepath.Join(tmpDir, "test.txt"), os.ModePerm),
		IOEG.Fold[IOE.IOEither[error, []byte]](
			IOO.Of[error],
			F.Ignore1of1[[]byte](IOO.None[error]),
		),
	)

	fmt.Println(resIOO())
}

func TestIOEitherToOption2(t *testing.T) {
	tmpDir := t.TempDir()
	content := []byte("abc")

	resIOO := F.Pipe3(
		content,
		IOEF.WriteFile(filepath.Join(tmpDir, "test.txt"), os.ModePerm),
		IOE.Swap[error, []byte],
		IOE.ToIOOption[[]byte, error],
	)

	fmt.Println(resIOO())
}
