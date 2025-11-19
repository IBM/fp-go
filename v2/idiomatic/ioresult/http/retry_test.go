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

package http

import (
	"net"
	"net/http"
	"testing"
	"time"

	AR "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	E "github.com/IBM/fp-go/v2/idiomatic/result"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

var expLogBackoff = R.ExponentialBackoff(250 * time.Millisecond)

// our retry policy with a 1s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

type PostItem struct {
	UserID uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestRetryHttp(t *testing.T) {
	// URLs to try, the first URLs have an invalid hostname
	urls := AR.From("https://jsonplaceholder1.typicode.com/posts/1", "https://jsonplaceholder2.typicode.com/posts/1", "https://jsonplaceholder3.typicode.com/posts/1", "https://jsonplaceholder4.typicode.com/posts/1", "https://jsonplaceholder.typicode.com/posts/1")
	client := MakeClient(&http.Client{})

	action := func(status R.RetryStatus) IOResult[*PostItem] {
		return F.Pipe1(
			MakeGetRequest(urls[status.IterNumber]),
			ReadJSON[*PostItem](client),
		)
	}

	check := E.Fold(
		F.Flow2(
			errors.As[*net.DNSError](),
			O.IsSome[*net.DNSError],
		),
		F.Constant1[*PostItem](false),
	)

	_, err := ioresult.Retrying(testLogPolicy, action, check)()
	assert.NoError(t, err)
}
