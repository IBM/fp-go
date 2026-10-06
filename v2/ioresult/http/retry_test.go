// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	AR "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
	HD "github.com/IBM/fp-go/v2/http/headers"
	"github.com/IBM/fp-go/v2/ioeither"
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

// testPost is the item served by the test server.
var testPost = PostItem{UserID: 1, Id: 1, Title: "title", Body: "body"}

// unresolvableHosts fails the lookup of hosts in the reserved ".invalid" top level
// domain (RFC 2606) with a DNS error, without depending on the network.
func unresolvableHosts() *http.Transport {
	dialer := &net.Dialer{}
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err == nil && strings.HasSuffix(host, ".invalid") {
				return nil, &net.DNSError{Err: "no such host", Name: host, IsNotFound: true}
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}
}

// postServer serves testPost as JSON.
func postServer(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(HD.ContentType, C.JSON)
		_ = json.NewEncoder(w).Encode(testPost)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestRetryHttp(t *testing.T) {
	// URLs to try, the first URLs have a hostname that cannot be resolved
	srv := postServer(t)
	urls := AR.From("http://host1.invalid/posts/1", "http://host2.invalid/posts/1", "http://host3.invalid/posts/1", "http://host4.invalid/posts/1", srv.URL+"/posts/1")
	client := MakeClient(&http.Client{Transport: unresolvableHosts()})

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

	item := ioeither.Retrying(testLogPolicy, action, check)()
	assert.Equal(t, E.Right[error](&testPost), item)
}
