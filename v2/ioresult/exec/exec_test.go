//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache LicensVersion 2.0 (the "License");
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

	RA "github.com/IBM/fp-go/v2/array"
	B "github.com/IBM/fp-go/v2/bytes"
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/exec"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

func TestOpenSSL(t *testing.T) {
	// execute the openSSL binary
	version := F.Pipe1(
		Command("openssl")(RA.From("version"))(B.Monoid.Empty()),
		ioeither.Map[error](F.Flow3(
			exec.StdOut,
			B.ToString,
			strings.TrimSpace,
		)),
	)

	assert.True(t, E.IsRight(version()))
}
