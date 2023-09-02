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

package mostlyadequate

import (
	"context"
	"fmt"
	"net/http"

	R "github.com/IBM/fp-go/context/readerioeither"
	H "github.com/IBM/fp-go/context/readerioeither/http"
	F "github.com/IBM/fp-go/function"
)

type PostItem struct {
	UserId uint   `json:"userId"`
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func getTitle(item PostItem) string {
	return item.Title
}

func idxToUrl(idx int) string {
	return fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", idx+1)
}

func renderString(destinations string) func(string) string {
	return func(events string) string {
		return fmt.Sprintf("<div>Destinations: [%s], Events: [%s]</div>", destinations, events)
	}
}

func Example_renderPage() {
	// prepare the http client
	client := H.MakeClient(http.DefaultClient)

	// get returns the title of the nth item from the REST service
	get := F.Flow4(
		idxToUrl,
		H.MakeGetRequest,
		H.ReadJson[PostItem](client),
		R.Map(getTitle),
	)

	res := F.Pipe2(
		R.Of(renderString),                // start with a function with 2 unresolved arguments
		R.Ap[func(string) string](get(1)), // resolve the first argument
		R.Ap[string](get(2)),              // in parallel resolve the second argument
	)

	// finally invoke in context and start
	fmt.Println(res(context.TODO())())

	// Output:
	// Right[<nil>, string](<div>Destinations: [qui est esse], Events: [ea molestias quasi exercitationem repellat qui ipsa sit aut]</div>)

}
