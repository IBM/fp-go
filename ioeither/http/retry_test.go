package http

import (
	"net"
	"net/http"
	"testing"
	"time"

	AR "github.com/ibm/fp-go/array"
	E "github.com/ibm/fp-go/either"
	HE "github.com/ibm/fp-go/either/http"
	"github.com/ibm/fp-go/errors"
	F "github.com/ibm/fp-go/function"
	IOE "github.com/ibm/fp-go/ioeither"
	O "github.com/ibm/fp-go/option"
	R "github.com/ibm/fp-go/retry"
	"github.com/stretchr/testify/assert"
)

var expLogBackoff = R.ExponentialBackoff(250 * time.Millisecond)

// our retry policy with a 1s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

type PostItem struct {
	UserId uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestRetryHttp(t *testing.T) {
	// URLs to try, the first URLs have an invalid hostname
	urls := AR.From("https://jsonplaceholder1.typicode.com/posts/1", "https://jsonplaceholder2.typicode.com/posts/1", "https://jsonplaceholder3.typicode.com/posts/1", "https://jsonplaceholder4.typicode.com/posts/1", "https://jsonplaceholder.typicode.com/posts/1")
	client := MakeClient(&http.Client{})

	action := func(status R.RetryStatus) IOE.IOEither[error, *PostItem] {
		return F.Pipe2(
			HE.GetRequest(urls[status.IterNumber]),
			IOE.FromEither[error, *http.Request],
			IOE.Chain(ReadJson[*PostItem](client)),
		)
	}

	check := E.Fold(
		F.Flow2(
			errors.As[*net.DNSError](),
			O.Fold(F.Constant(false), F.Constant1[*net.DNSError](true)),
		),
		F.Constant1[*PostItem](false),
	)

	item := IOE.Retrying(testLogPolicy, action, check)()
	assert.True(t, E.IsRight(item))
}
