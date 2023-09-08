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
	IOO "github.com/IBM/fp-go/iooption"
	N "github.com/IBM/fp-go/number"
	O "github.com/IBM/fp-go/option"
	M "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"
)

type (
	PostItem struct {
		UserId uint   `json:"userId"`
		Id     uint   `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}

	Player struct {
		Id   int
		Name string
	}

	LocalStorage = map[string]Player
)

var (
	playerAlbert = Player{
		Id:   1,
		Name: "Albert",
	}
	playerTheresa = Player{
		Id:   2,
		Name: "Theresa",
	}
	localStorage = LocalStorage{
		"player1": playerAlbert,
		"player2": playerTheresa,
	}

	// getFromCache :: String -> IO User
	getFromCache = func(name string) IOO.IOOption[Player] {
		return func() O.Option[Player] {
			return M.MonadLookup(localStorage, name)
		}
	}

	// game :: User -> User -> String
	game = F.Curry2(func(a, b Player) string {
		return fmt.Sprintf("%s vs %s", a.Name, b.Name)
	})
)

func (player Player) getName() string {
	return player.Name
}

func (player Player) getId() int {
	return player.Id
}

func (item PostItem) getTitle() string {
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
		R.Map(PostItem.getTitle),
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

func Example_solution10A() {
	safeAdd := F.Curry2(func(a, b O.Option[int]) O.Option[int] {
		return F.Pipe3(
			N.Add[int],
			O.Of[func(int) func(int) int],
			O.Ap[func(int) int](a),
			O.Ap[int](b),
		)
	})

	fmt.Println(safeAdd(O.Of(2))(O.Of(3)))
	fmt.Println(safeAdd(O.None[int]())(O.Of(3)))
	fmt.Println(safeAdd(O.Of(2))(O.None[int]()))

	// Output:
	// Some[int](5)
	// None[int]
	// None[int]
}

func Example_solution10B() {

	safeAdd := F.Curry2(T.Untupled2(F.Flow2(
		O.SequenceTuple2[int, int],
		O.Map(T.Tupled2(N.MonoidSum[int]().Concat)),
	)))

	fmt.Println(safeAdd(O.Of(2))(O.Of(3)))
	fmt.Println(safeAdd(O.None[int]())(O.Of(3)))
	fmt.Println(safeAdd(O.Of(2))(O.None[int]()))

	// Output:
	// Some[int](5)
	// None[int]
	// None[int]
}

func Example_solution10C() {
	// startGame :: IO String
	startGame := F.Pipe2(
		IOO.Of(game),
		IOO.Ap[func(Player) string](getFromCache("player1")),
		IOO.Ap[string](getFromCache("player2")),
	)

	startGameTupled := F.Pipe2(
		T.MakeTuple2("player1", "player2"),
		IOO.TraverseTuple2(getFromCache, getFromCache),
		IOO.Map(T.Tupled2(func(a, b Player) string {
			return fmt.Sprintf("%s vs %s", a.Name, b.Name)
		})),
	)

	fmt.Println(startGame())
	fmt.Println(startGameTupled())

	// Output:
	// Some[string](Albert vs Theresa)
	// Some[string](Albert vs Theresa)
}
