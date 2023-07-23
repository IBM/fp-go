// Copyright (c) 2023 IBM Corp.
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

	req1 := MakeGetRequest("https://jsonplaceholder.typicode.com/posts/1")

	readItem := ReadJson[PostItem](client)

	resp1 := readItem(req1)

	resE := resp1(context.Background())()

	fmt.Println(resE)
}
