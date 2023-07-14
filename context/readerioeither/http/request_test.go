package http

import (
	"context"
	"fmt"
	"testing"

	H "net/http"
)

type PostItem struct {
	UserId uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func TestSendSingleRequest(t *testing.T) {

	client := MakeClient(H.DefaultClient)

	req1 := NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)

	readItem := ReadJson[PostItem](client)

	resp1 := readItem(req1)

	resE := resp1(context.Background())()

	fmt.Println(resE)
}
