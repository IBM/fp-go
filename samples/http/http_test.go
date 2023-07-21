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
	T "github.com/IBM/fp-go/tuple"
	"github.com/stretchr/testify/assert"
)

type PostItem struct {
	UserId uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type CatFact struct {
	Fact string `json:"fact"`
}

func idxToUrl(idx int) string {
	return fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", idx+1)
}

// TestMultipleHttpRequests shows how to execute multiple HTTP requests in parallel assuming
// that the response structure of all requests is identical, which is why we can use [R.TraverseArray]
func TestMultipleHttpRequests(t *testing.T) {
	// prepare the http client
	client := H.MakeClient(HTTP.DefaultClient)
	// readSinglePost sends a GET request and parses the response as [PostItem]
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

// TestHeterogeneousHttpRequests shows how to execute multiple HTTP requests in parallel when
// the response structure of these requests is different. We use [R.TraverseTuple2] to account for the different types
func TestHeterogeneousHttpRequests(t *testing.T) {
	// prepare the http client
	client := H.MakeClient(HTTP.DefaultClient)
	// readSinglePost sends a GET request and parses the response as [PostItem]
	readSinglePost := H.ReadJson[PostItem](client)
	// readSingleCatFact sends a GET request and parses the response as [CatFact]
	readSingleCatFact := H.ReadJson[CatFact](client)

	data := F.Pipe3(
		T.MakeTuple2("https://jsonplaceholder.typicode.com/posts/1", "https://catfact.ninja/fact"),
		T.Map2(H.MakeGetRequest, H.MakeGetRequest),
		R.TraverseTuple2(
			readSinglePost,
			readSingleCatFact,
		),
		R.ChainFirstIOK(IO.Logf[T.Tuple2[PostItem, CatFact]]("Log Result: %v")),
	)

	result := data(context.Background())

	fmt.Println(result())
}
