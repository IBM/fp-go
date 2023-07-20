package http

import (
	"context"
	"fmt"
	"testing"

	HTTP "net/http"

	A "github.com/IBM/fp-go/array"
	R "github.com/IBM/fp-go/context/readerioeither"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	"github.com/stretchr/testify/assert"
)

type PostItem struct {
	UserId uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func idxToUrl(idx int) string {
	return fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", idx+1)
}

func TestMultipleHttpRequests(t *testing.T) {
	// prepare the http client
	client := H.MakeClient(HTTP.DefaultClient)
	readSinglePost := H.ReadJson[PostItem](client)

	// total number of http requests
	count := 10

	data := F.Pipe3(
		A.MakeBy(count, idxToUrl),
		R.TraverseArray(F.Flow3(
			H.MakeGetRequest,
			readSinglePost,
			R.ChainFirstIOK(IO.Logf[PostItem]("Log Single: %v")),
		)),
		R.ChainFirstIOK(IO.Logf[[]PostItem]("Log Result: %v")),
		R.Map(A.Size[PostItem]),
	)

	result := data(context.Background())

	assert.Equal(t, E.Of[error](count), result())
}
