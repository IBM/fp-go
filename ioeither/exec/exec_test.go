//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

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
