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

package builder

import (
	"net/http"
	"net/url"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/http/builder"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

func TestBuilderWithQuery(t *testing.T) {
	// add some query
	withLimit := R.WithQueryArg("limit")("10")
	withURL := R.WithURL("http://www.example.org?a=b")

	b := F.Pipe2(
		R.Default,
		withLimit,
		withURL,
	)

	req := F.Pipe3(
		b,
		Requester,
		ioeither.Map[error](func(r *http.Request) *url.URL {
			return r.URL
		}),
		ioeither.ChainFirstIOK[error](func(u *url.URL) io.IO[Void] {
			return io.FromImpure(func() {
				q := u.Query()
				assert.Equal(t, "10", q.Get("limit"))
				assert.Equal(t, "b", q.Get("a"))
			})
		}),
	)

	assert.True(t, E.IsRight(req()))
}
